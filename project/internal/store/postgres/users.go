package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
)

func (db *DB) Users() store.UsersRepository {
	if db.userbd == nil {
		db.userbd = NewUsersRepository(db.conn)
	}

	return db.userbd
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UsersRepository {
	return &UsersRepository{conn: conn}
}

func (u UsersRepository) Create(ctx context.Context, users *models.User) error {
	_, err := u.conn.Exec("INSERT INTO userbd(username, email, password) VALUES ($1, $2, $3)",
		users.Username,
		users.Email,
		users.Password)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) All(ctx context.Context, filter *models.UsersFilter) ([]*models.User, error) {
	userbd := make([]*models.User, 0)
	if err := u.conn.Select(&userbd, "SELECT * FROM userbd"); err != nil {
		return nil, err
	}

	return userbd, nil
}


func (u UsersRepository) ByID(ctx context.Context, email string) (*models.User, error) {
	panic("implement me")
}

