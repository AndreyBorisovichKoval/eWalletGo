// C:\GoProject\src\eWalletGo_TestTask\pkg\repository\gorm.go

package repository

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Переводим ошибки GORM (конвертируем) в кастомные ошибки...
func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warning.Printf("Record not found error: %v...", err)
		return errs.ErrRecordNotFound
	}

	// Check for uniqueness violation error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		logger.Warning.Printf("Uniqueness violation: %v...", err)
		return errs.ErrUniquenessViolation
	}

	// Добавить логирование для других ошибок по мере необходимости...
	logger.Error.Printf("Unhandled error: %v...", err)

	// return err
	return errs.ErrSomethingWentWrong
}
