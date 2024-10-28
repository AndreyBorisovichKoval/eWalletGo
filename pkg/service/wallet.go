// C:\GoProject\src\eWalletGo_TestTask\pkg\service\wallet.go

package service

import (
	"eWalletGo_TestTask/db"
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/pkg/repository"
	"errors"
)

// CheckWalletExistence checks for wallet existence by ID and returns the result...
func CheckWalletExistence(walletID string) (bool, error) {
	exists, err := repository.CheckWalletExists(walletID)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return false, errs.ErrWalletNotFound
		}
		return false, errs.ErrSomethingWentWrong
	}
	return exists, nil
}

// GetMonthlyRechargeSummary returns summary data for the wallet for the specified month and year
func GetMonthlyRechargeSummary(walletID string, year int, month int) (int64, float64, error) {
	totalCount, totalAmount, err := repository.GetMonthlyRechargeSummary(walletID, year, month)
	if err != nil {
		logger.Error.Printf("[service.GetMonthlyRechargeSummary] Error retrieving data: %v", err)
		return 0, 0, err
	}
	return totalCount, totalAmount, nil
}

// GetWalletBalance retrieves the current wallet balance...
func GetWalletBalance(walletID string) (float64, error) {
	balance, err := repository.GetWalletBalance(walletID)
	if err != nil {
		if err == errs.ErrWalletNotFound {
			logger.Warning.Printf("[service.GetWalletBalance] Wallet not found: %s", walletID)
			return 0, errs.ErrWalletNotFound
		}
		logger.Error.Printf("[service.GetWalletBalance] Error retrieving wallet balance for wallet ID %s: %v", walletID, err)
		return 0, errs.ErrSomethingWentWrong
	}
	return balance, nil
}

// RecalculateWalletBalance recalculates the wallet balance based on all transactions...
func RecalculateWalletBalance(walletID string) (float64, error) {
	// Check if wallet exists
	exists, err := repository.CheckWalletExists(walletID)
	if err != nil {
		return 0, err
	}
	if !exists {
		logger.Warning.Printf("[service.RecalculateWalletBalance] Wallet not found: %s", walletID)
		return 0, errs.ErrWalletNotFound
	}

	// Get recalculated balance based on transactions
	newBalance, err := repository.CalculateBalanceFromTransactions(walletID)
	if err != nil {
		logger.Error.Printf("[service.RecalculateWalletBalance] Error recalculating balance for wallet ID %s: %v", walletID, err)
		return 0, err
	}

	// Update balance in the account
	if err := repository.UpdateWalletBalanceDirectly(walletID, newBalance); err != nil {
		logger.Error.Printf("[service.RecalculateWalletBalance] Error updating wallet balance for wallet ID %s: %v", walletID, err)
		return 0, err
	}

	return newBalance, nil
}

// RechargeWallet recharges the wallet balance and checks limits
func RechargeWallet(walletID string, amount float64) error {
	tx := db.GetDBConn().Begin()
	if tx.Error != nil {
		logger.Error.Println("[service.RechargeWallet] Error starting transaction:", tx.Error)
		return errs.ErrSomethingWentWrong
	}

	// Check wallet existence and get limit
	walletWithLimit, err := repository.GetWalletWithLimit(walletID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if the new balance exceeds the limit
	newBalance := walletWithLimit.Balance + amount
	// if newBalance > walletWithLimit.MaxLimit {
	if newBalance > walletWithLimit.MaxLimit || newBalance < 0 {
		tx.Rollback()
		return errs.ErrLimitExceeded // Error if new balance exceeds the limit
	}

	// Update balance
	err = repository.UpdateWalletBalanceTx(walletID, amount, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add transaction record
	err = repository.CreateTransactionTx(walletWithLimit.AccountID, amount, "recharge", tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error.Println("[service.RechargeWallet] Error committing transaction:", err)
		return errs.ErrSomethingWentWrong
	}

	logger.Info.Printf("[service.RechargeWallet] Wallet recharged: %s, Amount: %.2f", walletID, amount)
	return nil
}
