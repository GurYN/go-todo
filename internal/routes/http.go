package routes

import (
	"github.com/gofiber/fiber/v2"
)

func HttpRoute(app *fiber.App) {
	r := app.Group("/")

	r.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "GO Todo - Welcome",
		})
	})
}
