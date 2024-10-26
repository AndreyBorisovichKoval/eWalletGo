// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\errors.go

package controllers

import (
	"eWalletGo_TestTask/errs"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleError обрабатывает ошибки и возвращает соответствующий ответ клиенту...
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.ErrWalletNotFound):
		c.JSON(http.StatusNotFound, newErrorResponse(err.Error()))
	case errors.Is(err, errs.ErrSomethingWentWrong):
		c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
	default:
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
