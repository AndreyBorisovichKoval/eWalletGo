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
