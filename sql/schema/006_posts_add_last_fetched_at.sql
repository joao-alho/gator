-- +goose Up
alter table posts add column last_fetched_at timestamp;


-- +goose Down
alter table posts drop column last_fetched_at;
