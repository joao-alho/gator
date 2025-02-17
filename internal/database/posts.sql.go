// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
insert into posts (id, created_at, updated_at, last_fetched_at, title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning id, created_at, updated_at, title, url, description, published_at, feed_id, last_fetched_at
`

type CreatePostParams struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt sql.NullTime
	Title         string
	Url           string
	Description   sql.NullString
	PublishedAt   sql.NullTime
	FeedID        uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.LastFetchedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
select id, created_at, updated_at, title, url, description, published_at, feed_id, last_fetched_at from posts
order by last_fetched_at asc nulls first, published_at desc nulls last
limit $1
`

func (q *Queries) GetPosts(ctx context.Context, limit int32) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markPostFetched = `-- name: MarkPostFetched :exec
update posts
set updated_at = $2, last_fetched_at = $3
where id = $1
`

type MarkPostFetchedParams struct {
	ID            uuid.UUID
	UpdatedAt     time.Time
	LastFetchedAt sql.NullTime
}

func (q *Queries) MarkPostFetched(ctx context.Context, arg MarkPostFetchedParams) error {
	_, err := q.db.ExecContext(ctx, markPostFetched, arg.ID, arg.UpdatedAt, arg.LastFetchedAt)
	return err
}
