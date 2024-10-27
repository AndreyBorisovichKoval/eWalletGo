// C:\GoProject\src\eWalletGo_TestTask\pkg\service\test_data.go

package service

import (
	"eWalletGo_TestTask/logger" // add logger
	"eWalletGo_TestTask/pkg/repository"
	"errors"
	"fmt"
)

// InsertTestData loads test data from an SQL file and checks for existing data before adding
func InsertTestData() error {
	// Check if data already exists in the "users" table (choose a table that is always populated with test data)
	exists, err := repository.CheckTestDataExists("users")
	if err != nil {
		logger.Error.Printf("Error checking data existence: %v", err) // add error logging
		return fmt.Errorf("error checking data existence: %v", err)
	}
	if exists {
		return errors.New("test data already added")
	}

	// Load test data from file
	filePath := "C:/GoProject/src/eWalletGo_TestTask/db/insert_test_data.sql"
	sqlData, err := repository.LoadSQLFile(filePath)
	if err != nil {
		logger.Error.Printf("Error loading SQL file: %v", err) // add error logging
		return fmt.Errorf("error loading SQL file: %v", err)
	}

	// Execute SQL query
	err = repository.ExecSQLQuery(sqlData)
	if err != nil {
		logger.Error.Printf("Error executing SQL query: %v", err) // add error logging
		return fmt.Errorf("error executing SQL query: %v", err)
	}

	return nil
}
