package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Periodicity string

const (
	PeriodicityDaily     Periodicity = "daily"
	PeriodicityDailyEven Periodicity = "daily_even"
	PeriodicityDailyOdd  Periodicity = "daily_odd"
	PediodicityMonthly   Periodicity = "monthly"
)

type Task struct {
	ID               int64       `json:"id"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Status           Status      `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Periodicity      Periodicity `json:"periodicity"`
	PeriodicityValue int         `json:"periodicity_value"`
	PublishDate      time.Time   `json:"publish_date"`
	IsActive         bool        `json:"is_active"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (p Periodicity) Valid() bool {
	switch p {
	case PeriodicityDaily, PeriodicityDailyEven, PeriodicityDailyOdd, PediodicityMonthly:
		return true
	default:
		return false
	}
}
