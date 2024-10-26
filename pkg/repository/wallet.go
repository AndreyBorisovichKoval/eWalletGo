// C:\GoProject\src\eWalletGo_TestTask\pkg\repository\wallet.go

package repository

import (
	"eWalletGo_TestTask/db"
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/models"
	"errors"

	"gorm.io/gorm"
)

// CheckWalletExists проверяет, существует ли кошелек с данным ID...
func CheckWalletExists(walletID string) (bool, error) {
	var count int64
	err := db.GetDBConn().Model(&models.Wallet{}).Where("wallet_number = ?", walletID).Count(&count).Error
	if err != nil {
		logger.Error.Printf("[repository.CheckWalletExists] Error checking wallet existence: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errs.ErrWalletNotFound
		}
		return false, errs.ErrSomethingWentWrong
	}
	return count > 0, nil
}

// // UpdateWalletBalance обновляет баланс счета по WalletID...
// func UpdateWalletBalance(walletID string, amount float64) error {
// 	var account models.Account

// 	// Сначала находим счет, связанный с данным WalletID
// 	if err := db.GetDBConn().Table("accounts").
// 		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
// 		Where("wallets.wallet_number = ?", walletID).
// 		Select("accounts.id").
// 		First(&account).Error; err != nil {
// 		logger.Error.Printf("[repository.UpdateWalletBalance] Error finding account for wallet ID %s: %v", walletID, err)
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return errs.ErrWalletNotFound
// 		}
// 		return errs.ErrSomethingWentWrong
// 	}

// 	// Теперь обновляем баланс счета
// 	result := db.GetDBConn().Model(&models.Account{}).
// 		Where("id = ?", account.ID).
// 		Update("balance", gorm.Expr("balance + ?", amount))

// 	if result.Error != nil {
// 		logger.Error.Printf("[repository.UpdateWalletBalance] Error updating account balance for wallet ID %s: %v", walletID, result.Error)
// 		return errs.ErrSomethingWentWrong
// 	}

// 	// Проверка, что обновление затронуло строки
// 	if result.RowsAffected == 0 {
// 		logger.Warning.Printf("[repository.UpdateWalletBalance] No rows affected, possibly account not found for wallet ID: %s", walletID)
// 		return errs.ErrWalletNotFound
// 	}

// 	logger.Info.Printf("[repository.UpdateWalletBalance] Account balance updated successfully for wallet ID: %s", walletID)
// 	return nil
// }

// // CreateTransaction создает запись транзакции для счета (account)...
// func CreateTransaction(accountID uint, amount float64, transactionType string) error {
// 	transaction := models.Transaction{
// 		AccountID: accountID,
// 		Amount:    amount,
// 		Type:      transactionType,
// 	}
// 	if err := db.GetDBConn().Create(&transaction).Error; err != nil {
// 		logger.Error.Printf("[repository.CreateTransaction] Error recording transaction: %v", err)
// 		return errs.ErrSomethingWentWrong
// 	}
// 	return nil
// }

// // Получает ID счета, связанный с `wallet_number` через `wallet_id`
// func GetAccountIDByWalletID(walletNumber string) (uint, error) {
// 	var wallet models.Wallet
// 	err := db.GetDBConn().Select("id").Where("wallet_number = ?", walletNumber).First(&wallet).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			logger.Warning.Printf("[repository.GetAccountIDByWalletID] Wallet not found for wallet number: %s", walletNumber)
// 			return 0, errs.ErrWalletNotFound
// 		}
// 		logger.Error.Printf("[repository.GetAccountIDByWalletID] Error retrieving wallet ID for wallet number %s: %v", walletNumber, err)
// 		return 0, errs.ErrSomethingWentWrong
// 	}

// 	// Теперь используем `wallet.ID` для поиска `account` в таблице `accounts`
// 	var account models.Account
// 	err = db.GetDBConn().Select("id").Where("wallet_id = ?", wallet.ID).First(&account).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			logger.Warning.Printf("[repository.GetAccountIDByWalletID] Account not found for wallet ID: %d", wallet.ID)
// 			return 0, errs.ErrRecordNotFound
// 		}
// 		logger.Error.Printf("[repository.GetAccountIDByWalletID] Error retrieving account ID for wallet ID %d: %v", wallet.ID, err)
// 		return 0, errs.ErrSomethingWentWrong
// 	}
// 	return account.ID, nil
// }

// CheckWalletExistsTx проверяет, существует ли кошелек с данным ID в рамках транзакции...
func CheckWalletExistsTx(walletID string, tx *gorm.DB) (bool, error) {
	var count int64
	err := tx.Model(&models.Wallet{}).Where("wallet_number = ?", walletID).Count(&count).Error
	if err != nil {
		logger.Error.Printf("[repository.CheckWalletExistsTx] Error checking wallet existence: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errs.ErrRecordNotFound
		}
		return false, errs.ErrSomethingWentWrong
	}
	return count > 0, nil
}

// // UpdateWalletBalanceTx обновляет баланс кошелька в рамках транзакции...
// func UpdateWalletBalanceTx(walletID string, amount float64, tx *gorm.DB) error {
// 	err := tx.Model(&models.Account{}).
// 		Where("wallet_id = ?", walletID).
// 		Update("balance", gorm.Expr("balance + ?", amount)).Error
// 	if err != nil {
// 		logger.Error.Printf("[repository.UpdateWalletBalanceTx] Error updating wallet balance: %v", err)
// 		return errs.ErrSomethingWentWrong
// 	}
// 	return nil
// }

// UpdateWalletBalanceTx обновляет баланс кошелька в рамках транзакции...
func UpdateWalletBalanceTx(walletID string, amount float64, tx *gorm.DB) error {
	var account models.Account

	// Находим счет, связанный с данным walletID
	err := tx.Table("accounts").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Where("wallets.wallet_number = ?", walletID).
		Select("accounts.id").
		First(&account).Error
	if err != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceTx] Error finding account for wallet ID %s: %v", walletID, err)
		return errs.ErrAccountNotFound
	}

	// Обновляем баланс счета
	result := tx.Model(&models.Account{}).
		Where("id = ?", account.ID).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceTx] Error updating balance for account ID %d: %v", account.ID, result.Error)
		return errs.ErrSomethingWentWrong
	}

	// Проверяем, что обновление затронуло строки
	if result.RowsAffected == 0 {
		logger.Warning.Printf("[repository.UpdateWalletBalanceTx] No rows affected for account ID: %d", account.ID)
		return errs.ErrAccountNotFound
	}

	logger.Info.Printf("[repository.UpdateWalletBalanceTx] Account balance updated successfully for account ID: %d", account.ID)
	return nil
}

// CreateTransactionTx создает запись транзакции для счета (account) в рамках транзакции...
func CreateTransactionTx(accountID uint, amount float64, transactionType string, tx *gorm.DB) error {
	transaction := models.Transaction{
		AccountID: accountID,
		Amount:    amount,
		Type:      transactionType,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		logger.Error.Printf("[repository.CreateTransactionTx] Error recording transaction: %v", err)
		return errs.ErrSomethingWentWrong
	}
	return nil
}

// GetAccountIDByWalletIDTx возвращает AccountID, связанный с данным WalletID в рамках транзакции...
func GetAccountIDByWalletIDTx(walletID string, tx *gorm.DB) (uint, error) {
	var wallet models.Wallet
	err := tx.Select("id").Where("wallet_number = ?", walletID).First(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warning.Printf("[repository.GetAccountIDByWalletIDTx] Wallet not found for wallet number: %s", walletID)
			return 0, errs.ErrWalletNotFound
		}
		logger.Error.Printf("[repository.GetAccountIDByWalletIDTx] Error retrieving wallet ID for wallet number %s: %v", walletID, err)
		return 0, errs.ErrSomethingWentWrong
	}

	var account models.Account
	err = tx.Select("id").Where("wallet_id = ?", wallet.ID).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warning.Printf("[repository.GetAccountIDByWalletIDTx] Account not found for wallet ID: %d", wallet.ID)
			return 0, errs.ErrAccountNotFound
		}
		logger.Error.Printf("[repository.GetAccountIDByWalletIDTx] Error retrieving account ID for wallet ID %d: %v", wallet.ID, err)
		return 0, errs.ErrSomethingWentWrong
	}
	return account.ID, nil
}
