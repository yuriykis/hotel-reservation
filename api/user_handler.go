package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Bond",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"user": c.Params("id"),
	})
}
