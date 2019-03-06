package constants

const (
	//医生在线状态
	Status_outline = "1"
	Status_online = "2"
	Status_onbusy = "3"

	//问诊处理状态
	Status_treat_new = "1"
	Status_treat_onHandle = "2"
	Status_treat_complete = "3"
)

const (
	Table_clientInfo = "client_info"
	Table_webClient = "web_profile"
	Table_robot = "robot_profile"
	Table_TreatInfo = "treat_info"
	Table_Lecture = "lecture"
)

const  (
	//用户相关查询
	Query_account = "client_account = ?"
	Query_online_status = "online_status = ? AND client_type = ?"
	//问诊查询
	Query_treat = "client_account = ? AND treat_status = ?"
	Query_treat_status = "treat_status = ?"
)

const  (
	ClientType_doctor = "doctor"
	ClientType_service = "service"
	ClientType_robot = "robot"
	ClientType_admin = "admin"
)