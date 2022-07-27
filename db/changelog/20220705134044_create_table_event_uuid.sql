-- +goose Up

CREATE TABLE IF NOT EXISTS analytics.event_uuid
(
    uuid uuid UNIQUE
);
