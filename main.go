package main

import (
	"log"
	"m/v2/app/middleware"
	"m/v2/app/services"
	configuration "m/v2/config"
	"m/v2/database"
	"m/v2/routes"
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

	app.Get("/", middleware.RequireSession, func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTemporaryRedirect).Redirect("/dashboard/")
	})

	app.Get("/dashboard", middleware.RequireSession, func(c *fiber.Ctx) error {
		sess, err := database.SessionStore.Get(c)
		if err != nil {
			panic(err)
		}

		user, err := services.UserTemplateFromContext(c)
		if err != nil {
			log.Printf("Error: could not find user (%s)\n", err.Error())
			if err := sess.Destroy(); err != nil {
				log.Print("Error destroying session", err.Error())
			}
			log.Print("User not found from main")
			return fiber.NewError(fiber.StatusInternalServerError, "Session user not found")
		}

		return c.Render("dashboard", fiber.Map{
			"User": user,
		})
	})

	app.Static("/", "./static/public", fiber.Static{
		Compress:  true,
		ByteRange: true,
		Browse:    false,
		Index:     "/",
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("signup", fiber.Map{})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the logout endpoint 😉",
		})
	})

	authGroup := app.Group("/auth")

	routes.AuthRoutes(authGroup)

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
