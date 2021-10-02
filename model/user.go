package model

import (
	"time"
)

type UserRepository interface {
	CreateUser(username string) (int64, error)
	GetUsers() ([]User, error)
	GetUserByID(id int64) (User, error)
	GetUserByPublicID(id string) (User, error)
	UpdateUserByID(user User) (int64, error)
	DeleteUserByID(id int64) (int64, error)
}

// User ...
type User struct {
	ID       int64  `db:"id"`
	PublicID string `db:"public_id"`
	// profile
	Username string `db:"username"`

	UpdatedAt time.Time  `db:"updated_at"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// TODO: replace query builder(goqu)
const (
	insertUserQuery = `
		INSERT users(public_id, username)
		VALUE (?, ?);
	`

	selectUsersQuery = `
		SELECT * FROM users
		WHERE deleted_at IS NULL;
	`

	selectUserByIDQuery = `
	SELECT * FROM users
	WHERE deleted_at IS NULL AND id = ?;
	`

	selectUserByPublicIDQuery = `
	SELECT * FROM users
	WHERE deleted_at IS NULL AND public_id = ?;
	`

	updateUserByIDQuery = `
		UPDATE users
			SET username = :username
			WHERE id = :id;
	`

	deleteUserQuery = `
		UPDATE users
			SET deleted_at = NOW()
			WHERE id = ?
	`
)

func (repo *SqlxRepository) CreateUser(username string) (int64, error) {
	publicID, _ := GenerateRandomString(10)
	result, err := repo.db.Exec(insertUserQuery, publicID, username)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()

}
func (repo *SqlxRepository) GetUsers() ([]User, error) {
	var users []User
	err := repo.db.Select(&users, selectUsersQuery)
	return users, err

}
func (repo *SqlxRepository) GetUserByID(id int64) (User, error) {
	user := User{}
	err := repo.db.Get(&user, selectUserByIDQuery, id)
	return user, err
}

func (repo *SqlxRepository) GetUserByPublicID(id string) (User, error) {
	user := User{}
	err := repo.db.Get(&user, selectUserByPublicIDQuery, id)
	return user, err
}

func (repo *SqlxRepository) UpdateUserByID(user User) (int64, error) {
	result, err := repo.db.NamedExec(updateUserByIDQuery, user)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}
func (repo *SqlxRepository) DeleteUserByID(id int64) (int64, error) {
	result, err := repo.db.Exec(deleteUserQuery, id)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()

}
