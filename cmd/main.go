package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/guryn/go-todo/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New(fiber.Config{
		Views: html.New("./internal/templates", ".html"),
	})
	app.Static("/", "./internal/public")

	routes.HttpRoute(app)
	routes.ApiRoute(app)
	routes.WebsocketRoute(app)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
