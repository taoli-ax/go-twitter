package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server 持有我们应用的所有依赖项
type Server struct {
	router *mux.Router
	store  Store
}

// NewServer 是我们的构造函数，用于创建和初始化 Server
func NewServer() (*Server, error) {
	// 初始化一个新的内存存储
	store := NewInMemoryStore()

	// 创建一个新的 Server 实例
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
	}

	// 初始化路由
	s.routes()

	return s, nil
}

// routes 方法用于注册我们所有的 API 端点
func (s *Server) routes() {
	// Commit 1: 注册和登录路由
	s.router.HandleFunc("/register", s.handleRegister()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")

	// Commit 3: 推文相关的路由会在这里添加
}

// ServeHTTP 使得我们的 Server 类型满足 http.Handler 接口
// 这样我们就可以将它直接传递给 http.ListenAndServe
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
