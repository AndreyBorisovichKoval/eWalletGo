// C:\GoProject\src\eWalletGo_TestTask\pkg\repository\test_data.go

package repository

import (
	"eWalletGo_TestTask/db"
	"fmt"
	"io/ioutil"
	"os"
)

// CheckTestDataExists проверяет, есть ли уже данные в указанной таблице
func CheckTestDataExists(tableName string) (bool, error) {
	var count int64
	err := db.GetDBConn().Table(tableName).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("ошибка проверки данных в таблице %s: %v", tableName, err)
	}
	return count > 0, nil
}

// LoadSQLFile читает содержимое SQL файла и возвращает его как строку
func LoadSQLFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия SQL файла: %v", err)
	}
	defer file.Close()

	sqlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения SQL файла: %v", err)
	}

	return string(sqlBytes), nil
}

// ExecSQLQuery выполняет переданный SQL-запрос на базе данных
func ExecSQLQuery(sqlQuery string) error {
	if err := db.GetDBConn().Exec(sqlQuery).Error; err != nil {
		return fmt.Errorf("ошибка выполнения SQL запроса: %v", err)
	}
	return nil
}
