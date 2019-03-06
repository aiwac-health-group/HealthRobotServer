package models


//客服人员健康讲座工作量
type ClientLecture struct{
	Client
	CountLecture int  
}

//客服人员健康讲座工作量
type ClientReport struct{
	Client
	CountReport int  
}

//客服人员健康讲座工作量
type ClientRegist struct{
	Client
	CountRegist int  
}

type ClientTotalWork struct{
	Client
	CountLecture int `json:"lecturesNum"`
	CountReport int  `json:"registersNum"`
	CountRegist int  `json:"recommendNum"`
}

//客服人员健康讲座工作量
type StatisticRequest struct{
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"` 
}