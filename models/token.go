package models

type Token struct {
	Base
	RawToken string `gorm:"column:raw_token;varchar(128);not null;"`
	ClientID int64 `gorm:"column:client_id;not null;"`
	ClientType string `gorm:"column:client_type;not null"`
	ExpressIn int64 `gorm:"column:express_in"`
}