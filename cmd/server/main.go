package main

import (
	"chat-app/internal/handler"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"chat-app/internal/middleware"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	dsn := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to DB")
	}

	userRepo := &repository.UserRepository{DB: db}
	userHandler := &handler.UserHandler{UserRepository: userRepo}
	authRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	msgRepo := repository.NewMessageRepository(db)
	msgHandler := handler.NewMessageHandler(msgRepo)
	wsHandler := handler.NewWebSocketHandler(msgRepo)

	r := gin.Default()

	// セッションミドルウェアの追加（CookieStore使用）
    store := cookie.NewStore([]byte("super-secret-key"))
    r.Use(sessions.Sessions("chat_session", store))

	r.LoadHTMLGlob("web/templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// ユーザー一覧を表示
	r.GET("/users", userHandler.GetAllUser)
	// ユーザーをIDで取得
	r.GET("/user/:id", userHandler.GetUserByID)

	// ユーザーページの表示
	r.GET("/mypage", middleware.RequireLogin(), func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
	
		c.JSON(http.StatusOK, gin.H{
			"message":  "You are logged in!",
			"user_id":  userID,
		})
	})

	// ユーザー登録
	r.POST("/register", authHandler.Register)
	// ユーザーログイン
	r.POST("/login", authHandler.Login)
	// ユーザーログアウト
	r.POST("/logout", authHandler.Logout)
	// ユーザーログインページ
	r.GET("/login-page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	// ユーザー登録ページ
	r.GET("/register-page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	

	// websocket接続
	r.GET("/ws", wsHandler.Handle)
	// チャット確認ページ
	r.GET("/ws-test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ws.html", nil)
	})

	r.GET("/messages/:room_id", msgHandler.GetMessages)


	r.Run(":" + os.Getenv("PORT"))
}
