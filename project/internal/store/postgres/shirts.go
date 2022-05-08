package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
)

func (db *DB) Shirts() store.ShirtsRepository {
	if db.shirtbd == nil {
		db.shirtbd = NewShirtsRepository(db.conn)
	}

	return db.shirtbd
}

type ShirtsRepository struct {
	conn *sqlx.DB
}

func NewShirtsRepository(conn *sqlx.DB) store.ShirtsRepository {
	return &ShirtsRepository{conn: conn}
}

func (s ShirtsRepository) Create(ctx context.Context, shirt *models.Shirt) error {
	_, err := s.conn.Exec("INSERT INTO shirtbd(name, color, size, condition, description, phonenumber, price) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		shirt.Name,
		shirt.Color,
		shirt.Size,
		shirt.Condition,
		shirt.Description,
		shirt.Phonenumber,
		shirt.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s ShirtsRepository) All(ctx context.Context, filter *models.ShirtsFilter) ([]*models.Shirt, error) {
	shirtbd := make([]*models.Shirt, 0)
	if err := s.conn.Select(&shirtbd, "SELECT * FROM shirtbd"); err != nil {
		return nil, err
	}

	return shirtbd, nil
}

func (s ShirtsRepository) ByID(ctx context.Context, id int) (*models.Shirt, error) {
	shirt := new(models.Shirt)
	if err := s.conn.Get(shirt, "SELECT id, name, size, condition, description, phonenumber, price  FROM shirtbd WHERE id=$1", id); err != nil {
		return nil, err
	}

	return shirt, nil
}

func (s ShirtsRepository) Update(ctx context.Context, shirt *models.Shirt) error {
	_, err := s.conn.Exec("UPDATE shirtbd SET name = $1, color = $2, size = $3, condition = $4, description = $5, phonenumber = $6, price = $7 WHERE id = $8",
		shirt.Name,
		shirt.Color,
		shirt.Size,
		shirt.Condition,
		shirt.Description,
		shirt.Phonenumber,
		shirt.Price,
		shirt.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s ShirtsRepository) Delete(ctx context.Context, id int) error {
	_, err := s.conn.Exec("DELETE FROM shirtbd WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
