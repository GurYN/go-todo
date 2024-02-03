package routes

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/guryn/go-todo/internal/models"
	"github.com/guryn/go-todo/internal/services"
)

func WebsocketRoute(app *fiber.App) {
	r := app.Group("/ws/todo")
	wsClients := []websocket.Conn{}

	r.Get("/", websocket.New(func(c *websocket.Conn) {
		todoService := services.NewTodoService()
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
				result, err := todoService.GetList()
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

				result, err := todoService.Create(todo)
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
				err := todoService.Complete(req.Data.(string))
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
				err := todoService.Delete(req.Data.(string))
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
