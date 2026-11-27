package main

import (
	"flag"
	"log/slog"
	"os"

	"crysnal/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "面板监听地址")
	flag.Parse()

	if err := server.Run(*addr); err != nil {
		slog.Error("启动失败", "err", err)
		os.Exit(1)
	}
}
