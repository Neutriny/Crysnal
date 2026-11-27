// Package server 提供面板 HTTP 服务：内嵌前端 SPA（web/dist）并预留 /api 前缀。
package server

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strings"

	"crysnal/internal/system"
)

//go:embed all:dist
var distFS embed.FS

// New 构建面板 HTTP Handler。addr 形如 ":8080"。
func New(addr string) (http.Handler, error) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("定位内嵌前端资源: %w", err)
	}
	// 构建产物和字体不会改变内容，启用强缓存
	immutable := http.FileServer(http.FS(sub))

	mux := http.NewServeMux()
	sys := system.NewHandler()
	mux.HandleFunc("POST /api/system/browse", sys.BrowseFile)
	mux.HandleFunc("GET /api/system/fs/list", sys.ListDir)
	mux.Handle("/api/", http.NotFoundHandler()) // 其余 API 尚未实现，先返回 404 而不是回退到 index.html
	mux.Handle("/assets/", immutable)
	mux.Handle("/fonts/", immutable)
	mux.Handle("/", spaHandler(sub))
	return mux, nil
}

// Run 启动服务并阻塞，直到进程被终止。
func Run(addr string) error {
	handler, err := New(addr)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听 %s: %w", addr, err)
	}

	slog.Info("Crysnal 面板已启动", "addr", listener.Addr().String())
	if err := http.Serve(listener, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP 服务退出: %w", err)
	}
	return nil
}

// spaHandler 处理静态文件与 SPA 路由回退：存在的文件直接返回，否则回退到 index.html，
// 使前端 history 路由（如 /instance/create）刷新后仍能加载。
func spaHandler(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 规范化路径，避免目录穿越与重复斜杠
		p := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
		if p != "/" {
			if f, err := fsys.Open(strings.TrimPrefix(p, "/")); err == nil {
				_ = f.Close()
				r.URL.Path = p
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
