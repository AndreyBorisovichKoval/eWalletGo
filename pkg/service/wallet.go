// C:\GoProject\src\eWalletGo_TestTask\pkg\service\wallet.go

package service

import (
	"eWalletGo_TestTask/db"
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/pkg/repository"
	"errors"
)

// CheckWalletExistence проверяет наличие кошелька по ID и возвращает результат...
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

// GetMonthlyRechargeSummary возвращает суммарные данные по кошельку для указанного месяца и года
func GetMonthlyRechargeSummary(walletID string, year int, month int) (int64, float64, error) {
	totalCount, totalAmount, err := repository.GetMonthlyRechargeSummary(walletID, year, month)
	if err != nil {
		logger.Error.Printf("[service.GetMonthlyRechargeSummary] Ошибка получения данных: %v", err)
		return 0, 0, err
	}
	return totalCount, totalAmount, nil
}

// GetWalletBalance возвращает текущий баланс кошелька...
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

// RecalculateWalletBalance пересчитывает баланс кошелька на основе всех транзакций...
func RecalculateWalletBalance(walletID string) (float64, error) {
	// Проверяем существование кошелька
	exists, err := repository.CheckWalletExists(walletID)
	if err != nil {
		return 0, err
	}
	if !exists {
		logger.Warning.Printf("[service.RecalculateWalletBalance] Wallet not found: %s", walletID)
		return 0, errs.ErrWalletNotFound
	}

	// Получаем пересчитанный баланс на основе транзакций
	newBalance, err := repository.CalculateBalanceFromTransactions(walletID)
	if err != nil {
		logger.Error.Printf("[service.RecalculateWalletBalance] Error recalculating balance for wallet ID %s: %v", walletID, err)
		return 0, err
	}

	// Обновляем баланс в аккаунте
	if err := repository.UpdateWalletBalanceDirectly(walletID, newBalance); err != nil {
		logger.Error.Printf("[service.RecalculateWalletBalance] Error updating wallet balance for wallet ID %s: %v", walletID, err)
		return 0, err
	}

	return newBalance, nil
}

// RechargeWallet пополняет баланс кошелька и проверяет лимиты
func RechargeWallet(walletID string, amount float64) error {
	tx := db.GetDBConn().Begin()
	if tx.Error != nil {
		logger.Error.Println("[service.RechargeWallet] Ошибка начала транзакции:", tx.Error)
		return errs.ErrSomethingWentWrong
	}

	// Проверка существования кошелька и получение лимита
	walletWithLimit, err := repository.GetWalletWithLimit(walletID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Проверка, превышает ли новое пополнение лимит
	newBalance := walletWithLimit.Balance + amount
	if newBalance > walletWithLimit.MaxLimit {
		tx.Rollback()
		return errs.ErrLimitExceeded // Ошибка, если новый баланс превышает лимит
	}

	// Обновляем баланс
	err = repository.UpdateWalletBalanceTx(walletID, amount, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Добавляем запись транзакции
	err = repository.CreateTransactionTx(walletWithLimit.AccountID, amount, "recharge", tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		logger.Error.Println("[service.RechargeWallet] Ошибка при выполнении коммита транзакции:", err)
		return errs.ErrSomethingWentWrong
	}

	logger.Info.Printf("[service.RechargeWallet] Кошелек пополнен: %s, Сумма: %.2f", walletID, amount)
	return nil
}
