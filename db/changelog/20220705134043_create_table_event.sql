-- +goose Up

CREATE TABLE IF NOT EXISTS analytics.event
(
    task_id          int REFERENCES analytics.task(id),
    occurred_at      timestamptz NOT NULL,
    event_type       text NOT NULL,
    event_user       text NOT NULL,
    approvers_number int
);
