package controllers

import (
	"fmt"
	"m/v2/app/models"
	"m/v2/app/repos"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func CreateTask(c *fiber.Ctx) error {
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
	if r, err := repos.GetAllTasks(); err == nil && len(*r) > 0 {
		result := []models.TaskResponse{}
		if err := copier.Copy(&result, &r); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Cannot map results",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"message": "Cannot fetch results",
	})
}

func DeleteTask(c *fiber.Ctx) error {
	// get parameter value
	paramId := c.Params("id")
	var id uint
	// convert parameter value string to int
	if v, err := strconv.ParseUint(paramId, 10, 32); err == nil {
		id = uint(v)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}

	// find Todo and return
	if err := repos.DeleteTask(id); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}

	// if Todo not available
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Todo not found",
	})
}
