package dictionary

import (
	"github.com/gofiber/fiber/v2"
)

type Dictionary struct {
	English string `json:"english"`
	Chinese string `json:"chinese"`
}

func New(c *fiber.Ctx) error {
	dictionary := new(Dictionary)
	if err := c.BodyParser(dictionary); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	return c.JSON(dictionary)
}
