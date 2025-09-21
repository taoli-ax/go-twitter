package main

import (
	"log"
	"net/http"
)

func main() {
	// 创建一个新的服务器实例，这里会包含我们的路由和依赖
	server, err := NewServer()
	if err != nil {
		log.Fatalf("could not create server: %v", err)
	}

	// 启动 HTTP 服务器
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server.router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
