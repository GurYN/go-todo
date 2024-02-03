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
	wsClients := []websocket.Conn{}

	r.Get("/", websocket.New(func(c *websocket.Conn) {
		wsClients = append(wsClients, *c)

		for {
			_, msg, read_err := c.ReadMessage()
			if read_err != nil {
				c.Close()
				break
			}
			req := models.WebsocketRequest{}
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

				if c.WriteJSON(models.WebsocketResponse{Answer: "todo_list", Data: result}) != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to send todos\"}}"))
					continue
				}
			case "create":
				todo := models.Todo{}
				todo.Title = req.Data.(map[string]interface{})["title"].(string)

				result, err := create(todo)
				if err != nil {
					c.WriteJSON([]byte("{\"error\": \"Failed to create todo\"}}"))
					continue
				}

				for _, ws := range wsClients {
					if ws.WriteJSON(models.WebsocketResponse{Answer: "added_item", Data: result}) != nil {
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

				for _, ws := range wsClients {
					if ws.WriteJSON(models.WebsocketResponse{Answer: "completed_item", Data: req.Data.(string)}) != nil {
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

				for _, ws := range wsClients {
					if ws.WriteJSON(models.WebsocketResponse{Answer: "deleted_item", Data: req.Data.(string)}) != nil {
						ws.WriteJSON([]byte("{\"error\": \"Failed to send todo\"}}"))
						continue
					}
				}
			default:
				c.WriteJSON([]byte("{\"error\": \"Unknown command\"}}"))
			}
		}

		for i, ws := range wsClients {
			if ws.Conn == c.Conn {
				wsClients = append(wsClients[:i], wsClients[i+1:]...)
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

	output := []models.Todo{}
	err := json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func create(d models.Todo) (models.Todo, error) {
	request := fiber.Post(fmt.Sprintf("%s/api/todo", os.Getenv("API_URL")))
	input, err := json.Marshal(d)
	if err != nil {
		return models.Todo{}, err
	}

	_, data, errs := request.Body(input).Bytes()
	if errs != nil {
		return models.Todo{}, errs[0]
	}

	output := models.Todo{}
	parse_err := json.Unmarshal(data, &output)
	if parse_err != nil {
		return models.Todo{}, parse_err
	}

	return output, nil
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
