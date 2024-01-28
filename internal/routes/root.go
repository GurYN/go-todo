package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guryn/go-todo/internal/models"
	"github.com/guryn/go-todo/internal/repositories"
)

func Root(app *fiber.App) {
	r := app.Group("/")
	todos := repositories.NewTodoRepository()

	r.Get("/", func(c *fiber.Ctx) error {
		todosList, err := todos.FindAll()
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
			"Tasks": todosList,
		})
	})

	r.Get("/details/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}

		todo, err := todos.FindById(id)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.SendString(fmt.Sprintf("%v", todo))
	})

	r.Post("/create", func(c *fiber.Ctx) error {
		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
			return c.SendString(err.Error())
		}

		todos.Create(todo)

		return c.Redirect("/")
	})

	r.Get("/complete/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}

		todo, err := todos.FindById(id)
		if err != nil {
			return c.SendString(err.Error())
		}

		todo.Completed = true
		_, err = todos.Update(todo)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.Redirect("/")
	})

	r.Get("/delete/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}

		err = todos.Delete(id)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.Redirect("/")
	})
}
