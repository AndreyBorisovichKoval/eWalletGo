// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\test_data.go

package controllers

import (
	"eWalletGo_TestTask/logger"
	"eWalletGo_TestTask/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InsertTestData calls the service to insert test data into the database
func InsertTestData(c *gin.Context) {
	err := service.InsertTestData()
	if err != nil {
		// Check if the error indicates that data already exists
		if err.Error() == "test data already added" {
			c.JSON(http.StatusConflict, gin.H{"error": "Test data has already been added"})
			return
		}
		// Log and send a generic error message
		logger.Error.Printf("[InsertTestData] Error inserting test data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert test data"})
		return
	}

	logger.Info.Println("[InsertTestData] Test data successfully added")
	c.JSON(http.StatusOK, gin.H{"message": "Test data successfully added"})
}
