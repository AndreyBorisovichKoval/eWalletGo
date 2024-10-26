// C:\GoProject\src\eWalletGo_TestTask\pkg\service\test_data.go

package service

import (
	"eWalletGo_TestTask/logger" // добавляем логгер
	"eWalletGo_TestTask/pkg/repository"
	"errors"
	"fmt"
)

// InsertTestData загружает тестовые данные из SQL файла и проверяет их наличие перед добавлением
func InsertTestData() error {
	// Проверяем, существуют ли уже данные в таблице "users" (выберите таблицу, которая всегда заполняется тестовыми данными)
	exists, err := repository.CheckTestDataExists("users")
	if err != nil {
		logger.Error.Printf("Ошибка проверки существования данных: %v", err) // добавляем логирование ошибки
		return fmt.Errorf("ошибка проверки существования данных: %v", err)
	}
	if exists {
		return errors.New("тестовые данные уже добавлены")
	}

	// Загружаем тестовые данные из файла
	filePath := "C:/GoProject/src/eWalletGo_TestTask/db/insert_test_data.sql"
	sqlData, err := repository.LoadSQLFile(filePath)
	if err != nil {
		logger.Error.Printf("Ошибка загрузки SQL файла: %v", err) // добавляем логирование ошибки
		return fmt.Errorf("ошибка загрузки SQL файла: %v", err)
	}

	// Выполняем SQL-запрос
	err = repository.ExecSQLQuery(sqlData)
	if err != nil {
		logger.Error.Printf("Ошибка выполнения SQL запроса: %v", err) // добавляем логирование ошибки
		return fmt.Errorf("ошибка выполнения SQL запроса: %v", err)
	}

	return nil
}
