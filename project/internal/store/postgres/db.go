package postgres

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"lectures-6/internal/store"
)

type DB struct {
	conn *sqlx.DB
	laptopsbd store.LaptopsRepository
	snowboardbd store.SnowboardsRepository
	shirtbd store.ShirtsRepository
	toybd store.ToysRepository
	userbd store.UsersRepository
}



func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		panic(err)
	}

	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}




