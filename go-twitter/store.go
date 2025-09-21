package main

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// Store 定义了我们应用所需的所有数据存储方法
type Store interface {
	// CreateUser 用户相关
	CreateUser(username, password string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

// User 代表一个用户对象
type User struct {
	ID           int
	Username     string
	PasswordHash string
}

// InMemoryStore 是 Store 接口的一个内存实现
type InMemoryStore struct {
	mu         sync.RWMutex
	users      map[string]*User
	nextUserID int
}

// NewInMemoryStore 创建一个新的 InMemoryStore 实例
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users:      make(map[string]*User),
		nextUserID: 1,
	}
}

// CreateUser 创建一个新用户并将其存入内存
func (s *InMemoryStore) CreateUser(username, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查用户名是否已存在
	if _, exists := s.users[username]; exists {
		return nil, fmt.Errorf("username '%s' already exists", username)
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %v", err)
	}

	// 创建新用户
	user := &User{
		ID:           s.nextUserID,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	// 存储用户
	s.users[username] = user
	s.nextUserID++

	return user, nil
}

// GetUserByUsername 通过用户名从内存中获取用户
func (s *InMemoryStore) GetUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, fmt.Errorf("user '%s' not found", username)
	}

	return user, nil
}
