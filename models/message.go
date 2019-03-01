package models

//定义了挂号请求和问诊请求等消息模型

//问诊请求模型
type Treat struct {
	Base
	Account string
	ClientName string
	Others string
	//问诊请求的处理状态，0表示未处理，1表示正在处理，2表示处理完成
	Status int8
}

//挂号请求模型
type Registration struct {

}