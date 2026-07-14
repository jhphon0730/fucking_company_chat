package web

import (
	"g_kk_ch/internal/chat"
	"g_kk_ch/internal/infrastructure"
	"g_kk_ch/internal/infrastructure/web/websocket"
	"g_kk_ch/internal/middleware"
	"g_kk_ch/internal/user"
)

func (s *server) registerRoutes() error {
	db, err := infrastructure.GetDB()
	if err != nil {
		return err
	}

	/* USER */
	userRepository := user.NewUserRepository(db)
	chatRepository := chat.NewChatRepository(db)

	userService := user.NewUserService(userRepository)
	chatService := chat.NewChatService(chatRepository)

	/* WebSocket */
	hub := websocket.NewHub(chatService)
	go hub.Run()

	userHandler := user.NewUserHandler(userService)
	chatHandler := chat.NewChatHandler(hub, chatService, userService)

	/* CHAT */

	v1 := s.router.Group("/api/v1")
	{
		v1_auth := v1.Group("/auth")
		{
			v1_auth.GET("", userHandler.FindAll)
			v1_auth.POST("/register", userHandler.Register)
			v1_auth.POST("/login", userHandler.Login)
		}

		v1_room := v1.Group("/room", middleware.JWTAuthMiddleware())
		{
			v1_room.GET("", chatHandler.GetUserRooms)
			v1_room.POST("", chatHandler.CreateRoom)
			v1_room.POST("/:room_id/participants", chatHandler.AddParticipants)
			v1_room.PATCH("/:room_id/read", chatHandler.ReadRoom)
		}
	}

	webSocketHandler := websocket.NewWebSocketHandler(hub, userService)
	s.router.GET("/ws", webSocketHandler.ConnectWebSocket)

	return nil
}
