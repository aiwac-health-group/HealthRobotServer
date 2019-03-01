package models

type ClientInfo struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique" json:"account"`
	ClientName string `gorm:"column:client_name;type:varchar(128);not null" json:"name"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)" json:"-"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null" json:"-"`
	OnlineStatus string `gorm:"column:online_status;default:0;type:varchar(1)" json:"-"` //1表示离线，2表示在线，3表示在忙
	Department string `gorm:"column:department;type:varchar(12)" json:"class"`
}

