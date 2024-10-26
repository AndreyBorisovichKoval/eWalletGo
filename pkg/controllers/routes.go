// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\routes.go

package controllers

import (
	"eWalletGo_TestTask/configs"
	"eWalletGo_TestTask/logger"

	_ "eWalletGo_TestTask/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// PingPong обрабатывает запрос ping и отвечает сообщением pong
func PingPong(c *gin.Context) {
	logger.Info.Printf("Маршрут '%s' вызван методом '%s'", c.FullPath(), c.Request.Method)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// InitRoutes инициализирует все маршруты с необходимыми middleware и Swagger
func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Устанавливаем режим работы Gin из конфигурации
	gin.SetMode(configs.AppSettings.AppParams.GinMode)

	// Маршрут для документации Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Маршрут для проверки состояния сервера
	router.GET("/ping", PingPong)

	// Группа маршрутов для управления кошельком
	// walletG := router.Group("/wallet", middleware.AuthMiddleware)
	walletG := router.Group("/wallet")
	{
		walletG.GET("/check/:wallet_id", CheckWalletExistence) // Проверка существования кошелька...
		// walletG.POST("/recharge", RechargeWallet)                                  // Пополнение кошелька...
		// walletG.GET("/monthly-summary", GetMonthlyRechargeSummary)                 // Суммарные данные пополнений за текущий месяц, либо за указанный период (нужны/необходимы qwery-параметры года и месяца)...
		// walletG.GET("/balance/:wallet_id", GetWalletBalance)                       // Получение текущего баланса кошелька...
		// walletG.PATCH("/recalculate-balance/:wallet_id", RecalculateWalletBalance) // Перерасчёт баланса на основе транзакций...
	}

	// Маршрут для вставки тестовых данных
	router.POST("/insert-test-data", InsertTestData)

	return router
}
