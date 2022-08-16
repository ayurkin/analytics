-- +goose Up

ALTER TABLE analytics.task ADD CONSTRAINT unique_id UNIQUE (id);