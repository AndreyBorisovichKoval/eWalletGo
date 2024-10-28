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
	secretKey    = "Fred_secret_key" // Set a unique secret key
)

// validateHMAC checks if the received HMAC matches the expected one.
func validateHMAC(body, digest string) bool {
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(body))
	expectedMAC := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedMAC), []byte(digest))
}

// AuthMiddleware checks for the presence and validity of X-UserId and X-Digest headers.
func AuthMiddleware(c *gin.Context) {
	userID := c.GetHeader(userIDHeader)
	digest := c.GetHeader(digestHeader)

	if userID == "" || digest == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication headers"})
		c.Abort()
		return
	}

	// Read the request body and restore it for further processing
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		c.Abort()
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Restore request body

	// Validate HMAC
	if !validateHMAC(string(body), digest) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid digest"})
		c.Abort()
		return
	}

	c.Set("userID", userID)
	c.Next()
}
