package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `json:"DeletedAt" gorm:"index"`

	ID       uuid.UUID `json:"ID"       gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username string    `json:"Username" gorm:"uniqueIndex;not null"`
	Email    string    `json:"Email"    gorm:"uniqueIndex;not null"`
	Password string    `json:"Password" gorm:"not null"`
	Role     string    `json:"Role"     gorm:"default:user;not null"`
	Verified bool      `json:"Verified" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
