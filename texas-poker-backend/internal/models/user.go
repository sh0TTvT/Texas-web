// 用户数据模型
// 作用：定义用户相关的数据结构和数据库操作方法

package models

import (
	"database/sql"
	"time"
)

// User 用户模型
type User struct {
	ID         int64     `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"-" db:"password_hash"` // 密码不返回给客户端
	Chips      int       `json:"chips" db:"chips"`
	TotalGames int       `json:"total_games" db:"total_games"`
	TotalWins  int       `json:"total_wins" db:"total_wins"`
	AvatarURL  string    `json:"avatar_url" db:"avatar_url"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新个人信息请求结构
type UpdateProfileRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

// UserResponse 用户响应结构（不包含敏感信息）
type UserResponse struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Chips      int       `json:"chips"`
	TotalGames int       `json:"total_games"`
	TotalWins  int       `json:"total_wins"`
	AvatarURL  string    `json:"avatar_url"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToResponse 将User转换为UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Chips:      u.Chips,
		TotalGames: u.TotalGames,
		TotalWins:  u.TotalWins,
		AvatarURL:  u.AvatarURL,
		Status:     u.Status,
		CreatedAt:  u.CreatedAt,
	}
}

// CreateUser 创建用户
func CreateUser(db *sql.DB, username, email, passwordHash string) (*User, error) {
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES (?, ?, ?)
	`
	result, err := db.Exec(query, username, email, passwordHash)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(db, id)
}

// GetUserByID 根据ID获取用户
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, chips, total_games, 
		       total_wins, avatar_url, status, created_at, updated_at
		FROM users WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Chips, &user.TotalGames, &user.TotalWins,
		&user.AvatarURL, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, chips, total_games, 
		       total_wins, avatar_url, status, created_at, updated_at
		FROM users WHERE username = ?
	`
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Chips, &user.TotalGames, &user.TotalWins,
		&user.AvatarURL, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, chips, total_games, 
		       total_wins, avatar_url, status, created_at, updated_at
		FROM users WHERE email = ?
	`
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Chips, &user.TotalGames, &user.TotalWins,
		&user.AvatarURL, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(db *sql.DB, id int64, email, avatarURL string) error {
	query := `
		UPDATE users SET email = ?, avatar_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, email, avatarURL, id)
	return err
}

// UpdateUserChips 更新用户筹码
func UpdateUserChips(db *sql.DB, id int64, chips int) error {
	query := `
		UPDATE users SET chips = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, chips, id)
	return err
} 