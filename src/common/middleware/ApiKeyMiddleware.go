package middleware

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func ValidateAPIKey(encodedApiKey string) gin.HandlerFunc {
	apiKey, err := base64.StdEncoding.DecodeString(encodedApiKey)
	if err != nil {
		panic("Error al decodificar la API Key: " + err.Error())
	}
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		providedApiKey := c.GetHeader("X-API-Key")
		if providedApiKey != string(apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key no válida o ausente"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func SetupApiKeyMiddleware() gin.HandlerFunc {
	encodedApiKey := viper.GetString("AUTH_API_KEY")
	if encodedApiKey == "" {
		logger.Error("[API_KEY_MIDDLEWARE_CONFIG_FAILED] API_KEY no está configurado en el archivo de configuración")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "AUTH_API_KEY no está configurado"})
		}
	}
	return ValidateAPIKey(encodedApiKey)
}
