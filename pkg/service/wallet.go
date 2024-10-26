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

// RechargeWallet пополняет баланс кошелька и добавляет запись транзакции в рамках одной транзакции базы данных...
func RechargeWallet(walletID string, amount float64) error {
	// Начинаем транзакцию
	tx := db.GetDBConn().Begin()
	if tx.Error != nil {
		logger.Error.Println("[service.RechargeWallet] Ошибка начала транзакции:", tx.Error)
		return errs.ErrSomethingWentWrong
	}

	// Проверяем существование кошелька
	exists, err := repository.CheckWalletExistsTx(walletID, tx)
	if err != nil {
		tx.Rollback()
		logger.Error.Printf("[service.RechargeWallet] Ошибка проверки кошелька с ID %s: %v", walletID, err)
		return err
	}
	if !exists {
		tx.Rollback()
		logger.Warning.Printf("[service.RechargeWallet] Кошелек не найден: %s", walletID)
		return errs.ErrWalletNotFound
	}

	// Получаем AccountID
	accountID, err := repository.GetAccountIDByWalletIDTx(walletID, tx)
	if err != nil {
		tx.Rollback()
		logger.Error.Printf("[service.RechargeWallet] Ошибка получения AccountID для кошелька с ID %s: %v", walletID, err)
		return errs.ErrAccountNotFound
	}

	// Обновляем баланс кошелька
	if err := repository.UpdateWalletBalanceTx(walletID, amount, tx); err != nil {
		tx.Rollback()
		logger.Error.Printf("[service.RechargeWallet] Ошибка обновления баланса кошелька с ID %s: %v", walletID, err)
		return err
	}

	// Добавляем запись транзакции
	if err := repository.CreateTransactionTx(accountID, amount, "recharge", tx); err != nil {
		tx.Rollback()
		logger.Error.Printf("[service.RechargeWallet] Ошибка добавления транзакции для AccountID %d: %v", accountID, err)
		return err
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		logger.Error.Println("[service.RechargeWallet] Ошибка при выполнении коммита транзакции:", err)
		return errs.ErrSomethingWentWrong
	}

	logger.Info.Printf("[service.RechargeWallet] Пополнение кошелька успешно завершено: %s, Сумма: %.2f", walletID, amount)
	return nil
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
