package middleware

import (
	"net/http"
	"strings"

	"github.com/pchawandi/xm-company/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	admin_role = "admin"
	user_role  = "user"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization Header"})
			c.Abort()
			return
		}

		tokenStr := header[len(BearerSchema):]
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if c.Request.Method != http.MethodGet && claims.Role != admin_role {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "insufficient user previleges"})
			c.Abort()
			return
		}

		c.Next()
	}
}
