package reaction

import (
	"context"
	"database/sql"
	"unlinked/entities"
)

const createTablesIfNotExists = `
CREATE TABLE IF NOT EXISTS reactions
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT UNIQUE               NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS photos_reactions
(
    photo_id    INT8 REFERENCES photos (id)    NOT NULL,
    user_id     INT8 REFERENCES users (id)     NOT NULL,
    reaction_id INT8 REFERENCES reactions (id) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now()      NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT now()      NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS photos_reactions_uidx ON photos_reactions (photo_id, user_id, reaction_id);
`

func (r *reactionService) CreateTableIfNotExists(ctx context.Context) (err error) {
	if _, err = r.db.ExecContext(ctx, createTablesIfNotExists); err != nil {
		return err
	}
	return nil
}

const insertReaction = `
INSERT INTO reactions(name)
VALUES ($1)
RETURNING id`

func (r *reactionService) Add(ctx context.Context, reactionName string) (id int64, err error) {
	var rows *sql.Rows
	if rows, err = r.db.QueryContext(ctx, insertReaction, reactionName); err != nil {
		return 0, err
	}

	if err = rows.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

const insertReactionsPhotos = `
INSERT INTO photos_reactions (photo_id, user_id, reaction_id)
VALUES ($1, $2, $3);
`

func (r *reactionService) AddReactionToPhoto(ctx context.Context, photoID, reactionID, userID int64) (err error) {
	if _, err = r.db.QueryContext(ctx, insertReactionsPhotos, photoID, userID, reactionID); err != nil {
		return err
	}

	return nil
}

const selectReactionsCountByPhotoID = `
SELECT r.id, r.name, count(*)
FROM photos_reactions AS pr
         INNER JOIN reactions AS r ON pr.reaction_id = r.id
WHERE photo_id = $1
GROUP BY r.id, r.name;
`

func (r *reactionService) GetAllReactionsCountByPhotoID(
	ctx context.Context,
	photoID int64,
) (_ []*entities.ReactionInfo, err error) {
	var rows *sql.Rows
	if rows, err = r.db.QueryContext(ctx, selectReactionsCountByPhotoID, photoID); err != nil {
		return nil, err
	}

	var reactionsInfo = make([]*entities.ReactionInfo, 0, 100)
	for rows.Next() {
		tmp := &entities.ReactionInfo{}
		if err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.Count); err != nil {
			return nil, err
		}
		reactionsInfo = append(reactionsInfo, tmp)
	}

	return reactionsInfo, nil
}

const reactionsCountByPhotoID = `
SELECT count(*)
FROM photos_reactions
WHERE photo_id = $1
  AND reaction_id = $2;
`

func (r *reactionService) GetReactionsCountByPhotoID(
	ctx context.Context,
	photoID, reactionID int64,
) (count int64, err error) {
	var rows *sql.Rows
	if rows, err = r.db.QueryContext(ctx, reactionsCountByPhotoID, photoID, reactionID); err != nil {
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
