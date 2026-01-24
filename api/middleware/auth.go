package middleware

import (
	"fmt"
	"project_m_backend/pkg/auth/jwt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	golangJwt "github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(jwtService *jwt.JWTService) fiber.Handler{
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key:[]byte(jwtService.GetSecretKey())},
		KeyFunc: func(token *golangJwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*golangJwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(jwtService.GetSecretKey()), nil
		},
		SuccessHandler: func(c *fiber.Ctx) error{
			userToken := c.Locals("user").(*golangJwt.Token)
			claims := userToken.Claims.(golangJwt.MapClaims)
			encryptedUserID := claims["sub"].(string)
			userID, err := jwtService.DecryptUserID(encryptedUserID)
			if err != nil{
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
			}

			c.Locals("user_id", userID)
			return c.Next()
		},
		ErrorHandler: func (c *fiber.Ctx, err error)  error{
			fmt.Println("JWT error:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired access token"})
		},
	
	})
}