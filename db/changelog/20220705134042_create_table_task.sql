-- +goose Up

CREATE TABLE IF NOT EXISTS analytics.task
(
    id                        int PRIMARY KEY,
    status                    text NOT NULL,
    created_at                timestamptz NOT NULL,
    last_mail_at              timestamptz,
    total_time                interval,
    approvers_number          int NOT NULL,
    current_approvers_number  int NOT NULL
);
