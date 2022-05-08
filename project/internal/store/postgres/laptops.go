package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
)

func (db *DB) Laptops() store.LaptopsRepository {
	if db.laptopsbd == nil {
		db.laptopsbd = NewLaptopsRepository(db.conn)
	}

	return db.laptopsbd
}

type LaptopsRepository struct {
	conn *sqlx.DB
}

func NewLaptopsRepository(conn *sqlx.DB) store.LaptopsRepository {
	return &LaptopsRepository{conn: conn}
}

func (l LaptopsRepository) Create(ctx context.Context, laptop *models.Laptop) error {
	_, err := l.conn.Exec("INSERT INTO laptopsbd(name, year, display, memory, storage, condition, description, phonenumber, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		laptop.Name,
		laptop.Year,
		laptop.Display,
		laptop.Memory,
		laptop.Storage,
		laptop.Condition,
		laptop.Description,
		laptop.Phonenumber,
		laptop.Price)
	if err != nil {
		return err
	}

	return nil
}

func (l LaptopsRepository) All(ctx context.Context, filter *models.LaptopsFilter) ([]*models.Laptop, error) {
	laptopsbd := make([]*models.Laptop, 0)
	basicQuery := "SELECT * FROM laptopsbd"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := l.conn.Select(&laptopsbd, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return laptopsbd, nil
	}

	if err := l.conn.Select(&laptopsbd, "SELECT * FROM laptopsbd"); err != nil {
		return nil, err
	}

	return laptopsbd, nil
}

func (l LaptopsRepository) ByID(ctx context.Context, id int) (*models.Laptop, error) {
	laptop := new(models.Laptop)
	if err := l.conn.Get(laptop, "SELECT id, name, year, display, memory, storage, condition, description, phonenumber, price  FROM laptopsbd WHERE id=$1", id); err != nil {
		return nil, err
	}

	return laptop, nil
}

func (l LaptopsRepository) Update(ctx context.Context, laptop *models.Laptop) error {
	_, err := l.conn.Exec("UPDATE laptopsbd SET name = $1, year = $2, display = $3, memory = $4, storage = $5, condition = $6, description = $7, phonenumber = $8, price = $9 WHERE id = $10",
		laptop.Name,
		laptop.Year,
		laptop.Display,
		laptop.Memory,
		laptop.Storage,
		laptop.Condition,
		laptop.Description,
		laptop.Phonenumber,
		laptop.Price,
		laptop.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l LaptopsRepository) Delete(ctx context.Context, id int) error {
	_, err := l.conn.Exec("DELETE FROM laptopsbd WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
