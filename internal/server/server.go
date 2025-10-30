package server

import (
	"MusicDev33/mdapi3/internal/zokyo"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func TestRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "okay",
	})
}

func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Test route
	r.GET("/test", TestRoute)

	// Zokyo routes group
	zokyoGroup := r.Group("/zokyo")
	{
		zokyoGroup.POST("/code", zokyo.CreateNewChatRoute)
		zokyoGroup.GET("/convs/:username", zokyo.GetConversationsRoute)
		zokyoGroup.GET("/msgs/:convId", zokyo.GetChatsByConvIdRoute)
		zokyoGroup.GET("/verify/:username", zokyo.CheckWhitelistUserRoute)
		zokyoGroup.POST("/auth", zokyo.AuthRoute)
		zokyoGroup.POST("/login/create", zokyo.CreateLoginRoute)
		zokyoGroup.DELETE("/convs/:convId", zokyo.DeleteConversationByIdRoute)
	}

	return &Server{
		router: r,
	}
}

func (s *Server) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}
