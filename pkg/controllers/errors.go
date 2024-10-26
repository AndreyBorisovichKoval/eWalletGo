// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\errors.go

package controllers

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleError обрабатывает ошибки и возвращает соответствующий ответ клиенту...
func handleError(c *gin.Context, err error) {
	logger.Error.Printf("Error occurred: %v", err) // Логирование всех ошибок для диагностики

	switch {
	case errors.Is(err, errs.ErrWalletNotFound):
		c.JSON(http.StatusNotFound, newErrorResponse(err.Error()))
	case errors.Is(err, errs.ErrAccountNotFound):
		c.JSON(http.StatusNotFound, newErrorResponse(err.Error()))
	case errors.Is(err, errs.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, newErrorResponse(err.Error()))
	case errors.Is(err, errs.ErrUniquenessViolation):
		c.JSON(http.StatusConflict, newErrorResponse(err.Error()))
	case errors.Is(err, errs.ErrInvalidRequest):
		c.JSON(http.StatusBadRequest, newErrorResponse(err.Error()))
	case errors.Is(err, errs.ErrSomethingWentWrong):
		logger.Error.Printf("Internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
	default:
		logger.Error.Printf("Unhandled error: %v", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse("internal server error"))
	}
}

// ErrorResponse представляет структуру для обработки ошибок...
type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}
