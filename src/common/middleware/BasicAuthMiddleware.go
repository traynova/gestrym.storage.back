package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	BasicAuthUsernameKey = "BASIC_AUTH_USERNAME"
	BasicAuthPasswordKey = "BASIC_AUTH_PASSWORD"
)

func SetupBasicAuthMiddleware() gin.HandlerFunc {
	return BasicAuthMiddleware()
}

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"mensaje": "Preflight request"})
			return
		}

		auth := c.GetHeader("Authorization")
		if auth == "" {
			logger.Warn("[BASIC_AUTH] El encabezado de autorización está ausente")
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[BASIC_AUTH] autorización requerida"})
			return
		}

		username, password, ok := parseBasicAuth(auth)
		if !ok {
			logger.Warn("[BASIC_AUTH] Formato de autenticación básica inválido")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[BASIC_AUTH] formato de autenticación inválido"})
			return
		}

		if !validateCredentials(username, password) {
			logger.Warn("[BASIC_AUTH] Intento de autenticación fallido para el usuario:", username)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[BASIC_AUTH] credenciales inválidas"})
			return
		}

		// Autenticación exitosa
		c.Next()
	}
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	if !strings.HasPrefix(auth, "Basic ") {
		return "", "", false
	}
	payload, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return "", "", false
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", false
	}
	return pair[0], pair[1], true
}

func validateCredentials(username, password string) bool {
	basicAuthUsername := viper.GetString(BasicAuthUsernameKey)
	basicAuthPassword := viper.GetString(BasicAuthPasswordKey)

	if basicAuthUsername == "" || basicAuthPassword == "" {
		logger.Error("[BASIC_AUTH] Las variables de entorno BASIC_AUTH_USERNAME y/o BASIC_AUTH_PASSWORD no están configuradas")
		return false
	}

	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(basicAuthUsername)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(basicAuthPassword)) == 1

	return usernameMatch && passwordMatch
}
