// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\routes.go

package controllers

import (
	"eWalletGo_TestTask/configs"
	"eWalletGo_TestTask/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// PingPong handles the ping request and responds with a pong message
func PingPong(c *gin.Context) {
	logger.Info.Printf("Route '%s' called with method '%s'", c.FullPath(), c.Request.Method)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// InitRoutes initializes all routes with necessary middleware and Swagger
func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Set Gin mode from configuration
	gin.SetMode(configs.AppSettings.AppParams.GinMode)

	// Route for Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Route to check server status
	router.GET("/ping", PingPong)

	// Wallet management route group
	walletG := router.Group("/wallet", AuthMiddleware)

	{
		walletG.GET("/check/:wallet_id", CheckWalletExistence)                     // Check wallet existence...
		walletG.POST("/recharge", RechargeWallet)                                  // Recharge wallet...
		walletG.GET("/monthly-summary", GetMonthlyRechargeSummary)                 // Monthly recharge summary for the current month or specified period (requires query parameters for year and month)...
		walletG.GET("/balance/:wallet_id", GetWalletBalance)                       // Get current wallet balance...
		walletG.PATCH("/recalculate-balance/:wallet_id", RecalculateWalletBalance) // Recalculate balance based on transactions...
	}

	// Route for inserting test data
	router.POST("/insert-test-data", InsertTestData)

	return router
}
