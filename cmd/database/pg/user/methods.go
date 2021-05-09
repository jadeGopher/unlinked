package user

import (
	"context"
	"database/sql"
	"unlinked/entities"
)

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    avatar_url TEXT                      NOT NULL,
    name       TEXT                      NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);
CREATE TABLE IF NOT EXISTS friends_relationships
(
    user_id   INT8 REFERENCES users (id) NOT NULL,
    friend_id INT8 REFERENCES users (id) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS friends_relationships_uid ON friends_relationships (user_id, friend_id);`

// CreateTableIfNotExists using to init tables instead of normal migrations, because i have no time
func (u *userService) CreateTableIfNotExists(ctx context.Context) (err error) {
	if _, err = u.db.ExecContext(ctx, createUsersTable); err != nil {
		return err
	}
	return nil
}

const addUser = `
INSERT INTO users(name, avatar_url)
VALUES ($1, $2)
RETURNING id`

func (u *userService) Add(ctx context.Context, name, avatarURL string) (id int64, err error) {
	var rows *sql.Rows
	if rows, err = u.db.QueryContext(ctx, addUser, name, avatarURL); err != nil {
		return 0, err
	}

	if err = rows.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

const selectUserByID = `
SELECT u.id, u.name, u.avatar_url, u.created_at, u.updated_at
FROM users AS u WHERE u.id = $1;
`

func (u *userService) GetByID(ctx context.Context, userID int64) (_ *entities.User, err error) {
	var rows *sql.Rows
	if rows, err = u.db.QueryContext(ctx, selectUserByID, userID); err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var user = &entities.User{}
	if err = rows.Scan(&user.ID, &user.Name, &user.Avatar, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return user, nil
}

const getFriendsCountByID = `
SELECT count(*)
FROM friends_relationships as fr
         INNER JOIN users u ON fr.friend_id = u.id
WHERE fr.user_id = $1;
`

func (u *userService) GetFriendsCountByID(ctx context.Context, userID int64) (count int64, err error) {
	var rows *sql.Rows
	if rows, err = u.db.QueryContext(ctx, getFriendsCountByID, userID); err != nil {
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

const selectUserFriends = `
SELECT u.id, u.name, u.avatar_url, u.created_at, u.updated_at
FROM friends_relationships as fr
         INNER JOIN users u ON fr.friend_id = u.id
WHERE fr.user_id = $1
LIMIT $2 OFFSET $3;
`

func (u *userService) GetFriendsByID(ctx context.Context, userID, limit, offset int64) (_ []*entities.User, err error) {
	var rows *sql.Rows
	if rows, err = u.db.QueryContext(ctx, selectUserFriends, userID, limit, offset); err != nil {
		return nil, err
	}

	var users = make([]*entities.User, 0, limit)
	for rows.Next() {
		user := &entities.User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Avatar, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

const selectUsersByReactionID = `
SELECT users.id, users.avatar_url, users.name, users.created_at, users.updated_at
FROM users
         INNER JOIN photos_reactions pr ON users.id = pr.user_id
WHERE reaction_id = $1
  AND photo_id = $2 LIMIT $3 OFFSET $4;
`

func (u *userService) GetUsersByReactionIDUnderPhoto(
	ctx context.Context,
	photoID, reactionID, limit, offset int64,
) (_ []*entities.User, err error) {
	var rows *sql.Rows
	if rows, err = u.db.QueryContext(ctx, selectUsersByReactionID, reactionID, photoID, limit, offset); err != nil {
		return nil, err
	}

	var users = make([]*entities.User, 0, limit)
	for rows.Next() {
		user := &entities.User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Avatar, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
