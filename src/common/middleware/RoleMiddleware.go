package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Role Constants
const (
	RoleCliente = 4
	RoleCoach   = 3
	RoleGym     = 2
	RoleAdmin   = 1
)

// RequireRoles verifica que el token JWT del usuario contenga alguno de los roles permitidos.
func RequireRoles(allowedRoles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDAny, exists := c.Get("role_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No hay un rol asociado en el contexto"})
			return
		}

		roleID := roleIDAny.(uint)

		// Verificamos si tiene el rol adecuado
		for _, r := range allowedRoles {
			if roleID == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acceso denegado, no tiene los permisos suficientes"})
	}
}
