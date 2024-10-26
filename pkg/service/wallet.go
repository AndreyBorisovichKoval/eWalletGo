// C:\GoProject\src\eWalletGo_TestTask\pkg\service\wallet.go

package service

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/pkg/repository"
	"fmt"
)

// CheckWalletExistence проверяет наличие кошелька по ID и возвращает результат...
func CheckWalletExistence(walletID string) error {
	exists, err := repository.CheckWalletExists(walletID)
	if err != nil {
		return fmt.Errorf("error checking wallet existence: %w", err)
	}
	if !exists {
		return errs.ErrWalletNotFound // кастомная ошибка, если кошелек не найден...
	}
	return nil
}
