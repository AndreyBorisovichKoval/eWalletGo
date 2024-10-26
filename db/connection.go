// C:\GoProject\src\eWalletGo_TestTask\db\connection.go

package db

import (
	"eWalletGo_TestTask/configs"
	"eWalletGo_TestTask/logger"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

// EnsureDatabaseExists проверяет наличие базы данных и создает её, если не существует.
func EnsureDatabaseExists() error {
	createDBConnStr := fmt.Sprintf("host=%s port=%s user=%s dbname=postgres password=%s TimeZone=Asia/Dushanbe",
		configs.AppSettings.PostgresParams.Host,
		configs.AppSettings.PostgresParams.Port,
		configs.AppSettings.PostgresParams.User,
		os.Getenv("DB_PASSWORD"))

	logger.Info.Println("Connecting to check if the database exists...")

	tempDB, err := gorm.Open(postgres.Open(createDBConnStr), &gorm.Config{NowFunc: func() time.Time {
		return time.Now().Local()
	}})
	if err != nil {
		logger.Error.Printf("Failed to connect to 'postgres': %v", err)
		return err
	}

	var exists bool
	checkDBQuery := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", configs.AppSettings.PostgresParams.Database)
	tempDB.Raw(checkDBQuery).Scan(&exists)

	if !exists {
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", configs.AppSettings.PostgresParams.Database)
		if err := tempDB.Exec(createDBQuery).Error; err != nil {
			logger.Error.Printf("Failed to create database '%s': %v", configs.AppSettings.PostgresParams.Database, err)
			return err
		}
		logger.Info.Printf("Database '%s' created successfully!", configs.AppSettings.PostgresParams.Database)
	} else {
		logger.Info.Printf("Database '%s' already exists.", configs.AppSettings.PostgresParams.Database)
	}

	sqlTempDB, _ := tempDB.DB()
	sqlTempDB.Close()

	return nil
}

// ConnectToDB устанавливает соединение с базой данных PostgreSQL.
func ConnectToDB() error {
	if err := EnsureDatabaseExists(); err != nil {
		return err
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s TimeZone=Asia/Dushanbe",
		configs.AppSettings.PostgresParams.Host,
		configs.AppSettings.PostgresParams.Port,
		configs.AppSettings.PostgresParams.User,
		configs.AppSettings.PostgresParams.Database,
		os.Getenv("DB_PASSWORD"))

	logger.Info.Println("Connecting to the database...")

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{NowFunc: func() time.Time {
		return time.Now().Local()
	}})
	if err != nil {
		logger.Error.Printf("Failed to connect to the database: %v", err)
		return err
	}

	logger.Info.Println("Successfully connected to the database!")
	dbConn = db
	return nil
}

// CloseDBConn закрывает соединение с базой данных.
func CloseDBConn() error {
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Error.Printf("Failed to retrieve SQL DB from gorm.DB: %v", err)
		return err
	}
	if err = sqlDB.Close(); err != nil {
		logger.Error.Printf("Failed to close database connection: %v", err)
		return err
	}
	logger.Info.Println("Database connection closed successfully!")
	return nil
}

// GetDBConn возвращает текущее соединение с базой данных.
func GetDBConn() *gorm.DB {
	return dbConn
}
