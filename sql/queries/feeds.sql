-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeed :one
select * from feeds where name = $1;

-- name: GetFeedFromUrl :one
select * from feeds where url = $1;

-- name: GetFeeds :many
select f.*, u.name as user_name from feeds f
inner join users u
on f.user_id = u.id;

-- name: MarkFeedFetched :exec
update feeds
set updated_at = $2, last_fetched_at = $3
where id = $1;

-- name: GetNextFeedToFetch :one
select *
from feeds
order by last_fetched_at asc nulls first
limit 1;
