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
	return func(ctx *gin.Context) {
		requestID := utils.GetRequestID(ctx)

		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			log.Error().Msg("JWT_SECRET_KEY is not set in environment variables")
			utils.HandleErrorResponse(ctx,
				utils.NewInternalServerError("SERVER_CONFIG_ERROR", "Server configuration error", fmt.Errorf("JWT_SECRET_KEY not set")),
				requestID)
			ctx.Abort()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Debug().Msg("Missing Authorization header")
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("MISSING_AUTH_HEADER", "Missing Authorization header", nil),
				requestID)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Debug().Msg("Invalid Authorization header format")
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("INVALID_TOKEN_FORMAT", "Invalid token format", nil),
				requestID)
			ctx.Abort()
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
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("INVALID_TOKEN", "Unauthorized", err),
				requestID)
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("claims", claims)
			log.Debug().Interface("claims", claims).Msg("Token claims set in context")
		} else {
			log.Debug().Msg("Invalid token claims")
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("INVALID_TOKEN_CLAIMS", "Invalid token claims", nil),
				requestID)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func CustomerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := utils.GetRequestID(ctx)

		claims, exists := ctx.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("NO_CLAIMS_FOUND", "No authentication claims found", nil),
				requestID)
			ctx.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || role != "customer" {
			log.Debug().Str("role", role).Msg("Access denied for non-customers.")
			utils.HandleErrorResponse(ctx,
				utils.NewForbiddenError("FORBIDDEN", "Access denied. Customer role required", nil),
				requestID)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := utils.GetRequestID(ctx)

		claims, exists := ctx.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("NO_CLAIMS_FOUND", "No authentication claims found", nil),
				requestID)
			ctx.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || role != "admin" {
			log.Debug().Str("role", role).Msg("Access denied for non-admin user")
			utils.HandleErrorResponse(ctx,
				utils.NewForbiddenError("FORBIDDEN", "Access denied. Admin role required", nil),
				requestID)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func AdminStaffMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := utils.GetRequestID(ctx)

		claims, exists := ctx.Get("claims")
		if !exists {
			log.Debug().Msg("No claims found in context")
			utils.HandleErrorResponse(ctx,
				utils.NewUnauthorizedError("NO_CLAIMS_FOUND", "No authentication claims found", nil),
				requestID)
			ctx.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || !(role == "admin" || role == "staff") {
			log.Debug().Str("role", role).Msg("Access denied for non-admin or non-staff user")
			utils.HandleErrorResponse(ctx,
				utils.NewForbiddenError("FORBIDDEN", "Access denied. Admin or staff role required!", nil),
				requestID)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func GetTokenClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	claims, exists := ctx.Get("claims")
	if !exists {
		return nil, fmt.Errorf("no token claims found in context")
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	return jwtClaims, nil
}
