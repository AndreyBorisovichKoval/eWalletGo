// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\errors.go

package controllers

import (
	"eWalletGo_TestTask/errs"
	"eWalletGo_TestTask/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleError handles errors and returns an appropriate response to the client...
func handleError(c *gin.Context, err error) {
	logger.Error.Printf("Error occurred: %v", err) // Logging all errors for diagnostics

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
	case errors.Is(err, errs.ErrLimitExceeded):
		c.JSON(http.StatusForbidden, newErrorResponse(err.Error())) // Code for limit exceeded
	case errors.Is(err, errs.ErrSomethingWentWrong):
		logger.Error.Printf("Internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
	default:
		logger.Error.Printf("Unhandled error: %v", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse("internal server error"))
	}
}

// ErrorResponse represents the structure for handling errors...
type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}
