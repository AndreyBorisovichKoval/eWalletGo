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

// CheckWalletExistence checks if the wallet exists
// @Summary Check wallet existence
// @Tags wallets
// @Description Checks if the wallet exists based on wallet ID
// @ID check-wallet-existence
// @Param wallet_id path string true "Wallet ID"
// @Success 200 {object} map[string]interface{} "wallet found"
// @Failure 404 {object} ErrorResponse "Wallet not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /wallet/{wallet_id}/exists [get]
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

// GetMonthlyRechargeSummary returns the total and count of recharges for a specified month
// @Summary Get monthly recharge summary
// @Tags wallets
// @Description Retrieves the total amount and count of recharge operations for a wallet in a specified month
// @ID get-monthly-recharge-summary
// @Param wallet_id query string true "Wallet ID"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} map[string]interface{} "Recharge summary retrieved"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 404 {object} ErrorResponse "Wallet not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /wallet/monthly-summary [get]
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

// GetWalletBalance retrieves the balance of a wallet
// @Summary Get wallet balance
// @Tags wallets
// @Description Retrieves the current balance of a specified wallet
// @ID get-wallet-balance
// @Param wallet_id path string true "Wallet ID"
// @Success 200 {object} map[string]interface{} "Wallet balance retrieved"
// @Failure 404 {object} ErrorResponse "Wallet not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /wallet/{wallet_id}/balance [get]
func GetWalletBalance(c *gin.Context) {
	walletID := c.Param("wallet_id")
	balance, err := service.GetWalletBalance(walletID)
	if err != nil {
		handleError(c, err)
		return
	}

	logger.Info.Printf("[controllers.GetWalletBalance] Balance retrieved for wallet: %s", walletID)
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// RecalculateWalletBalance recalculates the wallet balance based on transactions
// @Summary Recalculate wallet balance
// @Tags wallets
// @Description Recalculates the wallet balance based on transaction records
// @ID recalculate-wallet-balance
// @Param wallet_id path string true "Wallet ID"
// @Success 200 {object} map[string]interface{} "Wallet balance recalculated"
// @Failure 404 {object} ErrorResponse "Wallet not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /wallet/{wallet_id}/recalculate-balance [patch]
func RecalculateWalletBalance(c *gin.Context) {
	walletID := c.Param("wallet_id")
	newBalance, err := service.RecalculateWalletBalance(walletID)
	if err != nil {
		handleError(c, err)
		return
	}

	logger.Info.Printf("[controllers.RecalculateWalletBalance] Balance recalculated for wallet: %s, New Balance: %.2f", walletID, newBalance)
	c.JSON(http.StatusOK, gin.H{"wallet_id": walletID, "new_balance": newBalance})
}

// RechargeWallet recharges the wallet
// @Summary Recharge wallet
// @Tags wallets
// @Description Adds a specified amount to the wallet balance
// @ID recharge-wallet
// @Accept json
// @Produce json
// @Param wallet_id body string true "Wallet ID"
// @Param amount body float64 true "Amount to recharge"
// @Success 200 {object} map[string]interface{} "Wallet recharged successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 404 {object} ErrorResponse "Wallet not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /wallet/recharge [post]
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
