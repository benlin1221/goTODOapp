package main

import (
	"log"
	"m/v2/app/models"
	configuration "m/v2/config"
	"m/v2/database"
	apiRoutes "m/v2/routes/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type App struct {
	*fiber.App

	DB      *database.Database
	Session *session.Store
}

func main() {
	config := configuration.GetInstance()
	app := setupApp(config)

	if err := app.Listen(config.GetString("APP_ADDR")); err != nil {
		log.Panic(err)
	}
}

func setupApp(config *configuration.Config) App {
	database.Setup()

	app := App{
		App:     fiber.New(*config.GetFiberConfig()),
		Session: session.New(config.GetSessionConfig()),
	}

	app.DB = (&database.Database{
		DB: database.DBConn,
	})

	database.SessionStore = app.Session
	app.Session.RegisterType("")
	var typeUint uint = 1
	app.Session.RegisterType(typeUint)
	var typeBool bool = false
	app.Session.RegisterType(typeBool)

	setupRoutes(app.App)
	return app
}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		_, err := database.SessionStore.Get(c)
		if err != nil {
			panic(err)
		}

		var tasks []models.Task
		database.DBConn.Association("Tasks").Find(&tasks)

		template := "dashboard"

		return c.Render(template, fiber.Map{
			"Title": tasks,
		})
	})

	api := app.Group("/api")

	api.Get("", func(c *fiber.Ctx) error {
		_, err := database.SessionStore.Get(c)
		if err != nil {
			panic(err)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the api endpoint 😉",
		})
	})

	apiRoutes.TaskRoute(api.Group("/tasks"))
}
