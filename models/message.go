package models

//定义了挂号和问诊处理的等消息模型

//问诊请求模型
type TreatInfo struct {
	Base
	Account string `gorm:"column:client_account;type:varchar(11);not null;" json:"account"`
	ClientName string `gorm:"column:client_name;type:varchar(128);" json:"name"`
	HandleDoctor string `gorm:"column:handle_doctor;type:varchar(11);default:-" json:"-"`
	Others string `gorm:"column:others;type:varchar(128);" json:"others"`
	Status string `gorm:"column:treat_status;default:0;type:varchar(1)" json:"-"` //问诊请求的处理状态，1表示未处理，2表示正在处理，3表示处理完成
}

//问诊医生分配模型
type TreatAllocation struct {
	Patient string `json:"userAccount"`
	Doctor string `json:"doctorAccount"`
}


//挂号请求模型
type Registration struct {
	Base
	UserAccount 	string		`gorm:"column:user_account;type:varchar(11);not null"`
	Status			string		`gorm:"column:regist_status;default:1;type:varchar(1);not null"`
	RegistDesc		string		`gorm:"column:regist_descrbtion;type:varchar(50)"`
	HandleSevice	string		`gorm:"column:service_account;type:varchar(11)"`
	Regist
}

type Regist struct {
	Province		string		`gorm:"column:province;type:varchar(50);not null" json:"province,omitempty"`
	City			string		`gorm:"column:city;type:varchar(50);not null" json:"city,omitempty"`
	Hospital		string		`gorm:"column:hospital;type:varchar(50);not null" json:"hospital,omitempty"`
	Department 	 	string		`gorm:"column:department;type:varchar(50);not null" json:"department,omitempty"`
}

type RegistRequest struct {
	UserAccount 	string		`json:"user_account"`
	Request			Regist		`json:"regist_request"`
}

type RegistProcessRequest struct {
	ID				int64		`json:"id"`
	UserAccount 	string		`json:"user_account"`
	HandleSevice	string		`json:"service_account"`
}

type RegistProcessResponse struct {
	ID				int64		`json:"id"`
	ResponseDesc 	string		`json:"result"`
}

type RegistProcessingInfo struct {
	BaseResponse
	ID			 int64		`json:"id"`
	UserName	 string		`json:"userName"`
	UserAccoutn	 string		`json:"userAccount"`
	Class		 string		`json:"class"`
	Others		 string		`json:"others"`
}


