package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth/handler"
	handler2 "github.com/hardikm9850/GoChat/internal/chat/handler"
	http "github.com/hardikm9850/GoChat/internal/contacts/handler"
	"github.com/hardikm9850/authkit/jwt"
	"github.com/hardikm9850/authkit/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(
	r *gin.Engine,
	jwtManager *jwt.HS256Manager,
	wsHandler *handler2.WSHandler,
	authHandler *handler.AuthHandler,
	contactsHandler *http.ContactsHandler,
	conversationHandler *handler2.ConversationHandler,
	messagesHandler *handler2.MessagesHandler,
) {

	// Group creates a new router group. We should add all the routes that have common middlewares or the same path prefix.
	// For example, all the routes that use a common middleware for authorization could be grouped.

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GoChat backend running",
			"version": "v1.0",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	// swag init -g cmd/main.go --output ./docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// -------- AUTH --------
	auth := r.Group("/auth")
	{
		v1 := auth.Group("/v1")
		{
			v1.POST("/register", authHandler.Register)
			v1.POST("/login", authHandler.Login)
			v1.POST("/refresh", authHandler.Refresh)
		}
	}

	// -------- WEBSOCKET --------
	ws := r.Group("/ws")
	ws.Use(middleware.JWTAuth(*jwtManager))
	{
		ws.GET("/chat", wsHandler.HandleWebSocket)
	}

	// -------- Main API (protected)--------
	v1Auth := r.Group("/v1")
	v1Auth.Use(middleware.JWTAuth(*jwtManager))
	{
		v1Auth.GET("/me", func(c *gin.Context) {
			userID := c.GetString("userID")
			c.JSON(200, gin.H{"user_id": userID})
		})

		v1Auth.POST("/contacts/sync", contactsHandler.SyncContacts)

		conversations := v1Auth.Group("/conversations")
		{
			conversations.POST("", conversationHandler.CreateConversation)
			conversations.GET("", conversationHandler.GetMyConversations)
			conversations.GET("/:id", conversationHandler.GetConversation)
			conversations.GET("/:id/messages", messagesHandler.GetMessages)
		}
	}
}
