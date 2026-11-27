package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// requestDirList 调用 ListDir 并解析响应。
func requestDirList(t *testing.T, path, ext string) (*httptest.ResponseRecorder, dirListResponse) {
	t.Helper()
	query := url.Values{}
	if path != "" {
		query.Set("path", path)
	}
	if ext != "" {
		query.Set("ext", ext)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/system/fs/list?"+query.Encode(), nil)
	rec := httptest.NewRecorder()
	NewHandler().ListDir(rec, req)

	var body dirListResponse
	if rec.Code == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("解析响应失败: %v (body=%s)", err, rec.Body.String())
		}
	}
	return rec, body
}

func TestListDirRoots(t *testing.T) {
	rec, body := requestDirList(t, "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("状态码 = %d，期望 200: %s", rec.Code, rec.Body.String())
	}
	if len(body.Entries) == 0 {
		t.Fatal("盘符列表为空")
	}
	if body.Path != "" || body.Parent != "" {
		t.Fatalf("根列表的 path/parent 应为空，实际 %q/%q", body.Path, body.Parent)
	}
	for _, e := range body.Entries {
		if !e.IsDir {
			t.Fatalf("%q 应为目录", e.Name)
		}
	}
}

func TestListDirSortsAndFilters(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"b", "A"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"x.jar", "y.txt", "z.JAR"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	rec, body := requestDirList(t, dir, "jar")
	if rec.Code != http.StatusOK {
		t.Fatalf("状态码 = %d，期望 200: %s", rec.Code, rec.Body.String())
	}
	if body.Path != dir {
		t.Fatalf("path = %q，期望 %q", body.Path, dir)
	}
	if body.Parent != filepath.Dir(dir) {
		t.Fatalf("parent = %q，期望 %q", body.Parent, filepath.Dir(dir))
	}
	if body.Truncated {
		t.Fatal("小目录不应被截断")
	}

	// 目录在前按名排序，文件按扩展名过滤后同样按名排序
	want := []string{"A", "b", "x.jar", "z.JAR"}
	if len(body.Entries) != len(want) {
		t.Fatalf("条目数 = %d，期望 %d: %+v", len(body.Entries), len(want), body.Entries)
	}
	for i, name := range want {
		if body.Entries[i].Name != name {
			t.Fatalf("第 %d 项 = %q，期望 %q", i, body.Entries[i].Name, name)
		}
	}
	for _, e := range body.Entries[:2] {
		if !e.IsDir {
			t.Fatalf("%q 应为目录", e.Name)
		}
	}
	for _, e := range body.Entries[2:] {
		if e.IsDir || e.Size != 1 {
			t.Fatalf("%q 应为大小 1 的文件，实际 isDir=%v size=%d", e.Name, e.IsDir, e.Size)
		}
		if e.ModTime == "" {
			t.Fatalf("%q 缺少修改时间", e.Name)
		}
		if e.Path != filepath.Join(dir, e.Name) {
			t.Fatalf("%q 的 path = %q", e.Name, e.Path)
		}
	}
}

func TestListDirWithoutExtFilterKeepsAllFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.jar", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	_, body := requestDirList(t, dir, "")
	if len(body.Entries) != 2 {
		t.Fatalf("条目数 = %d，期望 2: %+v", len(body.Entries), body.Entries)
	}
}

func TestListDirErrors(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-dir.txt")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if rec, _ := requestDirList(t, filepath.Join(t.TempDir(), "missing"), ""); rec.Code != http.StatusNotFound {
		t.Fatalf("不存在目录的状态码 = %d，期望 404", rec.Code)
	}
	if rec, _ := requestDirList(t, file, ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("文件路径的状态码 = %d，期望 400", rec.Code)
	}
}

func TestParentDir(t *testing.T) {
	if runtime.GOOS != "windows" {
		if got := parentDir("/"); got != "" {
			t.Fatalf("根目录的 parent = %q，期望空", got)
		}
		return
	}
	if got := parentDir(`C:\`); got != "" {
		t.Fatalf(`C:\ 的 parent = %q，期望空`, got)
	}
	if got := parentDir(`C:\a\b`); got != `C:\a` {
		t.Fatalf(`C:\a\b 的 parent = %q，期望 C:\a`, got)
	}
}

func TestParseExtFilter(t *testing.T) {
	got := parseExtFilter(" JAR, .zip ,,")
	if len(got) != 2 {
		t.Fatalf("扩展名数量 = %d，期望 2: %v", len(got), got)
	}
	for _, want := range []string{".jar", ".zip"} {
		if _, ok := got[want]; !ok {
			t.Fatalf("缺少 %s: %v", want, got)
		}
	}
	if parseExtFilter("  ") != nil {
		t.Fatal("空白参数应返回 nil")
	}
}
