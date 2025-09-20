// 德州扑克后端服务 - 程序入口文件
// 作用：初始化服务器，配置路由和中间件，启动HTTP服务器和WebSocket服务

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"texas-poker-backend/internal/config"
	"texas-poker-backend/internal/database"
	"texas-poker-backend/internal/handlers"
	"texas-poker-backend/internal/middleware"
	"texas-poker-backend/internal/websocket"
)

func main() {
	// 加载配置
	cfg := config.Load()
	
	// 初始化数据库连接
	db, err := database.InitMySQL(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// 初始化Redis连接
	redis := database.InitRedis(cfg.RedisURL)
	defer redis.Close()
	
	// 初始化WebSocket管理器
	wsManager := websocket.NewManager()
	
	// 设置Gin模式
	gin.SetMode(cfg.GinMode)
	
	// 创建Gin引擎
	r := gin.New()
	
	// 添加基础中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	
	// 初始化处理器
	h := handlers.New(db, redis, wsManager)
	
	// 配置路由
	setupRoutes(r, h)
	
	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRoutes 配置所有路由
func setupRoutes(r *gin.Engine, h *handlers.Handler) {
	// API路由组
	api := r.Group("/api")
	{
		// 用户认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/logout", h.Logout)
			auth.GET("/profile", middleware.AuthRequired(), h.GetProfile)
			auth.PUT("/profile", middleware.AuthRequired(), h.UpdateProfile)
		}
		
		// 游戏路由
		rooms := api.Group("/rooms", middleware.AuthRequired())
		{
			rooms.GET("", h.GetRooms)
			rooms.POST("", h.CreateRoom)
			rooms.GET("/:id", h.GetRoom)
			rooms.POST("/:id/join", h.JoinRoom)
			rooms.POST("/:id/leave", h.LeaveRoom)
		}
		
		// 管理员路由
		admin := api.Group("/admin")
		{
			adminAuth := admin.Group("/auth")
			{
				adminAuth.POST("/login", h.AdminLogin)
			}
			
			adminAPI := admin.Group("", middleware.AdminRequired())
			{
				adminAPI.GET("/users", h.GetUsers)
				adminAPI.PUT("/users/:id", h.UpdateUser)
				adminAPI.GET("/rooms", h.GetRoomsAdmin)
				adminAPI.GET("/stats", h.GetStats)
			}
		}
	}
	
	// WebSocket路由
	r.GET("/ws", h.WebSocketHandler)
	
	// 静态文件服务（用于前端部署）
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")
} 