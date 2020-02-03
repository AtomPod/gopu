package database

import (
	"context"
	"time"

	"github.com/ngs24313/gopu/models"
	"github.com/ngs24313/gopu/utils/database/database"
)

//UserListQuery query params for list user
type UserListQuery struct {
	Query           string
	Offset          int
	Count           int
	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	WithCount       bool

	OrderBy  []string
	SortType []string
}

//UserListResult query result for list user
type UserListResult struct {
	Count int64
	Users []*models.User
}

//AccountDatabase account database
type AccountDatabase interface {
	database.Database

	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, u *models.User) error
	UpdateProfile(ctx context.Context, p *models.Profile) error

	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUser(ctx context.Context, q UserListQuery) (*UserListResult, error)

	UserIsExists(ctx context.Context, u *models.User) (bool, error)

	CountUser(ctx context.Context) (int64, error)
}
