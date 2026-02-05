package cron

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// JobFunc is the function signature for job handlers.
type JobFunc func(ctx context.Context) error

// Job represents a scheduled job.
type Job struct {
	mu sync.Mutex

	ID          string        // Unique job identifier
	Name        string        // Human-readable name
	Description string        // Job description
	Schedule    Schedule      // When to run
	Func        JobFunc       // The function to execute
	Timeout     time.Duration // Max execution time (0 = no limit)
	Enabled     bool          // Whether the job is active

	// Runtime state
	NextRun    time.Time // Next scheduled run
	LastRun    time.Time // Last run time
	LastError  error     // Last error (nil if successful)
	Running    bool      // Currently running
	RunCount   int64     // Total run count
	ErrorCount int64     // Total error count
}

// NewJob creates a new job with the given schedule.
func NewJob(id, name string, schedule Schedule, fn JobFunc) *Job {
	return &Job{
		ID:       id,
		Name:     name,
		Schedule: schedule,
		Func:     fn,
		Enabled:  true,
	}
}

// WithTimeout sets the job timeout.
func (j *Job) WithTimeout(d time.Duration) *Job {
	j.Timeout = d
	return j
}

// WithDescription sets the job description.
func (j *Job) WithDescription(desc string) *Job {
	j.Description = desc
	return j
}

// Status returns a human-readable status string.
func (j *Job) Status() string {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.Running {
		return "running"
	}
	if !j.Enabled {
		return "disabled"
	}
	if j.LastError != nil {
		return "error"
	}
	return "idle"
}

// ============================================================================
// Schedule Interface and Implementations
// ============================================================================

// Schedule determines when a job should run.
type Schedule interface {
	// Next returns the next time the job should run after the given time.
	Next(after time.Time) time.Time
}

// IntervalSchedule runs at fixed intervals.
type IntervalSchedule struct {
	Interval time.Duration
}

// Every creates an interval schedule.
func Every(d time.Duration) Schedule {
	return &IntervalSchedule{Interval: d}
}

// Next returns the next run time.
func (s *IntervalSchedule) Next(after time.Time) time.Time {
	return after.Add(s.Interval)
}

// DailySchedule runs at a specific time each day.
type DailySchedule struct {
	Hour   int
	Minute int
}

// Daily creates a daily schedule at the specified time.
func Daily(hour, minute int) Schedule {
	return &DailySchedule{Hour: hour, Minute: minute}
}

// Next returns the next run time.
func (s *DailySchedule) Next(after time.Time) time.Time {
	next := time.Date(
		after.Year(), after.Month(), after.Day(),
		s.Hour, s.Minute, 0, 0, after.Location(),
	)

	if !next.After(after) {
		next = next.Add(24 * time.Hour)
	}

	return next
}

// CronSchedule implements cron-style scheduling.
// Supports: minute hour day month weekday
// Examples: "0 9 * * *" (9am daily), "*/15 * * * *" (every 15 min)
type CronSchedule struct {
	Minutes  []int // 0-59
	Hours    []int // 0-23
	Days     []int // 1-31
	Months   []int // 1-12
	Weekdays []int // 0-6 (0 = Sunday)
}

// Cron creates a cron schedule from a cron expression.
func Cron(expr string) (Schedule, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}

	minutes, err := parseCronField(parts[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}

	hours, err := parseCronField(parts[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}

	days, err := parseCronField(parts[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("invalid day field: %w", err)
	}

	months, err := parseCronField(parts[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}

	weekdays, err := parseCronField(parts[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("invalid weekday field: %w", err)
	}

	return &CronSchedule{
		Minutes:  minutes,
		Hours:    hours,
		Days:     days,
		Months:   months,
		Weekdays: weekdays,
	}, nil
}

// MustCron creates a cron schedule, panicking on error.
func MustCron(expr string) Schedule {
	s, err := Cron(expr)
	if err != nil {
		panic(err)
	}
	return s
}

// Next returns the next run time.
func (s *CronSchedule) Next(after time.Time) time.Time {
	// Start from the next minute
	t := after.Truncate(time.Minute).Add(time.Minute)

	// Search for up to 4 years
	maxTime := after.Add(4 * 365 * 24 * time.Hour)

	for t.Before(maxTime) {
		if !contains(s.Months, int(t.Month())) {
			// Skip to next month
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
			continue
		}

		if !contains(s.Days, t.Day()) || !contains(s.Weekdays, int(t.Weekday())) {
			// Skip to next day
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
			continue
		}

		if !contains(s.Hours, t.Hour()) {
			// Skip to next hour
			t = t.Truncate(time.Hour).Add(time.Hour)
			continue
		}

		if !contains(s.Minutes, t.Minute()) {
			// Skip to next minute
			t = t.Add(time.Minute)
			continue
		}

		return t
	}

	// Fallback: just return 1 day later
	return after.Add(24 * time.Hour)
}

// parseCronField parses a single cron field.
func parseCronField(field string, min, max int) ([]int, error) {
	var result []int

	// Handle *
	if field == "*" {
		for i := min; i <= max; i++ {
			result = append(result, i)
		}
		return result, nil
	}

	// Handle */n (step)
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(field[2:])
		if err != nil || step <= 0 {
			return nil, fmt.Errorf("invalid step: %s", field)
		}
		for i := min; i <= max; i += step {
			result = append(result, i)
		}
		return result, nil
	}

	// Handle comma-separated values and ranges
	for _, part := range strings.Split(field, ",") {
		if strings.Contains(part, "-") {
			// Range
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, err
			}
			if start < min || end > max || start > end {
				return nil, fmt.Errorf("range out of bounds: %s", part)
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// Single value
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			if val < min || val > max {
				return nil, fmt.Errorf("value out of bounds: %d", val)
			}
			result = append(result, val)
		}
	}

	return result, nil
}

// contains checks if a slice contains a value.
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
