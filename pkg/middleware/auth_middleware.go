package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		userId, valid := isValidToken(authToken)
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

func isValidToken(token string) (string, bool) {
	validationURL := os.Getenv("AUTH_SERVICE_ENDPOINT_LOCAL")
	req, err := http.NewRequest("POST", validationURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", false
	}

	req.Header.Set("Access-Token", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Unauthorized for this action, invalid token response")
		return "", false
	}

	var responseBody struct {
		UserID string `json:"userId"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return "", false
	}

	return responseBody.UserID, true
}
