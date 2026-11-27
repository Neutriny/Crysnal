package system

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// 文件对话框等待用户操作，超时给足时间
const pickFileTimeout = 10 * time.Minute

// Handler 提供 /api/system/* 接口。
type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

type browseRequest struct {
	Title  string `json:"title"`
	Filter string `json:"filter"`
}

type browseResponse struct {
	// Path 为选中文件的完整路径；用户取消时为空字符串
	Path string `json:"path"`
}

// BrowseFile 处理 POST /api/system/browse：弹出文件选择对话框并返回选中路径。
func (h *Handler) BrowseFile(w http.ResponseWriter, r *http.Request) {
	var req browseRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "请求体不是合法 JSON")
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), pickFileTimeout)
	defer cancel()

	path, err := PickFile(ctx, req.Title, req.Filter)
	if err != nil {
		slog.Warn("文件选择对话框失败", "err", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, browseResponse{Path: path})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
