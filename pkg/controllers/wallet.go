// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\wallet.go

package controllers

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CheckWalletExistence проверяет существование кошелька и возвращает ответ клиенту...
func CheckWalletExistence(c *gin.Context) {
	walletID := c.Param("wallet_id")
	exists, err := service.CheckWalletExistence(walletID)
	if err != nil {
		handleError(c, err)
		return
	}

	if !exists {
		handleError(c, errs.ErrWalletNotFound)
		return
	}

	logger.Info.Printf("[controllers.CheckWalletExistence] Wallet found: %s", walletID)
	c.JSON(http.StatusOK, gin.H{"message": "wallet found"})
}

// RechargeWallet пополняет кошелек и возвращает ответ клиенту...
func RechargeWallet(c *gin.Context) {
	var req struct {
		WalletID string  `json:"wallet_id"`
		Amount   float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errs.ErrInvalidRequest)
		return
	}

	err := service.RechargeWallet(req.WalletID, req.Amount)
	if err != nil {
		handleError(c, err)
		return
	}

	logger.Info.Printf("[controllers.RechargeWallet] Wallet recharged successfully: %s, Amount: %.2f", req.WalletID, req.Amount)
	c.JSON(http.StatusOK, gin.H{"message": "Wallet recharged successfully"})
}

// GetMonthlyRechargeSummary возвращает количество и сумму пополнений за указанный месяц по `walletID`
func GetMonthlyRechargeSummary(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	walletID := c.Query("wallet_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		handleError(c, errs.ErrInvalidRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		handleError(c, errs.ErrInvalidRequest)
		return
	}

	totalCount, totalAmount, err := service.GetMonthlyRechargeSummary(walletID, year, month)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_count":  totalCount,
		"total_amount": totalAmount,
	})
}
