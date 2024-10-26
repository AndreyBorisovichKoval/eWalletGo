// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\test_data.go

package controllers

import (
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InsertTestData вызывает сервис для вставки тестовых данных в базу данных
func InsertTestData(c *gin.Context) {
	err := service.InsertTestData()
	if err != nil {
		// Проверяем, если ошибка указывает на уже существующие данные
		if err.Error() == "тестовые данные уже добавлены" {
			c.JSON(http.StatusConflict, gin.H{"error": "Ранее Тестовые Данные уже Были добавлены"})
			return
		}
		// Логируем и отправляем общее сообщение об ошибке
		logger.Error.Printf("[InsertTestData] Ошибка при вставке тестовых данных: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось вставить тестовые данные"})
		return
	}

	logger.Info.Println("[InsertTestData] Тестовые данные успешно добавлены")
	c.JSON(http.StatusOK, gin.H{"message": "Тестовые данные успешно добавлены"})
}
