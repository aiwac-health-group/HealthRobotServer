package models

//定义了挂号和问诊处理的等消息模型

//问诊请求模型
type TreatInfo struct {
	Base
	Account string
	ClientName string
	Others string
	//问诊请求的处理状态，0表示未处理，1表示正在处理，2表示处理完成
	Status string
}

//挂号请求模型
type Registration struct {

}

//问诊医生分配模型
type TreatAllocation struct {
	Patient string `json:"userAccount"`
	Doctor string `json:"doctorAccount"`
}