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
	// TODO: Add a way to set this to debug if needed
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// TODO: Add some middleware

	r.GET("/test", TestRoute)
	r.POST("/auth", zokyo.AuthRoute)
	r.POST("/login/create", zokyo.CreateLoginRoute)

	return &Server{
		router: r,
	}
}

func (s *Server) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}
