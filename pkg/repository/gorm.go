// C:\GoProject\src\eWalletGo_TestTask\pkg\repository\gorm.go

package repository

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// translateError переводит ошибки Go в кастомные ошибки и логирует их...
func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warning.Printf("Record not found: %v...", err)
		return errs.ErrRecordNotFound
	}

	// Проверка на нарушение уникальности
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		logger.Warning.Printf("Uniqueness violation: %v...", err)
		return errs.ErrUniquenessViolation
	}

	// Логирование других ошибок
	logger.Error.Printf("Unhandled error: %v...", err)

	return errs.ErrSomethingWentWrong
}
