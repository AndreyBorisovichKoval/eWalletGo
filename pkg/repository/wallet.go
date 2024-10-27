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

// CheckWalletExists checks if a wallet with the given ID exists...
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

// CheckWalletExistsTx checks if a wallet with the given ID exists within a transaction...
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

// GetAccountIDByWalletIDTx returns the AccountID associated with the given WalletID within a transaction...
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

// GetMonthlyRechargeSummary returns summary data for `walletID` for the specified month and year
func GetMonthlyRechargeSummary(walletID string, year int, month int) (int64, float64, error) {
	// Check if wallet exists
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

	// Query for summing recharges for the specified month and year
	err = db.GetDBConn().
		Model(&models.Transaction{}).
		Joins("JOIN accounts ON transactions.account_id = accounts.id").
		Joins("JOIN wallets ON accounts.wallet_id = wallets.id").
		Where("wallets.wallet_number = ? AND EXTRACT(YEAR FROM transactions.created_at) = ? AND EXTRACT(MONTH FROM transactions.created_at) = ? AND transactions.type = ?", walletID, year, month, "recharge").
		Count(&totalCount).
		Select("COALESCE(SUM(transactions.amount), 0)").Row().Scan(&totalAmount)

	if err != nil {
		logger.Error.Printf("[repository.GetMonthlyRechargeSummary] Error retrieving data: %v", err)
		return 0, 0, errs.ErrSomethingWentWrong
	}

	return totalCount, totalAmount, nil
}

// GetWalletBalance returns the current balance of the wallet by its ID...
func GetWalletBalance(walletID string) (float64, error) {
	var account models.Account
	err := db.GetDBConn().Table("accounts").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Where("wallets.wallet_number = ?", walletID).
		Select("accounts.balance").
		First(&account).Error

	// Apply translateError to convert the error
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

// CalculateBalanceFromTransactions calculates the new balance based on all transactions...
func CalculateBalanceFromTransactions(walletID string) (float64, error) {
	var newBalance float64

	// Sum all transactions for the given wallet
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

// UpdateWalletBalanceDirectly updates the wallet balance directly...
func UpdateWalletBalanceDirectly(walletID string, newBalance float64) error {
	var account models.Account

	// Find the account associated with the given walletID
	err := db.GetDBConn().Table("accounts").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Where("wallets.wallet_number = ?", walletID).
		Select("accounts.id").
		First(&account).Error

	if err != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceDirectly] Error finding account for wallet ID %s: %v", walletID, err)
		return errs.ErrAccountNotFound
	}

	// Update account balance
	result := db.GetDBConn().Model(&models.Account{}).
		Where("id = ?", account.ID).
		Update("balance", newBalance)

	if result.Error != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceDirectly] Error updating balance for account ID %d: %v", account.ID, result.Error)
		return errs.ErrSomethingWentWrong
	}

	// Check if any rows were affected by the update
	if result.RowsAffected == 0 {
		logger.Warning.Printf("[repository.UpdateWalletBalanceDirectly] No rows affected for account ID: %d", account.ID)
		return errs.ErrAccountNotFound
	}

	logger.Info.Printf("[repository.UpdateWalletBalanceDirectly] Account balance updated successfully for account ID: %d", account.ID)
	return nil
}

// GetWalletWithLimit returns wallet data and applies the limit (custom or default) within a transaction...
func GetWalletWithLimit(walletID string, tx *gorm.DB) (models.WalletWithLimit, error) {
	var wallet models.WalletWithLimit
	err := tx.Table("accounts").
		Select("accounts.balance, accounts.id AS account_id, COALESCE(limit_settings.custom_limit, limit_settings.default_limit) AS max_limit").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Joins("JOIN limit_settings ON limit_settings.client_type = wallets.client_type").
		Where("wallets.wallet_number = ?", walletID).
		First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return wallet, errs.ErrWalletNotFound
	} else if err != nil {
		logger.Error.Printf("[repository.GetWalletWithLimit] Error retrieving wallet data: %v", err)
		return wallet, errs.ErrSomethingWentWrong
	}
	return wallet, nil
}

// UpdateWalletBalanceTx updates the wallet balance within a transaction...
func UpdateWalletBalanceTx(walletID string, amount float64, tx *gorm.DB) error {
	var account models.Account

	// Find the account associated with the given walletID
	err := tx.Table("accounts").
		Joins("JOIN wallets ON wallets.id = accounts.wallet_id").
		Where("wallets.wallet_number = ?", walletID).
		Select("accounts.id").
		First(&account).Error
	if err != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceTx] Error finding account for wallet ID %s: %v", walletID, err)
		return errs.ErrAccountNotFound
	}

	// Update account balance
	result := tx.Model(&models.Account{}).
		Where("id = ?", account.ID).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		logger.Error.Printf("[repository.UpdateWalletBalanceTx] Error updating balance for account ID %d: %v", account.ID, result.Error)
		return errs.ErrSomethingWentWrong
	}

	// Check if any rows were affected by the update
	if result.RowsAffected == 0 {
		logger.Warning.Printf("[repository.UpdateWalletBalanceTx] No rows affected for account ID: %d", account.ID)
		return errs.ErrAccountNotFound
	}

	logger.Info.Printf("[repository.UpdateWalletBalanceTx] Account balance updated successfully for account ID: %d", account.ID)
	return nil
}

// CreateTransactionTx creates a transaction record for an account within a transaction...
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
