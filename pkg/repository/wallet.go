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

// GetMonthlyRechargeSummary возвращает суммарные данные по `walletID` для указанного месяца и года
func GetMonthlyRechargeSummary(walletID string, year int, month int) (int64, float64, error) {
	// Проверяем существование кошелька
	var walletCount int64
	err := db.GetDBConn().
		Model(&models.Wallet{}).
		Where("wallet_number = ?", walletID).
		Count(&walletCount).Error
	if err != nil || walletCount == 0 {
		logger.Warning.Printf("[repository.GetMonthlyRechargeSummary] Wallet not found for wallet number: %s", walletID)
		return 0, 0, errs.ErrWalletNotFound
	}

	var totalCount int64
	var totalAmount float64

	// Запрос для суммирования пополнений за указанный месяц и год
	err = db.GetDBConn().
		Model(&models.Transaction{}).
		Joins("JOIN accounts ON transactions.account_id = accounts.id").
		Joins("JOIN wallets ON accounts.wallet_id = wallets.id").
		Where("wallets.wallet_number = ? AND EXTRACT(YEAR FROM transactions.created_at) = ? AND EXTRACT(MONTH FROM transactions.created_at) = ? AND transactions.type = ?", walletID, year, month, "recharge").
		Count(&totalCount).
		Select("COALESCE(SUM(transactions.amount), 0)").Row().Scan(&totalAmount)

	if err != nil {
		logger.Error.Printf("[repository.GetMonthlyRechargeSummary] Ошибка получения данных: %v", err)
		return 0, 0, errs.ErrSomethingWentWrong
	}

	return totalCount, totalAmount, nil
}

// GetWalletBalance возвращает текущий баланс кошелька по его ID...
func GetWalletBalance(walletID string) (float64, error) {
	var account models.Account
	err := db.GetDBConn().Table("accounts").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Where("wallets.wallet_number = ?", walletID).
		Select("accounts.balance").
		First(&account).Error

	// Применяем translateError для конвертации ошибки
	if err != nil {
		err = translateError(err)
		if errors.Is(err, errs.ErrRecordNotFound) {
			logger.Warning.Printf("[repository.GetWalletBalance] Wallet not found for wallet ID: %s", walletID)
			return 0, errs.ErrWalletNotFound
		}
		logger.Error.Printf("[repository.GetWalletBalance] Error retrieving balance for wallet ID %s: %v", walletID, err)
		return 0, errs.ErrSomethingWentWrong
	}
	return account.Balance, nil
}

// CalculateBalanceFromTransactions вычисляет новый баланс на основе всех транзакций...
func CalculateBalanceFromTransactions(walletID string) (float64, error) {
	var newBalance float64

	// Суммируем все транзакции по данному кошельку
	err := db.GetDBConn().
		Table("transactions").
		Joins("JOIN accounts ON transactions.account_id = accounts.id").
		Joins("JOIN wallets ON accounts.wallet_id = wallets.id").
		Where("wallets.wallet_number = ?", walletID).
		Select("COALESCE(SUM(transactions.amount), 0)").Scan(&newBalance).Error

	if err != nil {
		logger.Error.Printf("[repository.CalculateBalanceFromTransactions] Error calculating balance for wallet ID %s: %v", walletID, err)
		return 0, errs.ErrSomethingWentWrong
	}

	return newBalance, nil
}

// UpdateWalletBalanceDirectly обновляет баланс кошелька напрямую...
func UpdateWalletBalanceDirectly(walletID string, newBalance float64) error {
	var account models.Account

	// Находим счет, связанный с данным walletID
	err := db.GetDBConn().Table("accounts").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Where("wallets.wallet_number = ?", walletID).
		Select("accounts.id").
		First(&account).Error

	if err != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceDirectly] Error finding account for wallet ID %s: %v", walletID, err)
		return errs.ErrAccountNotFound
	}

	// Обновляем баланс счета
	result := db.GetDBConn().Model(&models.Account{}).
		Where("id = ?", account.ID).
		Update("balance", newBalance)

	if result.Error != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceDirectly] Error updating balance for account ID %d: %v", account.ID, result.Error)
		return errs.ErrSomethingWentWrong
	}

	// Проверка, что обновление затронуло строки
	if result.RowsAffected == 0 {
		logger.Warning.Printf("[repository.UpdateWalletBalanceDirectly] No rows affected for account ID: %d", account.ID)
		return errs.ErrAccountNotFound
	}

	logger.Info.Printf("[repository.UpdateWalletBalanceDirectly] Account balance updated successfully for account ID: %d", account.ID)
	return nil
}
