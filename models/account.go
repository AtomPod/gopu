package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

//User user models
type User struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`

	Username string `gorm:"column:username;unique_key;index" json:"username"`
	Password string `gorm:"column:password" json:"-"`
	Email    string `gorm:"column:email;unique_key;index" json:"email"`

	Profile   Profile
	ProfileID uint

	Roles string `gorm:"-" json:"roles"`
}

//BeforeCreate for gorm set id
func (u *User) BeforeCreate(s *gorm.Scope) error {
	return s.SetColumn("id", xid.New().String())
}

//Profile profile of user
type Profile struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
	Avatar    string     `gorm:"column:avatar" json:"avatar"`
	Nickname  string     `gorm:"column:nickname" json:"nickname"`
	Company   string     `gorm:"column:company" json:"company"`
	Location  string     `gorm:"column:location" json:"location"`
}
