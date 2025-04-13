package main

import (
	"chat-app/internal/handler"
	"chat-app/internal/middleware"
	"chat-app/internal/repository"
	"chat-app/internal/service"
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(db, userRepo)
	userHandler := handler.NewUserHandler(userService)
	authRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, userRepo)
	roomHandler := handler.NewRoomHandler(roomService, userService)
	msgRepo := repository.NewMessageRepository(db)
	msgHandler := handler.NewMessageHandler(msgRepo, roomService)
	wsHandler := handler.NewWebSocketHandler(msgRepo, roomService)

	r := gin.Default()

	// CORSミドルウェアの追加
	r.Use(middleware.CORSMiddleware())

	// セッションミドルウェアの追加（CookieStore使用）
    store := cookie.NewStore([]byte("super-secret-key"))
    r.Use(sessions.Sessions("chat_session", store))

	r.LoadHTMLGlob("web/templates/*")

	r.Static("/static", "./web/static")

	// ユーザー一覧の取得
	r.GET("/users", userHandler.ListUsers)
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

	// ルームに入る
	r.GET("/chat/:room_id", func(c *gin.Context) {
		session := sessions.Default(c)
		userName := session.Get("user_name")
		if userName == nil {
			c.Redirect(http.StatusFound, "/login-page")
			return
		}
	
		roomID := c.Param("room_id")
		c.HTML(http.StatusOK, "ws.html", gin.H{
			"RoomID":   roomID,
			"UserName": userName,
		})
	})

	// ルームIDで履歴取得
	r.GET("/messages/:room_id", msgHandler.GetMessages)

	// ルーム作成
	r.POST("/rooms", roomHandler.CreateRoom)

	// ルーム一覧の取得
	r.GET("/rooms", roomHandler.ListRooms)


	// チャットトップページ
	r.GET("/chat", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			c.Redirect(http.StatusFound, "/login-page")
			return
		}
		c.HTML(http.StatusOK, "chat.html", nil)
	})
	


	r.Run(":" + os.Getenv("PORT"))
}
