// C:\GoProject\src\eWalletGo_TestTask\pkg\middleware\auth.go

package middleware

// const (
// 	userIDHeader = "X-UserId"
// 	digestHeader = "X-Digest"
// 	secretKey    = "your_secret_key" // Замените на реальный секретный ключ
// )

// // validateHMAC проверяет, совпадает ли полученный HMAC с ожидаемым.
// func validateHMAC(body, digest string) bool {
// 	h := hmac.New(sha1.New, []byte(secretKey))
// 	h.Write([]byte(body))
// 	expectedMAC := hex.EncodeToString(h.Sum(nil))
// 	return hmac.Equal([]byte(expectedMAC), []byte(digest))
// }

// // AuthMiddleware проверяет наличие и корректность заголовков X-UserId и X-Digest.
// func AuthMiddleware(c *gin.Context) {
// 	userID := c.GetHeader(userIDHeader)
// 	digest := c.GetHeader(digestHeader)

// 	if userID == "" || digest == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication headers"})
// 		c.Abort()
// 		return
// 	}

// 	// Читаем тело запроса
// 	body, err := c.GetRawData()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
// 		c.Abort()
// 		return
// 	}

// 	// Проверка HMAC
// 	if !validateHMAC(string(body), digest) {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid digest"})
// 		c.Abort()
// 		return
// 	}

// 	// Установка userID в контекст для дальнейшего использования
// 	c.Set("userID", userID)
// 	c.Request.Body = http.MaxBytesReader(c.Writer, http.NoBody, int64(len(body))) // Восстанавливаем тело запроса
// 	c.Next()
// }
