package repository

import (
	"context"
	"errors"
	"fmt"

	"simpletodo/internal/model"
	"simpletodo/internal/storage/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type Task struct {
	storage *postgres.Storage
	log     *zap.Logger
}

func NewTaskRepository(storage *postgres.Storage, log *zap.Logger) *Task {
	return &Task{storage: storage, log: log}
}

func (t *Task) Create(ctx context.Context, task *model.Task) error {
	const op = "repository.task.Create"

	query, args, err := t.storage.Builder.Insert("tasks").
		Columns("title", "done").
		Values(task.Title, task.Done).
		ToSql()
	if err != nil {
		t.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	_, err = t.storage.Pool.Exec(ctx, query, args...)
	if err != nil {
		t.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	return nil
}

func (t *Task) GetByID(ctx context.Context, id int) (*model.Task, error) {
	const op = "repository.task.GetByID"

	if id == 0 {
		t.log.Error(op, zap.Error(model.ErrInvalidID))

		return nil, model.ErrInvalidID
	}

	query, args, err := t.storage.Builder.
		Select("id", "title", "done").
		From("tasks").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		t.log.Error(op, zap.Error(err))

		return nil, fmt.Errorf("%s : %w", op, err)
	}

	var task model.Task

	err = t.storage.Pool.QueryRow(ctx, query, args...).
		Scan(&task.ID, &task.Title, &task.Done)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			t.log.Error(op, zap.Error(model.ErrNotFound))

			return nil, model.ErrNotFound
		}

		t.log.Error(op, zap.Error(err))

		return nil, fmt.Errorf("%s : %w", op, err)
	}

	t.log.Info(fmt.Sprintf("%s : successfully got task with id %d", op, id))

	return &task, nil
}

func (t *Task) GetAll(ctx context.Context) ([]model.Task, error) {
	const op = "repository.task.GetAll"

	t.log.Info(op + " : getting tasks")

	query, args, err := t.storage.Builder.
		Select("id", "title", "done").
		From("tasks").
		ToSql()
	if err != nil {
		t.log.Error(op, zap.Error(err))

		return nil, fmt.Errorf("%s : %w", op, err)
	}

	rows, err := t.storage.Pool.Query(ctx, query, args...)

	defer func() {
		rows.Close()
	}()

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			t.log.Error(op, zap.Error(model.ErrNotFound))

			return nil, model.ErrNotFound
		}

		t.log.Error(op, zap.Error(err))

		return nil, fmt.Errorf("%s : %w", op, err)
	}

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done); err != nil {
			t.log.Error(op, zap.Error(err))

			return nil, fmt.Errorf("%s : %w", op, err)
		}

		tasks = append(tasks, task)
	}

	t.log.Info(fmt.Sprintf("%s : successfully got %d tasks", op, len(tasks)))

	return tasks, nil
}

func (t *Task) Update(ctx context.Context, task *model.Task) error {
	const op = "repository.task.Update"

	if task == nil {
		t.log.Error(op, zap.Error(model.ErrInvalidTask))

		return model.ErrInvalidTask
	}

	if task.ID == 0 {
		t.log.Error(op, zap.Error(model.ErrInvalidID))

		return model.ErrInvalidID
	}

	query, args, err := t.storage.Builder.Update("tasks").
		Set("title", task.Title).
		Set("done", task.Done).
		Where(squirrel.Eq{"id": task.ID}).
		ToSql()
	if err != nil {
		t.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	cmdTag, err := t.storage.Pool.Exec(ctx, query, args...)
	if cmdTag.RowsAffected() == 0 {
		t.log.Error(op, zap.Error(model.ErrNotFound))

		return model.ErrNotFound
	}

	if err != nil {
		t.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	t.log.Info(fmt.Sprintf("%s :updated %d task with id %d", op, cmdTag.RowsAffected(), task.ID))

	return nil
}

func (t *Task) Delete(ctx context.Context, id int) error {
	const op = "repository.task.Delete"

	if id == 0 {
		t.log.Error(op, zap.Error(model.ErrInvalidID))

		return model.ErrInvalidID
	}

	query, args, err := t.storage.Builder.Delete("tasks").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		t.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	cmdTag, err := t.storage.Pool.Exec(ctx, query, args...)
	if cmdTag.RowsAffected() == 0 {
		t.log.Error(op, zap.Error(fmt.Errorf("%w with id %d", model.ErrNotFound, id)))

		return model.ErrNotFound
	}

	if err != nil {
		t.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	t.log.Info(fmt.Sprintf("%s : deleted %d task with id %d", op, cmdTag.RowsAffected(), id))

	return nil
}
