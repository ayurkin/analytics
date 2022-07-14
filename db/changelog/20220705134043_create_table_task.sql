-- +goose Up

create table analytics.task
(
    id                        int,
    status                    text,
    created_at                timestamptz,
    last_mail_at              timestamptz,
    total_time                interval,
    approvers_number          int,
    current_approvers_number  int
);
