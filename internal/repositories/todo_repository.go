package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/guryn/go-todo/internal/models"
)

type TodoRepository interface {
	FindAll() ([]models.Todo, error)
	FindById(id uuid.UUID) (models.Todo, error)
	Create(todo models.Todo) (models.Todo, error)
	Update(todo models.Todo) (models.Todo, error)
	Delete(id uuid.UUID) error
}

type TodoRepositoryImpl struct {
	db []models.Todo
}

func NewTodoRepository() TodoRepository {
	return &TodoRepositoryImpl{}
}

func (repository *TodoRepositoryImpl) FindAll() ([]models.Todo, error) {
	return repository.db, nil
}

func (repository *TodoRepositoryImpl) FindById(id uuid.UUID) (models.Todo, error) {
	for _, todo := range repository.db {
		if todo.ID == id {
			return todo, nil
		}
	}
	return models.Todo{}, errors.New("not found")
}

func (repository *TodoRepositoryImpl) Create(todo models.Todo) (models.Todo, error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	todo.ID = newId
	todo.Completed = false
	repository.db = append(repository.db, todo)
	return todo, nil
}

func (repository *TodoRepositoryImpl) Update(todo models.Todo) (models.Todo, error) {
	for index, t := range repository.db {
		if t.ID == todo.ID {
			repository.db[index] = todo
			return todo, nil
		}
	}
	return models.Todo{}, errors.New("not found")
}

func (repository *TodoRepositoryImpl) Delete(id uuid.UUID) error {
	for index, todo := range repository.db {
		if todo.ID == id {
			repository.db = append(repository.db[:index], repository.db[index+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
