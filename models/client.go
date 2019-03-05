package models

//账号基本信息
type ClientInfo struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
	OnlineStatus string `gorm:"column:online_status;default:0;type:varchar(1)"` //1表示离线，2表示在线，3表示在忙
}

//账号详细信息
type WebClient struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique;"`
	ClientName string `gorm:"column:client_name;type:varchar(128);"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
	Department string `gorm:"column:department;type:varchar(12)"` //医生的科室
	Brief string  `gorm:"column:brief;type:varchar(128)"` //医生简介
}

//医生列表项视图
type DoctorItem struct {
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique" json:"account"`
	ClientName string `gorm:"column:client_name;type:varchar(128)" json:"name"`
	Department string `gorm:"column:department;type:varchar(12)" json:"class"` //医生的科室
}

type Robot struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique;"`
	ClientName string `gorm:"column:client_name;type:varchar(128);"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
	Sex string `gorm:"column:sex;type:char(1)"`
	Birthday string `gorm:"column:birthday;type:char(8)"`
	Address string `gorm:"column:address;type:varchar(128)"`
	Wechat string `gorm:"column:wechat;type:varchar(128)"`
}