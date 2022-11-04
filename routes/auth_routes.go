package routes

import (
	"m/v2/app/services"

	"github.com/gofiber/fiber/v2"
)

// AuthRoutes containes all the auth routes
func AuthRoutes(route fiber.Router) {
	route.Post("/logout", services.Logout)
	route.Post("/signup", services.Signup)
	route.Post("/login", services.Login)
}
