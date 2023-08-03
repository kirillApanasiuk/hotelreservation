package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotelreservation/db"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("-- JWT auth")

		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnAuthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		fmt.Printf("wil expired at ---> %v\n", claims["expires"])
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userId := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userId)
		if err != nil {
			return ErrUnAuthorized()
		}
		// Set the current authenticated user to the context.
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("invalid signing method", token.Header["alg"])
				return nil, ErrUnAuthorized()
			}
			secret := os.Getenv("JWT_SECRET")
			fmt.Println("Never print secret", secret)

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secret), nil
		},
	)

	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token", err)
		return nil, ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnAuthorized()
	}

	return claims, nil
}
