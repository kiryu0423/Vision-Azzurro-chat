package main

import (
	"chat-app/internal/handler"
	"chat-app/internal/repository"
	"chat-app/internal/router"
	"chat-app/internal/service"
	"os"

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
	notifyHandler := handler.NewNotifyWSHandler()
	wsHandler := handler.NewWebSocketHandler(msgRepo, roomService, notifyHandler)
	wsNotifyHandler := handler.NewNotifyWSHandler()

	r := router.SetupRouter(userHandler, authHandler, roomHandler, msgHandler, wsHandler, wsNotifyHandler)

	r.Run(":" + os.Getenv("PORT"))
}
