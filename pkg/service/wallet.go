// C:\GoProject\src\eWalletGo_TestTask\pkg\service\wallet.go

package service

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/pkg/repository"
	"errors"
)

// CheckWalletExistence проверяет наличие кошелька по ID и возвращает результат...
func CheckWalletExistence(walletID string) (bool, error) {
	exists, err := repository.CheckWalletExists(walletID)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return false, errs.ErrWalletNotFound // Преобразуем общую ошибку в специфичную для кошелька
		}
		return false, errs.ErrSomethingWentWrong
	}
	return exists, nil
}
