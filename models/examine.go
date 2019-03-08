package models

// 数据库中的体检推荐表
type PhysicalExamine struct {
	Base
	Title			string		`gorm:"column:title;type:varchar(20);not null"`
	Abstract		string		`gorm:"column:abstract;type:varchar(128)"`
	Infor			string		`gorm:"column:infor;type:varchar(500);not null"`
	HandleSevice	string		`gorm:"column:service_account;type:varchar(11)"`
	Cover			string		`gorm:"column:cover;type:mediumtext"`
}


