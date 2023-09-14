package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	// user, ok := c.Context().UserValue("user").(*user.User)
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return fmt.Errorf("user not found in context")
	}
	if !user.IsAdmin {
		return fmt.Errorf("user is not admin")
	}
	return c.Next()
}
