-- name: CreatePost :one
insert into posts (id, created_at, updated_at, last_fetched_at, title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: GetPosts :many
select * from posts
order by last_fetched_at asc nulls first, published_at desc nulls last
limit $1;


-- name: MarkPostFetched :exec
update posts
set updated_at = $2, last_fetched_at = $3
where id = $1;
