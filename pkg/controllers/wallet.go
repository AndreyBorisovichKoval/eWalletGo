// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\wallet.go

package controllers

import (
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckWalletExistence проверяет существование кошелька и возвращает ответ клиенту...
func CheckWalletExistence(c *gin.Context) {
	walletID := c.Param("wallet_id")
	err := service.CheckWalletExistence(walletID)
	if err != nil {
		handleError(c, err)
		return
	}

	logger.Info.Printf("[controllers.CheckWalletExistence] Wallet found: %s", walletID)
	c.JSON(http.StatusOK, gin.H{"message": "wallet found"})
}
