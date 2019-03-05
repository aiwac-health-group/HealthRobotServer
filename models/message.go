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

//挂号请求模型
type Registration struct {

}

//问诊医生分配模型
type TreatAllocation struct {
	Patient string `json:"userAccount"`
	Doctor string `json:"doctorAccount"`
}