// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user_relate_url.sql

package db

import (
	"context"
	"time"
)

const createUserUrl = `-- name: CreateUserUrl :one
INSERT INTO user_relate_url (
  user_id, short_url, origin_url, status, expire_at
) VALUES (
  $1, $2, $3, 0, $4
)
RETURNING id, user_id, short_url, origin_url, status, expire_at, created_at
`

type CreateUserUrlParams struct {
	UserID    string    `json:"user_id"`
	ShortUrl  string    `json:"short_url"`
	OriginUrl string    `json:"origin_url"`
	ExpireAt  time.Time `json:"expire_at"`
}

func (q *Queries) CreateUserUrl(ctx context.Context, arg *CreateUserUrlParams) (*UserRelateUrl, error) {
	row := q.db.QueryRowContext(ctx, createUserUrl,
		arg.UserID,
		arg.ShortUrl,
		arg.OriginUrl,
		arg.ExpireAt,
	)
	var i UserRelateUrl
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ShortUrl,
		&i.OriginUrl,
		&i.Status,
		&i.ExpireAt,
		&i.CreatedAt,
	)
	return &i, err
}

const listUrlByUser = `-- name: ListUrlByUser :many
SELECT id, user_id, short_url, origin_url, status, expire_at, created_at
FROM user_relate_url
WHERE user_id = $1
ORDER BY id DESC
`

func (q *Queries) ListUrlByUser(ctx context.Context, userID string) ([]*UserRelateUrl, error) {
	rows, err := q.db.QueryContext(ctx, listUrlByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*UserRelateUrl{}
	for rows.Next() {
		var i UserRelateUrl
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ShortUrl,
			&i.OriginUrl,
			&i.Status,
			&i.ExpireAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateStatus = `-- name: UpdateStatus :one
UPDATE user_relate_url
SET status = $2
WHERE short_url = $1
RETURNING id, user_id, short_url, origin_url, status, expire_at, created_at
`

type UpdateStatusParams struct {
	ShortUrl string `json:"short_url"`
	Status   int32  `json:"status"`
}

func (q *Queries) UpdateStatus(ctx context.Context, arg *UpdateStatusParams) (*UserRelateUrl, error) {
	row := q.db.QueryRowContext(ctx, updateStatus, arg.ShortUrl, arg.Status)
	var i UserRelateUrl
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ShortUrl,
		&i.OriginUrl,
		&i.Status,
		&i.ExpireAt,
		&i.CreatedAt,
	)
	return &i, err
}
