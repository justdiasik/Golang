package store

import (
	"context"
	"lectures-6/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error

	Laptops() LaptopsRepository

	Snowboards() SnowboardsRepository
	Shirts() ShirtsRepository
	Toys() ToysRepository
	Users() UsersRepository
}


type LaptopsRepository interface {
	Create(ctx context.Context, laptop *models.Laptop) error
	All(ctx context.Context, filter *models.LaptopsFilter) ([]*models.Laptop, error)
	ByID(ctx context.Context, id int) (*models.Laptop, error)
	Update(ctx context.Context, laptop *models.Laptop) error
	Delete(ctx context.Context, id int) error
}

type SnowboardsRepository interface {
	Create(ctx context.Context, snowboard *models.Snowboard) error
	All(ctx context.Context, filter *models.SnowboardsFilter) ([]*models.Snowboard, error)
	ByID(ctx context.Context, id int) (*models.Snowboard, error)
	Update(ctx context.Context, snowboard *models.Snowboard) error
	Delete(ctx context.Context, id int) error
}

type ShirtsRepository interface {
	Create(ctx context.Context, shirt *models.Shirt) error
	All(ctx context.Context, filter *models.ShirtsFilter) ([]*models.Shirt, error)
	ByID(ctx context.Context, id int) (*models.Shirt, error)
	Update(ctx context.Context, shirt *models.Shirt) error
	Delete(ctx context.Context, id int) error
}

type ToysRepository interface {
	Create(ctx context.Context, toy *models.Toy) error
	All(ctx context.Context, filter *models.ToysFilter) ([]*models.Toy, error)
	ByID(ctx context.Context, id int) (*models.Toy, error)
	Update(ctx context.Context, toy *models.Toy) error
	Delete(ctx context.Context, id int) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context, filter *models.UsersFilter) ([]*models.User, error)
	ByID(ctx context.Context, email string) (*models.User, error)
}
