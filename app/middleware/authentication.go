package middleware

import (
	"m/v2/database"

	"github.com/gofiber/fiber/v2"
)

// RequireSession checks if the user is logged in.
func RequireSession(c *fiber.Ctx) error {
	sess, err := database.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendFile("./static/private/500.html")
	}
	if sess.Get("username") == nil {
		return c.Status(fiber.ErrUnauthorized.Code).Render("login", fiber.Map{})
	}

	return c.Next()
}
