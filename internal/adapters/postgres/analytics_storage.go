package postgres

import (
	"analytics/internal/domain/models"
	"context"
	"fmt"
	"github.com/google/uuid"
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
		return fmt.Errorf("begin tx failed: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event_uuid (uuid) VALUES ($1)`,
		event.UUID)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.task
				(id, status, created_at, last_mail_at, total_time, approvers_number, current_approvers_number)
			VALUES
				($1, 'created', $2, null, '0 years 0 mons 0 days 0 hours 0 mins 0.0 secs', $3, 0)`,
		event.TaskId, event.Time, event.ApproversNumber)
	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %v", err)
	}
	return nil
}

func (db *Database) AddMail(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event_uuid (uuid) VALUES ($1)`,
		event.UUID)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}
	_, err = tx.Exec(ctx,
		`UPDATE analytics.task
			 SET
				 status = 'waiting_response',
    			 last_mail_at = $2
			 WHERE id = $1`,
		event.TaskId, event.Time)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %v", err)
	}
	return nil
}

func (db *Database) AddApproveClick(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event_uuid (uuid) VALUES ($1)`,
		event.UUID)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
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
		return fmt.Errorf("query exec failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %v", err)
	}
	return nil
}

func (db *Database) AddRejectClick(ctx context.Context, event models.Event) error {
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event_uuid (uuid) VALUES ($1)`,
		event.UUID)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO analytics.event
				 (task_id, occurred_at, event_type, event_user, approvers_number)
			 VALUES ($1, $2, $3, $4, $5)`,
		event.TaskId, event.Time, event.Type, event.User, event.ApproversNumber)

	if err != nil {
		return fmt.Errorf("query exec failed: %v", err)
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
		return fmt.Errorf("query exec failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit failed failed: %v", err)
	}
	return nil
}

func (db *Database) GetTotalTaskResponseTime(ctx context.Context, taskId int32) (string, error) {
	var ival pqinterval.Interval
	var totalTaskResponseTime string

	err := db.DB.QueryRow(ctx, "SELECT total_time FROM analytics.task WHERE task.id = $1", taskId).Scan(&ival)
	if err != nil {
		return "", fmt.Errorf("query row failed: %v", err)
	}

	ivalValue, err := ival.Value()
	if err != nil {
		return "", fmt.Errorf("execute interval value failed: %v", err)
	}
	totalTaskResponseTime = fmt.Sprint(ivalValue)
	return totalTaskResponseTime, nil
}

func (db *Database) GetTasksCount(ctx context.Context, taskType string) (int32, error) {
	var approvedTasksCount int32
	var err error

	if taskType == "approved" {
		err = db.DB.QueryRow(ctx,
			"SELECT count(id) FROM analytics.task WHERE status = 'approved'").Scan(&approvedTasksCount)
	} else if taskType == "rejected" {
		err = db.DB.QueryRow(ctx,
			"SELECT count(id) FROM analytics.task WHERE status = 'rejected'").Scan(&approvedTasksCount)
	} else {
		return 0, fmt.Errorf("get tasks count failed: %s taskType not appicable", taskType)
	}

	if err != nil {
		return 0, fmt.Errorf("query row failed: %v", err)
	}

	return approvedTasksCount, nil
}

func (db *Database) CheckIdempotency(ctx context.Context, uuid uuid.UUID) (bool, error) {
	rows, err := db.DB.Query(ctx,
		"SELECT * FROM analytics.event_uuid WHERE uuid = $1", uuid)
	if err != nil {
		return false, fmt.Errorf("query exec failed: %v", err)
	}

	if !rows.Next() {
		return true, nil
	}
	return false, nil
}