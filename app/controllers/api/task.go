package controllers

import (
	"fmt"
	"m/v2/app/models"
	"m/v2/app/repos"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func CreateTodo(c *fiber.Ctx) error {
	var body models.TaskDTO
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}
	aTask, err := repos.CreateTask(body.Title, body.Assignee)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	result := &models.TaskResponse{}
	if err := copier.Copy(&result, &aTask); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func GetTasks(c *fiber.Ctx) error {
	pagination, err := getContextPagination(c, 1000, models.TaskResponse{})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if pagination.Sort == "" {
		pagination.Sort = "updated_at asc" // use as default sort behaviour
	}

	if r, err := repos.GetAllTasks(*pagination); err == nil && r.Rows != nil {
		result := []models.TaskResponse{}
		rows := r.Rows.([]models.Task)
		if err := copier.Copy(&result, &rows); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Cannot map results",
			})
		}
		r.Rows = result
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    r,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"message": "Cannot fetch results",
	})
}
