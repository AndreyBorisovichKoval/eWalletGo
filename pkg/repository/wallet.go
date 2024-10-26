// C:\GoProject\src\eWalletGo_TestTask\pkg\repository\wallet.go

package repository

import (
	"eWalletGo_TestTask/db"
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/models"
)

// CheckWalletExists проверяет, существует ли кошелек с данным ID...
func CheckWalletExists(walletID string) (bool, error) {
	var count int64
	err := db.GetDBConn().Model(&models.Wallet{}).Where("wallet_number = ?", walletID).Count(&count).Error
	if err != nil {
		logger.Error.Printf("[repository.CheckWalletExists] error checking wallet existence: %v", err)
		return false, translateError(err) // Используем translateError из gorm.go
	}
	return count > 0, nil
}
