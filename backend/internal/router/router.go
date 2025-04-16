package router

import (
	"chat-app/internal/handler"
	"chat-app/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	roomHandler *handler.RoomHandler,
	msgHandler *handler.MessageHandler,
	wsHandler *handler.WebSocketHandler,
	wsNotifyHandler *handler.NotifyWSHandler,
) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(middleware.CORSMiddleware())

	// Session
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("chat_session", store))

	// Routing
	r.GET("/users", userHandler.ListUsers)
	r.GET("/me", userHandler.Me)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)
	r.GET("/login-page", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/register-page", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	r.GET("/ws", wsHandler.Handle)
	r.GET("/ws-notify", wsNotifyHandler.Handle)

	r.GET("/chat/:room_id", func(c *gin.Context) {
		session := sessions.Default(c)
		userName := session.Get("user_name")
		if userName == nil {
			c.Redirect(302, "/login-page")
			return
		}
		c.HTML(200, "ws.html", gin.H{
			"RoomID":   c.Param("room_id"),
			"UserName": userName,
		})
	})

	r.GET("/messages/:room_id", msgHandler.GetMessages)
	r.POST("/rooms", roomHandler.CreateRoom)
	r.GET("/rooms", roomHandler.ListRooms)
	r.PUT("/rooms/:room_id/name", roomHandler.UpdateRoomName)

	r.GET("/chat", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			c.Redirect(302, "/login-page")
			return
		}
		c.HTML(200, "chat.html", nil)
	})

	// 既読管理
	r.POST("/rooms/:room_id/read", roomHandler.MarkRoomAsRead)

	return r
}
