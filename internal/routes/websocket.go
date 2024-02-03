package routes

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/guryn/go-todo/internal/models"
)

func WebsocketRoute(app *fiber.App) {
	r := app.Group("/ws/todo")
	websockets := []websocket.Conn{}

	r.Get("/", websocket.New(func(c *websocket.Conn) {
		websockets = append(websockets, *c)

		for {
			_, msg, read_err := c.ReadMessage()
			if read_err != nil {
				c.Close()
				break
			}
			req := models.ActionRequest{}
			if json.Unmarshal(msg, &req) != nil {
				c.WriteJSON([]byte("{\"error\": \"Failed to parse request\"}}"))
				continue
			}

			switch string(req.Action) {
			case "get_list":
				result, err := getList()
				if err != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to retrieve todos\"}}"))
					continue
				}

				if c.WriteJSON(models.TodoList{Answer: "todo_list", Data: result}) != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to send todos\"}}"))
					continue
				}
			case "create":
				result, err := create(req.Data)
				if err != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to create todo\"}}"))
					continue
				}

				for _, ws := range websockets {
					if ws.WriteJSON(models.TodoList{Answer: "added_item", Data: result}) != nil {
						ws.WriteJSON([]byte("{\"error\": \"Failed to send todo\"}}"))
						continue
					}
				}
			case "complete":
				err := complete(req.Data.(string))
				if err != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to complete todo\"}}"))
					continue
				}

				for _, ws := range websockets {
					if ws.WriteJSON(models.TodoList{Answer: "completed_item", Data: req.Data.(string)}) != nil {
						ws.WriteJSON([]byte("{\"error\": \"Failed to send todo\"}}"))
						continue
					}
				}
			case "delete":
				err := delete(req.Data.(string))
				if err != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to delete todo\"}}"))
					continue
				}

				for _, ws := range websockets {
					if ws.WriteJSON(models.TodoList{Answer: "deleted_item", Data: req.Data.(string)}) != nil {
						ws.WriteJSON([]byte("{\"error\": \"Failed to send todo\"}}"))
						continue
					}
				}
			default:
				c.WriteJSON([]byte("{\"error\": \"Unknown command\"}}"))
			}
		}

		for i, ws := range websockets {
			if ws.Conn == c.Conn {
				websockets = append(websockets[:i], websockets[i+1:]...)
				break
			}
		}
	}))
}

func getList() ([]models.Todo, error) {
	request := fiber.Get(fmt.Sprintf("%s/api/todo", os.Getenv("API_URL")))
	_, data, errs := request.Bytes()
	if errs != nil {
		return nil, errs[0]
	}

	lst := []models.Todo{}
	err := json.Unmarshal(data, &lst)
	if err != nil {
		return nil, err
	}

	return lst, nil
}

func create(d any) (models.Todo, error) {
	request := fiber.Post(fmt.Sprintf("%s/api/todo", os.Getenv("API_URL")))
	inp, err := json.Marshal(d)
	if err != nil {
		return models.Todo{}, err
	}

	_, data, errs := request.Body(inp).Bytes()
	if errs != nil {
		return models.Todo{}, errs[0]
	}

	item := models.Todo{}
	parse_err := json.Unmarshal(data, &item)
	if parse_err != nil {
		return models.Todo{}, parse_err
	}

	return item, nil
}

func complete(id string) error {
	request := fiber.Get(fmt.Sprintf("%s/api/todo/%s/complete", os.Getenv("API_URL"), id))
	_, _, errs := request.Bytes()
	if errs != nil {
		return errs[0]
	}

	return nil
}

func delete(id string) error {
	request := fiber.Delete(fmt.Sprintf("%s/api/todo/%s", os.Getenv("API_URL"), id))
	_, _, errs := request.Bytes()
	if errs != nil {
		return errs[0]
	}

	return nil
}
