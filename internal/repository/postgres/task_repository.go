package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (
			title,
			description,
			status,
			created_at,
			updated_at,
			periodicity_type,
			periodicity_value,
			publish_date,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING
			id,
			title,
			description,
			status,
			created_at,
			updated_at,
			periodicity_type,
			periodicity_value,
			publish_date,
			is_active
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
		task.Periodicity,
		task.PeriodicityValue,
		task.PublishDate,
		task.IsActive,
	)

	return scanTask(row)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			status,
			created_at,
			updated_at,
			periodicity_type,
			periodicity_value,
			publish_date,
			is_active
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return task, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			updated_at = $4,
			periodicity_type = $5,
			periodicity_value = $6,
			publish_date = $7,
			is_active = $8
		WHERE id = $9
		RETURNING
			id,
			title,
			description,
			status,
			created_at,
			updated_at,
			periodicity_type,
			periodicity_value,
			publish_date,
			is_active
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
		task.Periodicity,
		task.PeriodicityValue,
		task.PublishDate,
		task.IsActive,
		task.ID,
	)

	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			status,
			created_at,
			updated_at,
			periodicity_type,
			periodicity_value,
			publish_date,
			is_active
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task        taskdomain.Task
		status      string
		periodicity string
	)

	err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&periodicity,
		&task.PeriodicityValue,
		&task.PublishDate,
		&task.IsActive,
	)
	if err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)
	task.Periodicity = taskdomain.Periodicity(periodicity)
	return &task, nil
}
