package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// handleRegister 返回一个处理用户注册请求的 http.HandlerFunc
func (s *Server) handleRegister() http.HandlerFunc {
	// 定义请求体结构
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// 解码请求体
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 调用数据存储层来创建用户
		_, err := s.store.CreateUser(req.Username, req.Password)
		if err != nil {
			// 这里我们假设 store 层会返回一个具体的错误信息
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		// 返回成功响应
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
	}
}

// handleLogin 返回一个处理用户登录请求的 http.HandlerFunc
func (s *Server) handleLogin() http.HandlerFunc {
	// 定义请求体结构
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// 定义响应体结构
	type response struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// 解码请求体
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 从数据存储层获取用户
		user, err := s.store.GetUserByUsername(req.Username)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// 比较哈希密码和请求的密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// 登录成功，返回一个虚拟的 token
		// 在真实的课程中，这里会生成一个 JWT
		resp := response{Token: "fake-jwt-token"}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
