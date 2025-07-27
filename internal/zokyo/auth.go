package zokyo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthRoute(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(req)
}

func CreateLoginRoute(c *gin.Context) {

}
