package routes

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guryn/go-todo/internal/models"
	"github.com/guryn/go-todo/internal/repositories"
)

func ApiRoute(app *fiber.App) {
	r := app.Group("/api/todo")
	todos := repositories.NewTodoRepository()

	r.Get("/", func(c *fiber.Ctx) error {
		todosList, err := todos.FindAll()
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.JSON(todosList)
	})

	r.Get("/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}

		todo, err := todos.FindById(id)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.JSON(todo)
	})

	r.Post("/", func(c *fiber.Ctx) error {
		todo := models.Todo{}
		err := json.Unmarshal(c.Body(), &todo)
		if err != nil {
			return c.SendString("Failed to parse request: %v" + err.Error())
		}

		result, add_err := todos.Create(todo)
		if add_err != nil {
			return c.SendString(add_err.Error())
		}

		return c.JSON(result)
	})

	r.Put("/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}

		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
			return c.SendString(err.Error())
		}

		todo.ID = id
		_, err = todos.Update(todo)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.JSON(todo)
	})

	r.Get("/:id/complete", func(c *fiber.Ctx) error {
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

		return c.SendString("Completed")
	})

	r.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}

		err = todos.Delete(id)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.SendString("Deleted")
	})
}
