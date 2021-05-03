package photo

import (
	"context"
	"database/sql"
)

const createTableIfNotExists = `
CREATE TABLE IF NOT EXISTS photos
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    INT8 REFERENCES users (id) NOT NULL,
    url        TEXT                       NOT NULL,
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
