package service

import (
	"context"
	"simpletodo/internal/model"

	"go.uber.org/zap"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id int) (*model.Task, error)
	GetAll(ctx context.Context) ([]model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskService struct {
	log        *zap.Logger
	repository TaskRepository
}

func New(
	log *zap.Logger,
	repository TaskRepository,
) *TaskService {
	return &TaskService{
		log:        log,
		repository: repository,
	}
}

func (t *TaskService) Create(ctx context.Context, title string) error {
	task := model.NewTask(title)

	return t.repository.Create(ctx, task)
}

func (t *TaskService) GetAll(ctx context.Context) ([]model.Task, error) {
	return t.repository.GetAll(ctx)
}

func (t *TaskService) GetByID(ctx context.Context, id int) (*model.Task, error) {
	return t.repository.GetByID(ctx, id)
}

func (t *TaskService) Update(ctx context.Context, task *model.Task) error {
	return t.repository.Update(ctx, task)
}

func (t *TaskService) Delete(ctx context.Context, id int) error {
	return t.repository.Delete(ctx, id)
}
