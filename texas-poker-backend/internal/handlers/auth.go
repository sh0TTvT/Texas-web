// 用户认证处理器
// 作用：处理用户注册、登录、登出、个人信息等HTTP请求

package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"texas-poker-backend/internal/config"
	"texas-poker-backend/internal/models"
	"texas-poker-backend/internal/utils"
	"texas-poker-backend/internal/websocket"
)

// Handler HTTP处理器结构
type Handler struct {
	db        *sql.DB
	redis     *redis.Client
	wsManager *websocket.Manager
	config    *config.Config
}

// New 创建新的处理器实例
func New(db *sql.DB, redis *redis.Client, wsManager *websocket.Manager) *Handler {
	return &Handler{
		db:        db,
		redis:     redis,
		wsManager: wsManager,
		config:    config.Load(),
	}
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	if _, err := models.GetUserByUsername(h.db, req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if _, err := models.GetUserByEmail(h.db, req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "邮箱已被注册",
		})
		return
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败",
		})
		return
	}

	// 创建用户
	user, err := models.CreateUser(h.db, req.Username, req.Email, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户创建失败",
			"details": err.Error(),
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, "user", h.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token生成失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 获取用户信息（支持用户名或邮箱登录）
	var user *models.User
	var err error

	// 先尝试用户名登录
	user, err = models.GetUserByUsername(h.db, req.Username)
	if err != nil {
		// 再尝试邮箱登录
		user, err = models.GetUserByEmail(h.db, req.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户名或密码错误",
			})
			return
		}
	}

	// 检查账号状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "账号已被禁用",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, "user", h.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// Logout 用户登出
func (h *Handler) Logout(c *gin.Context) {
	// 在简单实现中，JWT是无状态的，客户端删除token即可
	// 在生产环境中，可以将token加入黑名单
	
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// GetProfile 获取用户个人信息
func (h *Handler) GetProfile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
		})
		return
	}

	// 获取用户信息
	user, err := models.GetUserByID(h.db, userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}

// UpdateProfile 更新用户个人信息
func (h *Handler) UpdateProfile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 检查邮箱是否被其他用户使用
	if req.Email != "" {
		if existingUser, err := models.GetUserByEmail(h.db, req.Email); err == nil {
			if existingUser.ID != userID.(int64) {
				c.JSON(http.StatusConflict, gin.H{
					"error": "邮箱已被其他用户使用",
				})
				return
			}
		}
	}

	// 更新用户信息
	err := models.UpdateUser(h.db, userID.(int64), req.Email, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新失败",
			"details": err.Error(),
		})
		return
	}

	// 获取更新后的用户信息
	user, err := models.GetUserByID(h.db, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user":    user.ToResponse(),
	})
} 