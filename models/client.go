package models

//管理员及客服基本信息
type ServiceInfo struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
	ClientName string `gorm:"column:client_name;type:varchar(128);"`
	OnlineStatus string `gorm:"column:online_status;default:1;type:varchar(1)"` //1表示离线，2表示在线，3表示在忙
}

//医生账号信息表
type DoctorInfo struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
	ClientName string `gorm:"column:client_name;type:varchar(128);"`
	Department string `gorm:"column:department;type:varchar(12)"` //科室
	Brief string  `gorm:"column:brief;type:varchar(128)"` //简介
	OnlineStatus string `gorm:"column:online_status;default:1;type:varchar(1)"` //1表示离线，2表示在线，3表示在忙
}
//医生列表项视图
type DoctorItem struct {
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique" json:"account"`
	ClientName string `gorm:"column:client_name;type:varchar(128)" json:"name"`
	Department string `gorm:"column:department;type:varchar(12)" json:"class"` //医生的科室
}

//机器人用户基本信息
type RobotInfo struct {
	Base
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;"`
	ClientPassword string `gorm:"column:client_password;type:varchar(128)"`
	ClientName string `gorm:"column:client_name;type:varchar(128);"`
	ClientType string `gorm:"column:client_type;type:varchar(7);not null"`
	Sex string `gorm:"column:sex;type:char(1)"`
	Birthday string `gorm:"column:birthday;type:char(8)"`
	Address string `gorm:"column:address;type:varchar(128)"`
	Wechat string `gorm:"column:wechat;type:varchar(128)"`
	OnlineStatus string `gorm:"column:online_status;default:1;type:varchar(1)"` //1表示离线，2表示在线，3表示在忙
	HealthStatus string `gorm:"column:health_status;default:1;type:varchar(1)"` //1表示需要为用户填写报告，2表示不需要
}

//机器人列表项视图
type RobotItem struct {
	ClientAccount string `json:"clientID"`
}


//客服人员工号姓名列表项视图
type Client struct{
	ClientAccount string `gorm:"column:client_account;type:varchar(11);not null;unique;" json:"acount"`
	ClientName string `gorm:"column:client_name;type:varchar(128);" json:"name"`
}