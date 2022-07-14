package postgres

import (
	"analytics/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/sanyokbig/pqinterval"
	"time"
)

type Datapoint struct {
	Timestamp time.Time
}

func (db *Database) CreateTask(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.task
				(id, status, created_at, last_mail_at, total_time, approvers_number, current_approvers_number)
			VALUES
				($1, 'created', $2, null, '0 years 0 mons 0 days 0 hours 0 mins 0.0 secs', $3, 0)`,
		event.TaskId, event.Time, event.ApproversNumber)
	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}

func (db *Database) AddMail(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}
	_, err = tx.Exec(ctx,
		`UPDATE analytics.task
			 SET
				 status = 'waiting_response',
    			 last_mail_at = $2
			 WHERE id = $1`,
		event.TaskId, event.Time)

	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}

func (db *Database) AddApproveClick(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}
	_, err = tx.Exec(ctx,
		`UPDATE analytics.task
			 SET
     			 status = CASE
                 	WHEN current_approvers_number + 1 = approvers_number THEN 'approved'
                 	ELSE 'response_received'
				 END,
    			 total_time = total_time + ($2::timestamp - last_mail_at),
				 current_approvers_number =  current_approvers_number + 1
			 WHERE id = $1`,
		event.TaskId, event.Time)
	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}

func (db *Database) AddRejectClick(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}
	_, err = tx.Exec(ctx,
		`UPDATE analytics.task
			 SET
    			 status = 'rejected',
    			 total_time = total_time + ($2::timestamp - last_mail_at),
				 current_approvers_number =  current_approvers_number + 1
			 WHERE id = $1`,
		event.TaskId, event.Time)
	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}

func (db *Database) GetTotalTaskResponseTime(ctx context.Context, taskId int32) (string, error) {
	var ival pqinterval.Interval
	var totalTaskResponseTime string

	rows, err := db.DB.Query(ctx, "SELECT total_time FROM analytics.task WHERE task.id = $1", taskId)
	if err != nil {
		return "", fmt.Errorf("query exec failed: %w", err)
	}

	if !rows.Next() {
		return "", errors.New("not found")
	}

	err = rows.Scan(&ival)
	if err != nil {
		return "", fmt.Errorf("scan exec failed: %w", err)
	}

	ivalValue, err := ival.Value()
	if err != nil {
		return "", fmt.Errorf("execute interval value failed: %w", err)
	}
	totalTaskResponseTime = fmt.Sprint(ivalValue)
	return totalTaskResponseTime, nil
}

func (db *Database) GetApprovedTasksCount(ctx context.Context) (int32, error) {
	var approvedTasksCount int32

	rows, err := db.DB.Query(ctx, "SELECT count(id) FROM analytics.task WHERE status = 'approved'")
	if err != nil {
		return -1, fmt.Errorf("query exec failed: %w", err)
	}

	if !rows.Next() {
		return -1, errors.New("not found")
	}

	err = rows.Scan(&approvedTasksCount)
	if err != nil {
		return -1, fmt.Errorf("scan exec failed: %w", err)
	}

	return approvedTasksCount, nil
}

func (db *Database) GetRejectedTasksCount(ctx context.Context) (int32, error) {
	var rejectedTasksCount int32

	rows, err := db.DB.Query(ctx,
		"SELECT count(id) FROM analytics.task WHERE status IN ('created', 'waiting_response', 'response_received', 'rejected')")
	if err != nil {
		return -1, fmt.Errorf("query exec failed: %w", err)
	}

	if !rows.Next() {
		return -1, errors.New("not found")
	}

	err = rows.Scan(&rejectedTasksCount)
	if err != nil {
		return -1, fmt.Errorf("scan exec failed: %w", err)
	}

	return rejectedTasksCount, nil
}
