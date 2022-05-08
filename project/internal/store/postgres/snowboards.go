package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
)

func (db *DB) Snowboards() store.SnowboardsRepository {
	if db.snowboardbd == nil {
		db.snowboardbd = NewSnowboardsRepository(db.conn)
	}

	return db.snowboardbd
}

type SnowboardsRepository struct {
	conn *sqlx.DB
}

func NewSnowboardsRepository(conn *sqlx.DB) store.SnowboardsRepository {
	return &SnowboardsRepository{conn: conn}
}

func (s SnowboardsRepository) Create(ctx context.Context, snowboard *models.Snowboard) error {
	_, err := s.conn.Exec("INSERT INTO snowboardbd(name, size, condition, description, phonenumber, price) VALUES ($1, $2, $3, $4, $5, $6)",
		snowboard.Name,
		snowboard.Size,
		snowboard.Condition,
		snowboard.Description,
		snowboard.Phonenumber,
		snowboard.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s SnowboardsRepository) All(ctx context.Context, filter *models.SnowboardsFilter) ([]*models.Snowboard, error) {
	snowboardbd := make([]*models.Snowboard, 0)
	if err := s.conn.Select(&snowboardbd, "SELECT * FROM snowboardbd"); err != nil {
		return nil, err
	}

	return snowboardbd, nil
}

func (s SnowboardsRepository) ByID(ctx context.Context, id int) (*models.Snowboard, error) {
	snowboard := new(models.Snowboard)
	if err := s.conn.Get(snowboard, "SELECT id, name, size, condition, description, phonenumber, price  FROM snowboardbd WHERE id=$1", id); err != nil {
		return nil, err
	}

	return snowboard, nil
}

func (s SnowboardsRepository) Update(ctx context.Context, snowboard *models.Snowboard) error {
	_, err := s.conn.Exec("UPDATE snowboardbd SET name = $1, size = $2, condition = $3, description = $4, phonenumber = $5, price = $6 WHERE id = $7",
		snowboard.Name,
		snowboard.Size,
		snowboard.Condition,
		snowboard.Description,
		snowboard.Phonenumber,
		snowboard.Price,
		snowboard.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s SnowboardsRepository) Delete(ctx context.Context, id int) error {
	_, err := s.conn.Exec("DELETE FROM snowboardbd WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
