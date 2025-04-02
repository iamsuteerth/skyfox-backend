package security

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			log.Error().Msg("JWT_SECRET_KEY is not set in environment variables")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Debug().Msg("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Debug().Msg("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Debug().Err(err).Msg("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
			log.Debug().Interface("claims", claims).Msg("Token claims set in context")
		} else {
			log.Debug().Msg("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized (No claims found)"})
			c.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || role != "admin" {
			log.Debug().Str("role", role).Msg("Access denied for non-admin user")
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden, Only Admins Allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func StaffMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized (No claims Found)"})
			c.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || (role != "admin" && role != "staff") {
			log.Debug().Str("role", role).Msg("Access denied for non-staff user")
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden, Only Staff Allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetTokenClaims(c *gin.Context) (jwt.MapClaims, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, fmt.Errorf("no token claims found in context")
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	return jwtClaims, nil
}
