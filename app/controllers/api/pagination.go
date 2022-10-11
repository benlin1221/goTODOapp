package controllers

import (
	"strconv"

	"m/v2/app/models"

	"github.com/gofiber/fiber/v2"
)

// Returns a Pagination struct with query params from context.
// Model is used as a whitelist for sort field names and is
// mapped to json tag names.
func getContextPagination(c *fiber.Ctx, maxLimit int, model interface{}) (*models.Pagination, error) {
	// get Query value
	queryLimit := c.Query("limit")
	var limit int
	// convert parameter value string to int
	if v, err := strconv.ParseInt(queryLimit, 10, 32); err == nil {
		limit = int(v)
	} else {
		limit = maxLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	queryPage := c.Query("page")
	var page int
	// convert parameter value string to int
	if v, err := strconv.ParseInt(queryPage, 10, 32); err == nil {
		page = int(v)
	}
	sortQuery := c.Query("sort")
	sortDescQuery := c.Query("sortDesc")
	sortDesc, err := strconv.ParseBool(sortDescQuery)
	if err != nil {
		sortDesc = false
	}
	pagination := &models.Pagination{
		Limit: limit,
		Page:  page,
	}
	if err := pagination.SetSort(sortQuery, sortDesc, model); err != nil {
		return nil, err
	}
	return pagination, nil
}
