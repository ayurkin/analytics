-- +goose Up

CREATE TABLE IF NOT EXISTS analytics.event_uuid
(
    id uuid UNIQUE
);
