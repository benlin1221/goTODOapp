package services

import (
	"errors"
	"fmt"
	"log"
	"m/v2/app/models"
	"m/v2/app/repos"
	"m/v2/app/utils"
	"m/v2/app/utils/password"
	"m/v2/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const signinMethodEmailPassword = "EmailPassword"
const signinMethodOAuth2AzureAD = "OAuth2AzureAD"

// Login service logs in a user
// if userAzureADLinkToken is present in body and not expired
// azure ad immutable id will be linked
func Login(ctx *fiber.Ctx) error {
	b := new(models.LoginDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		return err
	}
	u := &models.User{}

	err := repos.FindUserByUsername(u, b.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	if err := password.Verify(u.Password, b.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	sess, err := database.SessionStore.Get(ctx)
	if err != nil {
		panic(err)
	}

	// Set key/value
	sess.Set("ID", u.ID)
	sess.Set("username", u.Username)

	// save session
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return ctx.JSON(&models.UserResponse{
		ID:       u.ID,
		Username: u.Username,
	})
}

// Logout service logs out a user
func Logout(ctx *fiber.Ctx) error {
	sess, err := database.SessionStore.Get(ctx)
	if err != nil {
		panic(err)
	}

	// Destry session
	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "You are logged out 😉",
	})
}

// Signup service creates a user
func Signup(ctx *fiber.Ctx) error {
	b := new(models.SignupDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		return err
	}

	user := &models.User{
		Username: b.Username,
		Password: password.Generate(b.Password),
	}

	// Create a user, if error return
	if err := repos.CreateUser(user); err != nil {
		return fiber.NewError(fiber.StatusConflict, "Error creating account")
	}

	sess, err := database.SessionStore.Get(ctx)
	if err != nil {
		panic(err)
	}

	// Set key/value
	sess.Set("ID", user.ID)
	sess.Set("username", user.Username)

	// save session
	if err := sess.Save(); err != nil {
		panic(err)
	}

	s := fmt.Sprintf("Username: %s, ID: %d has signed up.", user.Username, user.ID)
	log.Printf(s)

	return ctx.JSON(&models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	})
}
