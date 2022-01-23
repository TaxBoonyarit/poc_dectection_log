package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.GetRespHeaders())
		return c.SendString("Hello, World!")
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
