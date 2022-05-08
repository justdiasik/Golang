package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
)

func (db *DB) Toys() store.ToysRepository {
	if db.toybd == nil {
		db.toybd = NewToysRepository(db.conn)
	}

	return db.toybd
}

type ToysRepository struct {
	conn *sqlx.DB
}

func NewToysRepository(conn *sqlx.DB) store.ToysRepository {
	return &ToysRepository{conn: conn}
}

func (t ToysRepository) Create(ctx context.Context, toys *models.Toy) error {
	_, err := t.conn.Exec("INSERT INTO toybd(name, condition, description, phonenumber, price) VALUES ($1, $2, $3, $4, $5)",
		toys.Name,
		toys.Condition,
		toys.Description,
		toys.Phonenumber,
		toys.Price)
	if err != nil {
		return err
	}

	return nil
}

func (t ToysRepository) All(ctx context.Context, filter *models.ToysFilter) ([]*models.Toy, error) {
	toybd := make([]*models.Toy, 0)
	if err := t.conn.Select(&toybd, "SELECT * FROM toybd"); err != nil {
		return nil, err
	}

	return toybd, nil
}

func (t ToysRepository) ByID(ctx context.Context, id int) (*models.Toy, error) {
	toy := new(models.Toy)
	if err := t.conn.Get(toy, "SELECT id, name, condition, description, phonenumber, price  FROM toybd WHERE id=$1", id); err != nil {
		return nil, err
	}

	return toy, nil
}

func (t ToysRepository) Update(ctx context.Context, toy *models.Toy) error {
	_, err := t.conn.Exec("UPDATE toybd SET name = $1, condition = $2, description = $3, phonenumber = $4, price = $5 WHERE id = $6",
		toy.Name,
		toy.Condition,
		toy.Description,
		toy.Phonenumber,
		toy.Price,
		toy.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t ToysRepository) Delete(ctx context.Context, id int) error {
	_, err := t.conn.Exec("DELETE FROM toybd WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
