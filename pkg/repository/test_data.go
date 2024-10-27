// C:\GoProject\src\eWalletGo_TestTask\pkg\repository\test_data.go

package repository

import (
	"eWalletGo_TestTask/db"
	"fmt"
	"io/ioutil"
	"os"
)

// CheckTestDataExists checks if data already exists in the specified table
func CheckTestDataExists(tableName string) (bool, error) {
	var count int64
	err := db.GetDBConn().Table(tableName).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error checking data in table %s: %v", tableName, err)
	}
	return count > 0, nil
}

// LoadSQLFile reads the content of an SQL file and returns it as a string
func LoadSQLFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening SQL file: %v", err)
	}
	defer file.Close()

	sqlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading SQL file: %v", err)
	}

	return string(sqlBytes), nil
}

// ExecSQLQuery executes the given SQL query on the database
func ExecSQLQuery(sqlQuery string) error {
	if err := db.GetDBConn().Exec(sqlQuery).Error; err != nil {
		return fmt.Errorf("error executing SQL query: %v", err)
	}
	return nil
}
