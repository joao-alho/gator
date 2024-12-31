-- name: CreateFeedFollow :one
with inserted_feed_follow as (insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
values ($1, $2, $3, $4, $5)
returning *
)
	select 
		ff.*,
		f.name feed_name,
		u.name user_name
	from inserted_feed_follow ff
	inner join feeds f
	on f.id = ff.feed_id
	inner join users u
	on u.id = ff.user_id
;

-- name: GetFeedFollowsForUser :many
select 
	ff.*,
	f.name feed_name,
	u.name user_name
from feed_follows ff
inner join feeds f
on f.id = ff.feed_id
inner join users u
on u.id = ff.user_id
where ff.user_id = $1
;

-- name: DeleteFeedFollow :exec
delete from feed_follows using feeds
  where feed_follows.user_id = $1 
    and feed_follows.feed_id = feeds.id
	and feeds.url = $2;
