package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/guryn/go-todo/internal/routes"
)

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./internal/templates", ".html"),
	})
	app.Static("/", "./internal/public")

	routes.Root(app)

	app.Listen(":3000")
}
