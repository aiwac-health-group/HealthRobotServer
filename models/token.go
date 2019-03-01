package models

type Token struct {
	Base
	RawToken string `gorm:"column:raw_token;varchar(256);not null;unique"`
	ClientAccount string `gorm:"column:client_account;varchar(11);not null;unique"`
	ClientType string `gorm:"column:client_type;varchar(7);not null"`
	ExpressIn int64 `gorm:"column:express_in"`
}