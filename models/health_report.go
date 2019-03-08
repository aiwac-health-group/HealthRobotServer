package models

type HealthReport struct {
	Base
	UserAccount 	string		`gorm:"column:user_account;type:varchar(11);not null"`
	Report		    string		`gorm:"column:report;type:varchar(200)"`
}

//测肤结果
type SkinTest struct {
	Base
	UserAccount 	string		  `gorm:"column:client_account;type:varchar(11);not null"`
	SkinDesc		string		  `gorm:"column:skin_desc;type:varchar(300);"`
	FaceURL			string		  `gorm:"column:skin_url;type:mediumtext;"`
}