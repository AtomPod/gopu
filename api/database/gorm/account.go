package gorm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	dao "github.com/ngs24313/gopu/api/database"
	"github.com/ngs24313/gopu/models"
	gormdb "github.com/ngs24313/gopu/utils/database/gorm"
)

//AccountDatabase account database
type AccountDatabase struct {
	gormdb.Database
}

func (d *AccountDatabase) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	db := d.Instance()
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *AccountDatabase) DeleteUser(ctx context.Context, id string) error {
	db := d.Instance()
	db = db.Where("id = ?", id).Delete(&models.User{})
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected == 0 {
		return dao.ErrNotFound
	}
	return nil
}

func (d *AccountDatabase) UpdateUser(ctx context.Context, u *models.User) error {
	db := d.Instance()
	db = db.Save(u)
	return db.Error
}

func (d *AccountDatabase) UpdateProfile(ctx context.Context, p *models.Profile) error {
	db := d.Instance()

	if err := db.Save(p).Error; err != nil {
		return err
	}
	return nil
}

func (d *AccountDatabase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	db := d.Instance()

	var user models.User
	if err := db.Preload("Profile").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dao.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *AccountDatabase) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	db := d.Instance()

	var user models.User
	if err := db.Preload("Profile").Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dao.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *AccountDatabase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	db := d.Instance()

	var user models.User
	if err := db.Preload("Profile").Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dao.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *AccountDatabase) UserIsExists(ctx context.Context, u *models.User) (bool, error) {
	db := d.Instance()

	var user models.User
	if err := db.Where("username = ? OR email = ?", u.Username, u.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if user.Username == u.Username {
		return true, dao.ErrUsernameAlreadyExists
	}
	return true, dao.ErrEmailAlreadyExists
}

func (d *AccountDatabase) CountUser(ctx context.Context) (int64, error) {
	db := d.Instance()

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *AccountDatabase) ListUser(ctx context.Context, q dao.UserListQuery) (*dao.UserListResult, error) {
	var result dao.UserListResult

	db := d.Instance()
	db = db.Preload("Profile")
	db = db.Table("users")
	if q.Query != "" {
		likeString := "%" + q.Query + "%"
		db = db.Where("(username LIKE ?) OR (email LIKE ?)", likeString, likeString)
	}

	if !q.CreateTimeStart.IsZero() {
		db = db.Where("createdAt >= ?", q.CreateTimeStart)
	}

	if !q.CreateTimeEnd.IsZero() {
		db = db.Where("createdAt < ?", q.CreateTimeEnd)
	}

	if q.WithCount {
		if err := db.Count(&result.Count).Error; err != nil {
			return nil, err
		}
	}

	if q.Offset > 0 {
		db = db.Offset(q.Offset)
	}

	var limit int = 20
	if q.Count > 0 {
		limit = q.Count
	}

	db = db.Limit(limit)

	for i, order := range q.OrderBy {
		sort := "asc"

		if i < len(q.SortType) {
			if q.SortType[i] == "desc" {
				sort = q.SortType[i]
			}
		}

		if d.canISort(order) {
			sort = strings.ToUpper(sort)
			db = db.Order(fmt.Sprintf("\"%s\" %s", order, sort))
		}
	}

	if len(q.OrderBy) == 0 {
		db = db.Order("id ASC")
	}

	if err := db.Find(&result.Users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dao.ErrNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (d *AccountDatabase) canISort(name string) bool {
	switch name {
	case "id", "username", "email", "createdAt", "updatedAt":
		return true
	default:
		return false
	}
}
