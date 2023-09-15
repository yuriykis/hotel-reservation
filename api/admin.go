package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	// user, ok := c.Context().UserValue("user").(*user.User)
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if !user.IsAdmin {
		return ErrUnauthorized()
	}
	return c.Next()
}
