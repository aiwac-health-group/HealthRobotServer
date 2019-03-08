package models

//http请求
//登录请求
type LoginRequest struct {
	Account string `json:"account"`
	Password string `json:"password"`
	IdentifyCode string `json:"identifyCode"`
}
//机器人获取token请求
type TokenGetRequest struct {
	Account string `json:"account"`
	OldToken string `json:"oldToken"`
}
//管理员增加账号请求
type AccountAddRequest struct {
	Account string `json:"account"`
	Name string `json:"name"`
	Password string `json:"pass"`
}
//管理员修改账号信息
type AccountModifyRequest struct {
	Account string `json:"account"`
	OperationType string `json:"type"`
	Value string `json:"value"`
}
//管理员修改医生信息请求
type DoctorProfileModifyRequest struct {
	Account string `json:"account"`
	Name string `json:"name"`
	Department string `json:"class"`
	Brief string `json:"blief"`
}
//医生获取用户测肤结果请求
type GetSkinRequest struct {
	Account string `json:"clientID"`
}
//客服体检推荐发布请求
type ExamineRequst struct {
	Title		string		`json:"title"`
	Abstract	string		`json:"blief"`
	Infor		string		`json:"text"`
}
//客服上传文本类讲座
type TextLectureUploadRequest struct{
	Title string `json:"title"`
	Blief string `json:"blief"`
	Text string `json:"text"`
}
//客服上传文件类讲座
type FileLectureUploadRequest struct{
	Title string `json:"title"`
	Filename string `json:"filename"`
	Blief string `json:"blief"`
}
//客服人员健康讲座工作量
type StatisticRequest struct{
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
}

//websocket请求
type WSRequest struct {
	BusinessCode string `json:"code"`
	WSRobotRequest
	WSWebRequest
}

//Web端的websocket请求字段
type WSWebRequest struct {
	Message string `json:"message,omitempty"`
}

//机器人端的websocket请求字段
type WSRobotRequest struct {
	//公共字字段
	Account    string `json:"account,omitempty"`
	UniqueID   string `json:"uniqueID,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	Time       string `json:"time,omitempty"`
	ExamID 		int64 	`json:"examID,omitempty"`
	RegisterID 	int64 	`json:"registerID,omitempty"`
	ReportID	int64	`json:"resultID,omitempty"`
	LectureID   string `json:"lectureID"`
	WSRobotProfile
	WSSkinRequest
	Regist
}
//详细信息修改请求
type WSRobotProfile struct {
	//个人信息
	Name string `json:"name"`
	Sex string 	`json:"sex,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Address string `json:"address,omitempty"`
	Wechat string `json:"wechat,omitempty"`
}

//安卓端上传测肤结果请求
type WSSkinRequest struct {
	PicURL 		string		`json:"face, omitempty"`
	Result		string		`json:"result, omitempty"`
}