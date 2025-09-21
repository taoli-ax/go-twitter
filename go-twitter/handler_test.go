package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler(t *testing.T) {
	// 子测试用例的结构
	testCases := []struct {
		name               string
		payload            map[string]string
		initialStoreState  map[string]*User // 模拟测试前数据库里已有的数据
		expectedStatusCode int
	}{
		{
			name:               "Success",
			payload:            map[string]string{"username": "testuser", "password": "password123"},
			initialStoreState:  nil, // 数据库是空的
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:    "User already exists",
			payload: map[string]string{"username": "existinguser", "password": "password123"},
			initialStoreState: map[string]*User{
				"existinguser": {ID: 1, Username: "existinguser"},
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "Invalid payload",
			payload:            nil, // 我们会发送一个无效的 JSON
			initialStoreState:  nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	// 遍历所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. 准备工作：为每个测试创建一个全新的、隔离的环境
			store := NewInMemoryStore()
			if tc.initialStoreState != nil {
				store.users = tc.initialStoreState
			}

			server := &Server{router: mux.NewRouter(), store: store}
			server.routes() // 注册路由

			// 2. 创建请求
			var reqBody []byte
			var err error
			if tc.payload != nil {
				reqBody, err = json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("could not marshal payload: %v", err)
				}
			} else {
				reqBody = []byte(`{"username":"test"`) // 无效的 JSON
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))

			// 3. 创建一个 ResponseRecorder 来捕获响应
			rr := httptest.NewRecorder()

			// 4. 执行请求
			server.ServeHTTP(rr, req)

			// 5. 断言：检查状态码是否符合预期
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatusCode)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	// 准备一个预先注册好的用户，密码是 "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	preexistingUser := &User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	testCases := []struct {
		name               string
		payload            map[string]string
		expectedStatusCode int
	}{
		{
			name:               "Success",
			payload:            map[string]string{"username": "testuser", "password": "password123"},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "User not found",
			payload:            map[string]string{"username": "nonexistentuser", "password": "password123"},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Incorrect password",
			payload:            map[string]string{"username": "testuser", "password": "wrongpassword"},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Invalid payload",
			payload:            nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. 准备工作
			store := NewInMemoryStore()
			store.users["testuser"] = preexistingUser // 将预注册用户放入 store

			server := &Server{router: mux.NewRouter(), store: store}
			server.routes()

			// 2. 创建请求
			var reqBody []byte
			var err error
			if tc.payload != nil {
				reqBody, err = json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("could not marshal payload: %v", err)
				}
			} else {
				reqBody = []byte(`{"username":"test"`)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))

			// 3. 创建 ResponseRecorder
			rr := httptest.NewRecorder()

			// 4. 执行请求
			server.ServeHTTP(rr, req)

			// 5. 断言
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatusCode)
			}

			// (可选) 对成功案例，可以进一步断言响应体中是否包含 token
			if tc.expectedStatusCode == http.StatusOK {
				var resp map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
				if _, ok := resp["token"]; !ok {
					t.Errorf("response body does not contain token")
				}
			}
		})
	}
}
