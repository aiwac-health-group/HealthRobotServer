package dao

import "github.com/jinzhu/gorm"

type TokenDao struct {
	engin *gorm.DB
}

func NewTokenDao(engin *gorm.DB) *TokenDao {
	return &TokenDao{
		engin:engin,
	}
}

