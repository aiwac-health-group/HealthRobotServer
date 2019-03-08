package models


//客服人员健康讲座工作量
type ClientLecture struct{
	ClientAccount string
	ClientName string
	CountLecture int
}

//客服人员健康讲座工作量
type ClientReport struct{
	ClientAccount string
	ClientName string
	CountReport int  
}

//客服人员健康讲座工作量
type ClientRegist struct{
	ClientAccount string
	ClientName string
	CountRegist int  
}

type ClientTotalWork struct{
	ClientAccount string `json:"acount"`
	ClientName string `json:"name"`
	CountLecture int `json:"lecturesNum"`
	CountReport int  `json:"registersNum"`
	CountRegist int  `json:"recommendNum"`
}

