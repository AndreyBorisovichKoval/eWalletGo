// C:\GoProject\src\eWalletGo_TestTask\models\response.go

package models

// Response structure for a standard API response
type Response struct {
	Status  string      `json:"status"`         // Response status: "success" or "error"
	Message string      `json:"message"`        // Description of the result
	Data    interface{} `json:"data,omitempty"` // Data returned by the API (if any)
}
