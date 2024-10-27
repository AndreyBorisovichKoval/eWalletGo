// C:\GoProject\src\eWalletGo_TestTask\main.go

package main

import (
	app "eWalletGo_TestTask/cmd"
	"eWalletGo_TestTask/utils"
)

// @title eWalletGo API
// @version 1.0
// @description API Server for eWalletGo Test Task Application

// @host localhost:57320
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Clear the console from old messages...
	utils.ClearConsole()

	app.RunApp()
}
