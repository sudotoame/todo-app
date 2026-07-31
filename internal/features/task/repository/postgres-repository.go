package repository

import (
	"context"
	"errors"
	apperror "todo-app/internal/core/app-error"
	"todo-app/internal/core/domains/task"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	Pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		Pool: pool,
	}
}

func (p *PostgresRepo) AddTask(ctx context.Context, t task.Task) (*task.Task, error) {
	sqlQuery := `
		INSERT INTO tasks (
		title, 
		description
		)
		VALUES ($1, $2)
		RETURNING id, title, description, created_at;
	`

	rows, err := p.Pool.Query(
		ctx,
		sqlQuery,
		t.Title,
		t.Description,
	)
	if err != nil {
		return nil, err
	}
	tr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Task])
	if err != nil {
		return nil, err
	}

	return &task.Task{
		ID:          tr.ID,
		Title:       tr.Title,
		Description: tr.Description,
		CreatedAt:   tr.CreatedAt,
	}, nil
}

func (p *PostgresRepo) GetTaskByID(ctx context.Context, id int) (*task.Task, error) {
	sqlQueryByID := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE id = $1;
	`

	var t task.Task
	err := p.Pool.QueryRow(ctx, sqlQueryByID, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Completed,
		&t.CreatedAt,
		&t.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrTaskNotFound
		} else {
			return nil, err
		}
	}

	return &t, nil
}

func (p *PostgresRepo) ListTasks(ctx context.Context) ([]task.Task, error) {
	sqlQuery := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks;
	`
	rows, err := p.Pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	t, err := pgx.CollectRows(rows, pgx.RowToStructByName[task.Task])

	return t, nil
}

func (p *PostgresRepo) ListTasksByComplete(ctx context.Context, completed *bool) ([]task.Task, error) {
	sqlQueryByComplete := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE completed = $1;
	`

	rows, err := p.Pool.Query(ctx, sqlQueryByComplete, completed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var task task.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (p *PostgresRepo) SetCompleteTask(ctx context.Context, id int, t task.Task) error {
	if t.Completed == false {
		t.CompletedAt = nil
	}

	sqlQuery := `
		UPDATE tasks
		SET completed=$1, completed_at=$2
		WHERE id = $3
	`

	_, err := p.Pool.Exec(ctx, sqlQuery, t.Completed, t.CompletedAt, id)

	return err
}

func (p *PostgresRepo) UpdateTask(ctx context.Context, id int, t task.Task) error {
	sqlQuery := `
		UPDATE tasks
		SET title=$1, description=$2
		WHERE id = $3
	`

	_, err := p.Pool.Exec(ctx, sqlQuery, t.Title, t.Description, id)

	return err
}

func (p *PostgresRepo) DeleteTask(ctx context.Context, id int) error {
	sqlQuery := `
		DELETE FROM tasks
		WHERE id = $1
	`

	_, err := p.Pool.Exec(ctx, sqlQuery, id)

	return err
}
