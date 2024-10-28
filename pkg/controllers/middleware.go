// C:\GoProject\src\eWalletGo_TestTask\pkg\controllers\middleware.go

package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	userIDHeader = "X-UserId"
	digestHeader = "X-Digest"
	secretKey    = "Fred_secret_key" // Установите уникальный секретный ключ
)

// validateHMAC проверяет, совпадает ли полученный HMAC с ожидаемым.
func validateHMAC(body, digest string) bool {
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(body))
	expectedMAC := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedMAC), []byte(digest))
}

// AuthMiddleware проверяет наличие и корректность заголовков X-UserId и X-Digest.
func AuthMiddleware(c *gin.Context) {
	userID := c.GetHeader(userIDHeader)
	digest := c.GetHeader(digestHeader)

	if userID == "" || digest == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication headers"})
		c.Abort()
		return
	}

	// Читаем тело запроса и восстанавливаем его для дальнейшей обработки
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		c.Abort()
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Восстанавливаем тело запроса

	// Проверка HMAC
	if !validateHMAC(string(body), digest) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid digest"})
		c.Abort()
		return
	}

	c.Set("userID", userID)
	c.Next()
}
