// C:\GoProject\src\eWalletGo_TestTask\models\response.go

package models

// Response структура для стандартного ответа API
type Response struct {
	Status  string      `json:"status"`         // Статус ответа: "success" или "error"
	Message string      `json:"message"`        // Описание результата выполнения
	Data    interface{} `json:"data,omitempty"` // Данные, возвращаемые API (если есть)
}
