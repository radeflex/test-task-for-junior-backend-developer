package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		Periodicity:      normalized.Periodicity,
		PeriodicityValue: normalized.PediodicityValue,
		PublishDate:      normalized.PublishDate,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now
	model.IsActive = calcActivity(model)

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:               id,
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		Periodicity:      normalized.Periodicity,
		PeriodicityValue: normalized.PediodicityValue,
		PublishDate:      normalized.PublishDate,
		UpdatedAt:        s.now(),
	}

	model.IsActive = calcActivity(model)

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if err := validatePeriodicity(input.Periodicity, input.PediodicityValue); err != nil {
		return CreateInput{}, err
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}
	if err := validatePeriodicity(input.Periodicity, input.PediodicityValue); err != nil {
		return UpdateInput{}, err
	}
	return input, nil
}

func validatePeriodicity(periodicity *taskdomain.Periodicity, value *int) error {
	if periodicity == nil {
		return nil
	}

	if !periodicity.Valid() {
		return fmt.Errorf("%w: invalid periodicity", ErrInvalidInput)
	}

	if (*periodicity == taskdomain.PeriodicityDaily ||
		*periodicity == taskdomain.PeriodicityMonthly) &&
		value == nil {
		return fmt.Errorf("%w: periodicity value is required", ErrInvalidInput)
	}

	if value == nil {
		return nil
	}

	v := *value

	switch *periodicity {
	case taskdomain.PeriodicityDaily:
		if v < 1 || v > 7 {
			return fmt.Errorf("%w: daily periodicity must be between 1 and 7", ErrInvalidInput)
		}

	case taskdomain.PeriodicityMonthly:
		if v < 1 || v > 30 {
			return fmt.Errorf("%w: monthly periodicity must be between 1 and 30", ErrInvalidInput)
		}
	}

	return nil
}

func (s *Service) UpdateActivity(ctx context.Context) error {
	tasks, err := s.repo.List(ctx)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		active := calcActivity(&t)
		if active != t.IsActive {
			if t.Status == taskdomain.StatusDone {
				t.Status = taskdomain.StatusNew
			}
			t.IsActive = active
			t.UpdatedAt = s.now()
			_, err := s.repo.Update(ctx, &t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func calcActivity(task *taskdomain.Task) bool {
	if task.PublishDate == nil || task.PublishDate.Before(time.Now()) {
		if task.Periodicity == nil {
			return true
		}
		switch *task.Periodicity {
		case taskdomain.PeriodicityDailyEven:
			return time.Now().Day()%2 == 0
		case taskdomain.PeriodicityDailyOdd:
			return time.Now().Day()%2 == 1
		case taskdomain.PeriodicityDaily:
			return int(time.Now().Weekday()) == *task.PeriodicityValue
		case taskdomain.PeriodicityMonthly:
			return time.Now().Day() == *task.PeriodicityValue
		}
	}
	return false
}
