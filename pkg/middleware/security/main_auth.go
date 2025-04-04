package security

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			log.Error().Msg("JWT_SECRET_KEY is not set in environment variables")
			utils.HandleErrorResponse(c,
				utils.NewInternalServerError("SERVER_CONFIG_ERROR", "Server configuration error", fmt.Errorf("JWT_SECRET_KEY not set")),
				requestID)
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Debug().Msg("Missing Authorization header")
			utils.HandleErrorResponse(c,
				utils.NewUnauthorizedError("MISSING_AUTH_HEADER", "Missing Authorization header", nil),
				requestID)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Debug().Msg("Invalid Authorization header format")
			utils.HandleErrorResponse(c,
				utils.NewUnauthorizedError("INVALID_TOKEN_FORMAT", "Invalid token format", nil),
				requestID)
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
			utils.HandleErrorResponse(c,
				utils.NewUnauthorizedError("INVALID_TOKEN", "Unauthorized", err),
				requestID)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
			log.Debug().Interface("claims", claims).Msg("Token claims set in context")
		} else {
			log.Debug().Msg("Invalid token claims")
			utils.HandleErrorResponse(c,
				utils.NewUnauthorizedError("INVALID_TOKEN_CLAIMS", "Invalid token claims", nil),
				requestID)
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		claims, exists := c.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			utils.HandleErrorResponse(c,
				utils.NewUnauthorizedError("NO_CLAIMS_FOUND", "No authentication claims found", nil),
				requestID)
			c.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || role != "admin" {
			log.Debug().Str("role", role).Msg("Access denied for non-admin user")
			utils.HandleErrorResponse(c,
				utils.NewBadRequestError("FORBIDDEN", "Access denied. Admin role required", nil),
				requestID)
			c.Abort()
			return
		}

		c.Next()
	}
}

func StaffMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		claims, exists := c.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			utils.HandleErrorResponse(c,
				utils.NewUnauthorizedError("NO_CLAIMS_FOUND", "No authentication claims found", nil),
				requestID)
			c.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || (role != "staff") {
			log.Debug().Str("role", role).Msg("Access denied for non-staff user")
			utils.HandleErrorResponse(c,
				utils.NewBadRequestError("FORBIDDEN", "Access denied. Staff or admin role required", nil),
				requestID)
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
