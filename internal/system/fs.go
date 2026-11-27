package system

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// maxDirEntries 限制单次目录列举的条目数，防止超大目录（如 node_modules）把响应撑爆。
const maxDirEntries = 2000

// dirEntry 是目录中的一项，Path 为绝对路径，前端直接用它导航，无需自己拼路径。
type dirEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

// dirListResponse 是 GET /api/system/fs/list 的响应。
// Path 为空表示当前展示的是「此电脑」（各盘符根目录）。
type dirListResponse struct {
	Path      string     `json:"path"`
	Parent    string     `json:"parent"`
	Entries   []dirEntry `json:"entries"`
	Truncated bool       `json:"truncated"`
}

// ListDir 处理 GET /api/system/fs/list?path=<绝对路径>&ext=jar,zip：
// 列举目录内容供前端文件浏览器使用。path 为空时返回盘符列表。
// ext 为可选的扩展名白名单（逗号分隔），仅用于过滤文件，目录始终返回。
func (h *Handler) ListDir(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := strings.TrimSpace(query.Get("path"))
	exts := parseExtFilter(query.Get("ext"))

	if path == "" {
		roots, err := volumeRoots()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, dirListResponse{Entries: roots})
		return
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "路径不合法: "+path)
		return
	}
	info, err := os.Stat(abs)
	if err != nil {
		writeError(w, http.StatusNotFound, "无法打开目录: "+abs)
		return
	}
	if !info.IsDir() {
		writeError(w, http.StatusBadRequest, "不是目录: "+abs)
		return
	}

	entries, truncated, err := readDirEntries(abs, exts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "读取目录失败: "+abs)
		return
	}
	writeJSON(w, http.StatusOK, dirListResponse{
		Path:      abs,
		Parent:    parentDir(abs),
		Entries:   entries,
		Truncated: truncated,
	})
}

// readDirEntries 列举 dir 下的条目：目录在前，其后按名称不区分大小写排序。
// exts 非空时只保留扩展名命中的文件。
func readDirEntries(dir string, exts map[string]struct{}) ([]dirEntry, bool, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, false, err
	}

	entries := make([]dirEntry, 0, len(items))
	truncated := false
	for _, item := range items {
		info, err := item.Info()
		if err != nil {
			continue // 权限不足或失效的符号链接，跳过而不是整页失败
		}
		if info.Mode()&os.ModeSymlink != 0 {
			// ReadDir 不跟随符号链接，这里跟随一次以判断目标是否为目录
			if target, err := os.Stat(filepath.Join(dir, item.Name())); err == nil {
				info = target
			}
		}

		isDir := info.IsDir()
		if !isDir && len(exts) > 0 {
			if _, ok := exts[strings.ToLower(filepath.Ext(item.Name()))]; !ok {
				continue
			}
		}

		entries = append(entries, dirEntry{
			Name:    item.Name(),
			Path:    filepath.Join(dir, item.Name()),
			IsDir:   isDir,
			Size:    info.Size(),
			ModTime: info.ModTime().Format(time.RFC3339),
		})
		if len(entries) >= maxDirEntries {
			truncated = true
			break
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})
	return entries, truncated, nil
}

// volumeRoots 返回可供浏览的根目录：Windows 为所有可访问的盘符，其他平台为 /。
func volumeRoots() ([]dirEntry, error) {
	if runtime.GOOS != "windows" {
		return []dirEntry{{Name: "/", Path: "/", IsDir: true}}, nil
	}
	roots := make([]dirEntry, 0, 8)
	for c := 'A'; c <= 'Z'; c++ {
		root := string(c) + `:\`
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			continue // 未分配盘符或无介质（光驱、未挂载的分区）
		}
		roots = append(roots, dirEntry{
			Name:  string(c) + ":",
			Path:  root,
			IsDir: true,
		})
	}
	return roots, nil
}

// parentDir 返回上级目录；path 已是根目录时返回空字符串。
func parentDir(path string) string {
	parent := filepath.Dir(path)
	if parent == path {
		return ""
	}
	return parent
}

// parseExtFilter 解析 "jar, .zip" 形式的扩展名白名单，统一成小写带点的形式。
func parseExtFilter(raw string) map[string]struct{} {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	exts := make(map[string]struct{})
	for _, part := range strings.Split(raw, ",") {
		ext := strings.ToLower(strings.TrimSpace(part))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		exts[ext] = struct{}{}
	}
	return exts
}
