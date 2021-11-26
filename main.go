package main

import (
	"os"

	"github.com/programzheng/program-english/dictionary"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Index!")
	})

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")
	dictionaryGroup := v1Group.Group("/dictionary")
	dictionaryGroup.Post("", dictionary.New)

	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
