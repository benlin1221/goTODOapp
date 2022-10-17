package routes

import (
	apiControllers "m/v2/app/controllers/api"

	"github.com/gofiber/fiber/v2"
)

func TaskRoute(route fiber.Router) {
	route.Get("", apiControllers.GetTasks)
	route.Post("", apiControllers.CreateTask)
	route.Delete("/:id", apiControllers.DeleteTask)
}
