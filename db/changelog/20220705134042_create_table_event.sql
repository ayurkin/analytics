-- +goose Up

create table analytics.event
(
    task_id          int,
    occurred_at      timestamptz,
    event_type       text,
    event_user       text,
    approvers_number int
);
