package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/guryn/go-todo/internal/models"
)

type TodoService interface {
	GetList() ([]models.Todo, error)
	Create(todo models.Todo) (models.Todo, error)
	Complete(id string) error
	Delete(id string) error
}

type TodoServiceImpl struct {
	apiURL string
}

func NewTodoService() TodoService {
	return &TodoServiceImpl{apiURL: os.Getenv("API_URL")}
}

func (service *TodoServiceImpl) GetList() ([]models.Todo, error) {
	request := fiber.Get(fmt.Sprintf("%s/api/todo", service.apiURL))
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

func (service *TodoServiceImpl) Create(d models.Todo) (models.Todo, error) {
	request := fiber.Post(fmt.Sprintf("%s/api/todo", service.apiURL))
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

func (service *TodoServiceImpl) Complete(id string) error {
	request := fiber.Get(fmt.Sprintf("%s/api/todo/%s/complete", service.apiURL, id))
	_, _, errs := request.Bytes()
	if errs != nil {
		return errs[0]
	}

	return nil
}

func (service *TodoServiceImpl) Delete(id string) error {
	request := fiber.Delete(fmt.Sprintf("%s/api/todo/%s", service.apiURL, id))
	_, _, errs := request.Bytes()
	if errs != nil {
		return errs[0]
	}

	return nil
}
