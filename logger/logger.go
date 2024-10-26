// C:\GoProject\src\eWalletGo_TestTask\logger\logger.go

package logger

import (
	"eWalletGo_TestTask/configs"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Warning *log.Logger
	Debug   *log.Logger
)

func Init() error {
	logParams := configs.AppSettings.LogParams

	// Проверка и создание директории для логов, если её нет
	if _, err := os.Stat(logParams.LogDirectory); os.IsNotExist(err) {
		err = os.Mkdir(logParams.LogDirectory, 0755)
		if err != nil {
			return err
		}
	}

	// Инициализация логгеров с помощью lumberjack
	lumberLogInfo := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogInfo),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	lumberLogError := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogError),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	lumberLogWarning := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogWarning),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	lumberLogDebug := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogDebug),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	// Настройка вывода логов Gin на stdout и в файл Info
	gin.DefaultWriter = io.MultiWriter(os.Stdout, lumberLogInfo)

	// Создание глобальных логгеров с установленным префиксом и форматом
	Info = log.New(gin.DefaultWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(lumberLogError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(lumberLogWarning, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(lumberLogDebug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
