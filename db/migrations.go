// C:\GoProject\src\eWalletGo_TestTask\db\migrations.go

package db

import (
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/models"
)

// MigrateDB выполняет миграцию базы данных, создавая необходимые таблицы.
func MigrateDB() error {
	logger.Info.Println("Starting database migration...")
	err := dbConn.AutoMigrate(
		models.User{},
		models.Phone{},
		models.Wallet{},
		models.Account{},
		models.Transaction{},
		models.LimitSettings{},
		models.UnverifiedUser{},
		models.UserSettings{},
		models.RequestHistory{},
	)
	if err != nil {
		logger.Error.Printf("Migration failed: %v", err)
		return err
	}

	logger.Info.Println("Database migration completed successfully!")
	return nil
}
