package api

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yuriykis/hotel-reservation/db"
)

func JWTAutentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnauthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return ErrUnauthorized()
		}
		expiresFloat := claims["expires"].(float64)
		// if time.Now().After(time.Unix(int64(expires), 0)) {
		// 	return fmt.Errorf("token expired")
		// }
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return NewError(fiber.StatusUnauthorized, "token expired")
		}
		userID, ok := claims["id"].(string)
		if !ok {
			return ErrUnauthorized()
		}
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}
		// set the current user in the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Unexpected signing method")
			return nil, ErrUnauthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		log.Println(secret)
		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err)
		return nil, ErrUnauthorized()
	}
	if !token.Valid {
		fmt.Println("Token is not valid")
		return nil, ErrUnauthorized()
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrUnauthorized()
}
