package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hotelreservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("You role is not admin")
	}
	return c.Next()
}
