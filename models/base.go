package models

import "time"

type Base struct {
	ID int64 `gorm:"column:id;PRIMARY_KEY;AUTO_INCREMENT;not null" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}