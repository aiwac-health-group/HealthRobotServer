package models

type ClientInfo struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique"`
	ClientName string `gorm:"column:client_name;type:varchar(128);not null"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
}

