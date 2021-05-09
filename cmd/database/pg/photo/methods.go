package photo

import (
	"context"
	"database/sql"
	"unlinked/entities"
)

const createTableIfNotExists = `
CREATE TABLE IF NOT EXISTS photos
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    INT8 REFERENCES users (id) NOT NULL,
    url        TEXT UNIQUE                NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()  NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now()  NOT NULL
);`

// CreateTableIfNotExists using to init tables instead of normal migrations, because i have no time
func (p *photoService) CreateTableIfNotExists(ctx context.Context) (err error) {
	if _, err = p.db.ExecContext(ctx, createTableIfNotExists); err != nil {
		return err
	}
	return nil
}

const addPhoto = `
INSERT INTO photos(user_id, url)
VALUES ($1, $2)
RETURNING id;
`

func (p *photoService) Add(ctx context.Context, userID int64, url string) (id int64, err error) {
	var rows *sql.Rows
	if rows, err = p.db.QueryContext(ctx, addPhoto, userID, url); err != nil {
		return 0, err
	}

	if err = rows.Scan(&id); err != nil {
		return 0, err
	}

	return 0, nil
}

const selectPhotosCountByUserID = `
SELECT count(*)
FROM photos
WHERE photos.user_id = $1;
`

func (p *photoService) GetPhotosCountByUserID(ctx context.Context, userID int64) (count int64, err error) {
	var rows *sql.Rows
	if rows, err = p.db.QueryContext(ctx, selectPhotosCountByUserID, userID); err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, sql.ErrNoRows
	}

	if err = rows.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

const selectPhotosByUserID = `
SELECT id, url, created_at, updated_at
FROM photos
WHERE photos.user_id = $1
LIMIT $2 OFFSET $3;
`

func (p *photoService) GetPhotosByUserID(
	ctx context.Context,
	userID, limit, offset int64,
) (_ []*entities.Photo, err error) {
	var rows *sql.Rows
	if rows, err = p.db.QueryContext(ctx, selectPhotosByUserID, userID, limit, offset); err != nil {
		return nil, err
	}

	var photos = make([]*entities.Photo, 0, limit)
	for rows.Next() {
		photo := &entities.Photo{}
		if err = rows.Scan(&photo.ID, &photo.Url, &photo.CreatedAt, &photo.UpdatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}
