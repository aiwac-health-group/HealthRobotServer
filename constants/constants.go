package constants

//账户类型变量
const  (
	ClientType_doctor = "doctor"
	ClientType_service = "service"
	ClientType_robot = "robot"
	ClientType_admin = "admin"
)

//状态变量
const (
	//在线状态
	Status_outline = "1"
	Status_online = "2"
	Status_onbusy = "3"

	//问诊处理状态
	Status_treat_new = "1"
	Status_treat_onHandle = "2"
	Status_treat_complete = "3"

	//语音挂号状态
	Status_noresponse = "1"
	Status_processing = "2"
	Status_completed  = "3"
)

//表名
const (
	Table_clientInfo = "client_info"
	Table_webClient = "web_profile"
	Table_robot = "robot_profile"
	Table_TreatInfo = "treat_info"
	Table_SkinInfo = "skin_info"
	Table_Lecture = "lecture"
)

//查询sql
const  (
	//用户相关查询
	Query_account = "client_account = ?"
	Query_online_status = "online_status = ? AND client_type = ?"
	//问诊查询
	Query_treat = "client_account = ? AND treat_status = ?"
	Query_treat_status = "treat_status = ?"
)

const (
	// 体检套餐文件地址
	ExaminePackageLink = "http://localhost/uploads/examineList.docx"
)

//讲座类型
const (
	Lecture_text = 1
	Lecture_audio = 2
	Lecture_video = 3
)