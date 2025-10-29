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

	// Zokyo routes
	r.POST("/code", zokyo.CreateNewChatRoute)
	r.GET("/convs/:username", zokyo.GetConversationsRoute)
	r.GET("/msgs/:convId", zokyo.GetChatsByConvIdRoute)
	r.GET("/verify/:username", zokyo.CheckWhitelistUserRoute)
	r.POST("/auth", zokyo.AuthRoute)
	r.POST("/login/create", zokyo.CreateLoginRoute)
	r.DELETE("/convs/:convId", zokyo.DeleteConversationByIdRoute)

	return &Server{
		router: r,
	}
}

func (s *Server) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}
