-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeed :one
select * from feeds where name = $1;

-- name: ResetFeeds :exec
delete from feeds;

-- name: GetFeeds :many
select * from feeds;
