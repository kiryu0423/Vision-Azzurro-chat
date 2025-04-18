package router

import (
	"chat-app/internal/handler"
	"chat-app/internal/middleware"

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

	// ✅ JWTで保護されたルーティンググループ
	auth := r.Group("/", middleware.JWTAuthMiddleware())
	{
		auth.GET("/chat", func(c *gin.Context) {
			c.HTML(200, "chat.html", nil)
		})

		auth.GET("/chat/:room_id", func(c *gin.Context) {
			roomID := c.Param("room_id")
			userName := c.GetString("user_name")

			c.HTML(200, "ws.html", gin.H{
				"RoomID":   roomID,
				"UserName": userName,
			})
		})

		auth.GET("/users", userHandler.ListUsers)
		auth.GET("/me", userHandler.Me)

		auth.GET("/messages/:room_id", msgHandler.GetMessages)

		auth.POST("/rooms", roomHandler.CreateRoom)
		auth.GET("/rooms", roomHandler.ListRooms)
		auth.PUT("/rooms/:room_id/name", roomHandler.UpdateRoomName)
		auth.GET("/rooms/:id/members", roomHandler.GetRoomMembers)

		// 既読管理
		auth.POST("/rooms/:room_id/read", roomHandler.MarkRoomAsRead)
		// グループ退会
		auth.DELETE("/rooms/:room_id/members/me", roomHandler.LeaveRoom)
		// グループ削除
		auth.DELETE("/rooms/:room_id", roomHandler.DeleteRoom)
	}

	// 認証不要
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

	return r
}
