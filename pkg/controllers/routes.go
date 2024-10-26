// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\routes.go

package controllers

import (
	_ "eWalletGo_TestTask/docs"
	"eWalletGo_TestTask/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingPong handles the ping request and responds with a pong message
func PingPong(c *gin.Context) {
	logger.Info.Printf("Route '%s' called with method '%s'", c.FullPath(), c.Request.Method)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
