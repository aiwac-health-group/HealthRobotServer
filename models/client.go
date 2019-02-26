package models

type ClientInfo struct {
	Base
	ClientName string `gorm:"column:client_name;type:varchar(128);not null;unique"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)"`
	ClientType string `gorm:"column:client_type;type:varchar(6);not null"`
}
