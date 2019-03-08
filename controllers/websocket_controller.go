package controllers

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/manager"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"log"
	"strconv"
	"strings"
)

const (
	BusinessRobotProfile = 6 //机器人注册或修改账号信息
	BussinessExamBlief = 7 //机器人发起体检摘要查询
	BussinessExamInfo = 8 //机器人发起体检详细信息查询
	BusinessLectureAudioAbstract = 9  //健康讲座视频摘要查询
	BusinessLectureVideoAbstract = 10  //健康讲座视频摘要查询
	BusinessLectureFileContent  = 11  //健康讲座音频和视频详情查询
	BusinessLectureTextAbstract = 12  //健康讲座文本摘要查询
	BusinessLectureTextContent  = 13  //健康讲座文本详情查询
	BussinessReportInfo = 15 //机器人用户发起的健康检测结果详细信息查询
	BussinessSkinUpload = 16 //安卓端上传测肤结果
	BusinessTreatRequest = 17 //机器人发起问诊请求
	BussinessRegistRequest = 19 //机器人进行在线挂号
	BussinessRegistRecord = 20 //机器人查询挂号历史记录
	BussinessRegistResult = 21 //机器人查询挂号结果
	BussinessExamPackage = 22 //机器人获取体检套餐文档
	Bussiness3NewExam = 23 //机器人查询体检推荐最新三条数据
	BusinessTreatHangOut = 24 //机器人挂断问诊请求
	BusinessOnlineDoctor = 2009 //客服获取在线医生列表
	BusinessTreatList = 2010 //客服获取待问诊列表
	BusinessDoctorHangOut = 2012 //医生主动挂断问诊电话
	BusinessDoctorRejectCall = 2013 //医生拒绝接听问诊电话
)

type WebsocketController struct {
	Ctx iris.Context
	Conn websocket.Connection
	Service services.WebsocketService
	WsManager manager.WSManager
}

func (c *WebsocketController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandlerWS().Serve, middleware.NewAuthToken().Serve)
	b.Handle("GET","/","Join")
}

//从token中提取的该用户的账号和账户类别信息
var (
	ws_account string
	ws_clientType string
	ws_clientName string
)

func (c *WebsocketController) Join() {
	//从token和数据库中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	ws_account = claims["Account"].(string)
	ws_clientType = claims["ClientType"].(string)

	if strings.EqualFold(ws_clientType, constants.ClientType_robot) {
		info := c.Service.SearchRobotClientInfo(ws_account)
		ws_clientName = info.ClientName
	} else if strings.EqualFold(ws_clientType, constants.ClientType_doctor) {
		info := c.Service.SearchDoctorClientInfo(ws_account)
		ws_clientName = info.ClientName
	} else {
		info := c.Service.SearchServiceClientInfo(ws_account)
		ws_clientName = info.ClientName
	}

	log.Println("New Websocket Connection: ",ws_account, ws_clientType, ws_clientName)

	//加入对应clientType的room, 每个room存放了相应用户类型的所有websocket连接
	c.Conn.Join(ws_clientType)
	//注册连接断开回调函数
	c.Conn.OnDisconnect(c.LoseConnection)
	//注册消息接收处理函数
	c.Conn.OnMessage(c.ReceiveRequest)
	//存储该用户和对应连接的映射关系
	c.WsManager.AddMapRelationship(ws_account,&(c.Conn))

	//更新用户状态为在线
	if strings.EqualFold(ws_clientType, constants.ClientType_robot) { //机器人用户上线
		_ = c.Service.UpdateRobotClientInfo(&models.RobotInfo{
			ClientAccount:ws_account,
			OnlineStatus:constants.Status_online,
		})
	} else if strings.EqualFold(ws_clientType, constants.ClientType_service) || strings.EqualFold(ws_clientType, constants.ClientType_admin) {
		_ = c.Service.UpdateServiceClientInfo(&models.ServiceInfo{
			ClientAccount:ws_account,
			OnlineStatus:constants.Status_online,
		})
	} else { //医生用户上线
		_ = c.Service.UpdateDoctorClientInfo(&models.DoctorInfo{
			ClientAccount:ws_account,
			OnlineStatus:constants.Status_online,
		})
		//获取在线医生列表,并将列表推送给所有的客服
		c.PushOnlineDoctorList()
	}

	//开启事件监听
	c.Conn.Wait()
}

func (c *WebsocketController) LoseConnection() {
	//删除用户与连接的映射关系
	c.WsManager.DeleteMapRelationship(ws_account)

	//更新用户状态为下线
	if strings.EqualFold(ws_clientType, constants.ClientType_robot) { //机器人用户下线
		_ = c.Service.UpdateRobotClientInfo(&models.RobotInfo{
			ClientAccount:ws_account,
			OnlineStatus:constants.Status_outline,
		})
	} else if strings.EqualFold(ws_clientType, constants.ClientType_service) || strings.EqualFold(ws_clientType, constants.ClientType_admin) {
		_ = c.Service.UpdateServiceClientInfo(&models.ServiceInfo{
			ClientAccount:ws_account,
			OnlineStatus:constants.Status_outline,
		})
	} else { //医生用户下线
		_ = c.Service.UpdateDoctorClientInfo(&models.DoctorInfo{
			ClientAccount:ws_account,
			OnlineStatus:constants.Status_outline,
		})
		//获取在线医生列表,并将列表推送给所有的客服
		c.PushOnlineDoctorList()
	}

	log.Printf("%s %s lose the connection", ws_account, ws_clientType)
}

func (c *WebsocketController) ReceiveRequest(data []byte) {
	//在这里解析收到的请求，根据请求中的业务号跳转到指定业务处理函数中进行处理
	var request models.WSRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Println("Websocket request from Explore Unmarshal err: ",err)
	}
	log.Printf("websocket request, %s", request)
	//根据request中的code字段配置相应的函数处理
	businessCode, _ := strconv.Atoi(request.BusinessCode)

	switch businessCode {
	case BusinessRobotProfile: c.RobotProfileHandler(&request)
	case BussinessExamBlief: c.QueryExamineBliefHandler(&request)
	case BussinessExamInfo: c.QueryExamineHandler(&request)
	case BussinessSkinUpload: c.SaveSkinTestHandler(&request)
	case BusinessLectureAudioAbstract: c.LectureAudioAbstractHandler(&request)
	case BusinessLectureVideoAbstract: c.LectureVideoAbstractHandler(&request)
	case BusinessLectureTextAbstract:  c.LectureTextAbstractHandler(&request)
	case BusinessLectureTextContent:   c.LectureTextContentHandler(&request)
	case BusinessLectureFileContent:   c.LectureTextContentHandler(&request)
	case BusinessTreatRequest: c.TreatRequestHandler(&request)
	case BusinessTreatHangOut: c.TreatHangOutHandler(&request)
	case BussinessRegistRequest:c.RegistRequestHandler(&request)
	case BusinessOnlineDoctor: c.DoctorListRequestHandler(&request)
	case BusinessTreatList: c.TreatWaitListHandler(&request)
	case BusinessDoctorHangOut: c.DoctorHangOutHandler(&request)
	case BusinessDoctorRejectCall: c.DoctorRejectCallHandler(&request)
	case BussinessRegistRecord:c.QueryRegistRecordHandler(&request)
	case BussinessRegistResult: c.QueryRegistResultHandler(&request)
	case BussinessExamPackage: c.GetExaminePackage(&request)
	case Bussiness3NewExam: c.Get3NewExamineList(&request)
	}
}

//6号业务处理
//处理机器人用户发起的个人信息注册及修改
func (c *WebsocketController) RobotProfileHandler(request *models.WSRequest) {
	//判断账号是否已经存在
	if client := c.Service.SearchRobotClientInfo(request.Account); client.ID == 0 {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0006",
			Status:"2001",
			Message:"the robot doesn't exist",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
			},
		})
		_ = c.Conn.Write(1,data)
		return
	}

	//添加详细信息
	var profile = models.RobotInfo{
		ClientAccount:request.Account,
		ClientName:request.Name,
		ClientType:"robot",
		Sex:request.Sex,
		Birthday:request.Birthday,
		Address:request.Address,
		Wechat:request.Wechat,
	}

	if err := c.Service.UpdateRobotClientInfo(&profile); err != nil {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0006",
			Status:"2001",
			Message:"system error",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
			},
		})
		_ = c.Conn.Write(1,data)
		return
	}

	data, _ := json.Marshal(models.WebsocketResponse{
		Code:"0006",
		Status:"2000",
		Message:"Register or update profile successfully",
		RobotResponse:models.RobotResponse{
			Account:ws_account,
			UniqueID:"",
		},
	})
	_ = c.Conn.Write(1,data)
	return
}

//0007号业务处理
//机器人用户发起的体检推荐摘要查询
func (c *WebsocketController) QueryExamineBliefHandler(request *models.WSRequest) {
	var list []models.PhysicalExamine
	list = c.Service.GetAllExamine()
	var items []interface{}
	for _, value := range list{
		items = append(items, iris.Map{
			"examID": value.ID,
			"name": value.Title,
			"description": value.Abstract,
			"updateTime": value.UpdatedAt,
			"cover": value.Cover,
		})
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Status: "2000",
		Code: "0007",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("ExamineList: %s", data)
	_ = c.Conn.Write(1,data)
}

//0008号业务处理
//机器人用户发起的体检推荐详情查询
func (c *WebsocketController) QueryExamineHandler(request *models.WSRequest) {
	examine := c.Service.SearchExamineInfo(request.ExamID)

	data, _ := json.Marshal(models.WebsocketResponse{
		Status: "2000",
		Code: "0008",
		RobotResponse:models.RobotResponse{
			ExamineResponse:models.ExamineResponse{
				Examine: examine.Infor,
			},
		},
	})
	log.Printf("Examine: %s", data)
	_ = c.Conn.Write(1, data)
}

//0009号业务处理
//获取音频课程摘要
func(c *WebsocketController) LectureAudioAbstractHandler(request *models.WSRequest){
	lectures := c.Service.GetLectureFileAbstractList(constants.Lecture_audio)
	var items []interface{}
	for _, value := range lectures {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0009",
		Status:"2000",
		Message:"Successful",
		Data: models.List{
			Items: items,
		},
		RobotResponse:models.RobotResponse{
			Account: request.Account,
			ClientType: request.ClientType,
			UniqueID: request.UniqueID,
		},
	})
	log.Printf("lectureList: %s", data)
	_ = c.Conn.Write(1,data)
}

//0010号业务处理
//获取视频课程摘要
func(c *WebsocketController) LectureVideoAbstractHandler(request *models.WSRequest){
	lectures := c.Service.GetLectureFileAbstractList(constants.Lecture_video)
	var items []interface{}
	for _, value := range lectures {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0010",
		Status:"2000",
		Message:"string",
		Data: models.List{
			Items: items,
		},
		RobotResponse:models.RobotResponse{
			Account: request.Account,
			ClientType: request.ClientType,
			UniqueID: request.UniqueID,
		},
	})
	log.Printf("lectureList: %s", data)
	_ = c.Conn.Write(1,data)
}

//0011号业务处理
//获取音视频讲座内容
func(c *WebsocketController) LectureFileContentHandler(request *models.WSRobotRequest){
	lectures := c.Service.GetLectureFileContent(request.LectureID)
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0011",
		Status:"2000",
		Message:"string",
		RobotResponse:models.RobotResponse{
			Account: request.Account,
			ClientType: request.ClientType,
			UniqueID: request.UniqueID,
			LectureResponse:models.LectureResponse{
				Link:lectures.Filename,
			},
		},
	})
	log.Printf("lectureList: %s", data)
	_ = c.Conn.Write(1,data)
}

//0012号业务处理
//获取文本讲座摘要
func(c *WebsocketController) LectureTextAbstractHandler(request *models.WSRequest){
	lectures := c.Service.GetLectureTextAbstractList()
	var items []interface{}
	for _, value := range lectures {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0012",
		Status:"2000",
		Message:"string",
		Data: models.List{
			Items: items,
		},
		RobotResponse:models.RobotResponse{
			Account: request.Account,
			ClientType: request.ClientType,
			UniqueID: request.UniqueID,
		},
	})
	log.Printf("lectureList: %s", data)
	_ = c.Conn.Write(1,data)
}

//0013号业务处理
//获取文本讲座内容
func(c *WebsocketController) LectureTextContentHandler(request *models.WSRequest) {

	lectures := c.Service.GetLectureTextContent(request.LectureID)
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0013",
		Status:"2000",
		Message:"string",
		RobotResponse:models.RobotResponse{
			Account: request.Account,
			ClientType: request.ClientType,
			UniqueID: request.UniqueID,
			LectureResponse:models.LectureResponse{
				LectureContext:lectures.Content,
			},
		},
	})
	log.Printf("lectureList: %s", data)
	_ = c.Conn.Write(1,data)
}

//0015号业务处理
//机器人用户发起的健康检测结果详细信息查询
func (c *WebsocketController) QueryReportHandler(request *models.WSRequest) {
	report := c.Service.GetReportByID(request.ReportID)
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0015",
		Status:"2000",
		Message:"Successful",
		RobotResponse:models.RobotResponse{
			Account: request.Account,
			ClientType: request.ClientType,
			UniqueID: request.UniqueID,
			HealthReportResponse:models.HealthReportResponse{
				Report:report.Report,
			},
		},
	})
	log.Printf("Report: %s", data)
	_ = c.Conn.Write(1, data)

}

//0016号业务处理
//安卓端测肤结果上传
func (c *WebsocketController) SaveSkinTestHandler(request *models.WSRequest) {
	//判断账号是否存在
	if client := c.Service.SearchRobotClientInfo(request.Account); client.ID == 0 {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0006",
			Status:"2001",
			Message:"the robot doesn't exist",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
			},
		})
		_ = c.Conn.Write(1,data)
		return
	}

	var skinInfo = models.SkinTest{
		UserAccount: request.Account,
		FaceURL: request.PicURL,
		SkinDesc: request.Result,
	}

	if err := c.Service.CreateSkinInfo(&skinInfo); err != nil {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0016",
			Status:"2001",
			Message:"system error",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
			},
		})
		_ = c.Conn.Write(1,data)
		return
	}

	data, _ := json.Marshal(models.WebsocketResponse{
		Code:"0016",
		Status:"2000",
		Message:"Upload skin result successfully",
		RobotResponse:models.RobotResponse{
			Account:ws_account,
			UniqueID:"",
		},
	})
	_ = c.Conn.Write(1,data)
	return
}

//17号业务处理
//处理用户发起的问诊请求
//把问诊请求存进数据库，并把等待列表推送给在线客服
func (c *WebsocketController) TreatRequestHandler(request *models.WSRequest) {
	//把问诊请求存进数据库
	c.Service.CreatTreatInfoRequest(&models.TreatInfo{
		Account:ws_account,
		ClientName:ws_clientName,
		Others:"",
	})

	//把正在等待的问诊请求列表推送给客服
	c.PushTreatWaitList()
}

//0019号业务处理
//机器人用户发起在线挂号
//把挂号信息存入数据库，并把等待挂号列表更新给在线客服
func (c *WebsocketController) RegistRequestHandler(request *models.WSRequest) {
	err := c.Service.CreateRegistInfo(&models.Registration{
		UserAccount: ws_account,
		Regist: request.Regist,
		Status: constants.Status_noresponse,
	})
	if err != nil {
		log.Println("接收挂号请求失败: ", err)
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0019",
			Status:"2001",
			Message:"挂号失败，系统错误",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
				ClientType:constants.ClientType_robot,
			},
		})
		_ = c.Conn.Write(1, data)
	} else {
		log.Println("接收挂号请求成功: ", err)
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0019",
			Status:"2000",
			Message:"挂号成功",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
				ClientType:constants.ClientType_robot,
			},
		})
		_ = c.Conn.Write(1, data)
	}
	c.PushOnlineRegistList()
}

//0020号业务处理
//机器人用户查询挂号历史记录
func (c *WebsocketController) QueryRegistRecordHandler(request *models.WSRequest) {
	list := c.Service.SearchRegistByUser(ws_account)

	if len(list) == 0 {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code: "0020",
			Status: "2001",
			Message: "该用户无挂号记录",
		})
		_ = c.Conn.Write(1,data)
		return
	}

	var items []interface{}
	for _, value := range list{
		items = append(items, iris.Map{
			"registerID": value.ID,
			"province": value.Province,
			"city": value.City,
			"hospital": value.Hospital,
			"department": value.Department,
			"registerStatus": value.Status,
			"description": value.RegistDesc,
			"createTime": value.CreatedAt,
			"updateTime": value.UpdatedAt,
		})
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0020",
		Status: "2000",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("registRecord: %s", data)
	_ = c.Conn.Write(1,data)
}

//0021号业务处理
//机器人用户使用挂号ID查询挂号结果
func (c *WebsocketController) QueryRegistResultHandler(request *models.WSRequest) {
	registInfo := c.Service.SearchRegistByID(request.RegisterID)
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0021",
		Status:"2000",
		Message:"返回挂号结果",
		RobotResponse: models.RobotResponse {
			Account:ws_account,
			ClientType:constants.ClientType_robot,
			UniqueID:"",
			RegistResponse: models.RegistResponse {
				Description: registInfo.RegistDesc,
				RegisterStatus: registInfo.Status,
				RegistInfo: models.Regist{
					Province: registInfo.Province,
					City: registInfo.City,
					Hospital: registInfo.Hospital,
					Department: registInfo.Department,
				},
				UpdateTime: registInfo.UpdatedAt,
				CreateTime: registInfo.CreatedAt,
			},
		},
	})
	log.Printf("Regist result: %s", data)
	_ = c.Conn.Write(1, data)
}

//0022号业务处理
//安卓端从服务器拉去体检套餐文档
func (c *WebsocketController) GetExaminePackage(request *models.WSRequest) {
	data, _ := json.Marshal(models.WebsocketResponse{
		Status: "2000",
		Code: "0008",
		RobotResponse:models.RobotResponse{
			ExamineResponse:models.ExamineResponse{
				ExaminePackage:constants.ExaminePackageLink,
			},
		},
	})
	log.Printf("ExaminePackage: %s", data)
	_ = c.Conn.Write(1, data)
}

//0023号业务处理
//机器人用户查询体检推荐最新三条数据
func (c *WebsocketController) Get3NewExamineList(request *models.WSRequest) {
	var list []models.PhysicalExamine
	list = c.Service.Get3NewExamine()
	var items []interface{}
	for _, value := range list{
		items = append(items, iris.Map{
			"examID": value.ID,
			"cover": value.Cover,
		})
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Status: "2000",
		Code: "0023",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("3NewExamineList: %s", data)
	_ = c.Conn.Write(1,data)
}

//24号业务处理
//机器人端挂断问诊请求
func (c *WebsocketController) TreatHangOutHandler(request *models.WSRequest) {
	//获取该用户问诊请求信息,更新其状态
	treat := c.Service.SearchNotCompleteTreatInfo("patient", ws_account)
	treat.Status = constants.Status_treat_complete
	c.Service.UpdateTreatInfoStatus(treat)
	//推送新的问诊列表给客服
	c.PushTreatWaitList()
	//同时更新对应医生的状态为空闲
	doctor := c.Service.SearchDoctorClientInfo(treat.HandleDoctor)
	if doctor.ID != 0 {
		doctor.OnlineStatus = constants.Status_online
		_ = c.Service.UpdateDoctorClientInfo(doctor)
	}
	//把空闲状态的医生列表推送给客服
	c.PushOnlineDoctorList()
}

//2008号业务处理
//推送挂号列表至所有客服
func (c *WebsocketController) PushOnlineRegistList() {
	list := c.Service.GetNonresponseRegist()
	var items []interface{}
	for _, value := range list {
		info := c.Service.SearchRobotClientInfo(value.UserAccount)
		userName := info.ClientName
		items = append(items, iris.Map{
			"id": value.ID,
			"account": value.UserAccount,
			"name": userName ,
			"class": value.Province + value.City + value.Hospital + value.Department,
			"others": "",
		})
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2008",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("registWaitList: %s", data)
	_ = c.Conn.To("service").EmitMessage(data)
}
//获取待挂号列表
func (c *WebsocketController) SendOnlineRegistList() {
	list := c.Service.GetNonresponseRegist()
	var items []interface{}
	for _, value := range list {
		info := c.Service.SearchRobotClientInfo(value.UserAccount)
		userName := info.ClientName
		items = append(items, iris.Map{
			"id": value.ID,
			"account": value.UserAccount,
			"name": userName ,
			"class": value.Province + value.City + value.Hospital + value.Department,
			"others": "",
		})
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2008",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("registWaitList: %s", data)
	_ = c.Conn.Write(1, data)
}
func (c *WebsocketController) GetRegistWaitListHandler(request *models.WSRequest) {
	if !strings.EqualFold(request.Message, "getServiceOnlineList") {
		return
	}
	c.SendOnlineRegistList()
}

//2009号业务处理
//推送列表至客服
func (c *WebsocketController) PushOnlineDoctorList() {
	doctors := c.Service.GetOnlineDoctor()
	var items []interface{}
	for _, value := range doctors {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2009",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("OnlineDoctorList: %s", data)
	_ = c.Conn.To("service").EmitMessage(data)
}
//获取在线医生列表
func (c *WebsocketController) GetOnlineDoctorList() {
	doctors := c.Service.GetOnlineDoctor()
	var items []interface{}
	for _, value := range doctors {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2009",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("OnlineDoctorList: %s", data)
	_ = c.Conn.Write(1,data)
}
func (c *WebsocketController) DoctorListRequestHandler(request *models.WSRequest) {
	if strings.EqualFold(request.Message, " getDoctorList ") {
		println("getDoctorList")
		c.GetOnlineDoctorList()
	}
}

//2010号业务处理
//推送等候问诊列表至所有客服
func (c *WebsocketController) PushTreatWaitList() {
	treats := c.Service.SearchNewTreatInfoList()
	var items []interface{}
	for _, value := range treats {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2010",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("TreatWaitList: %s", data)
	_ = c.Conn.To("service").EmitMessage(data)
}
//客服主动获取等待问诊列表
func (c *WebsocketController) GetTreatWaitList() {
	treats := c.Service.SearchNewTreatInfoList()
	var items []interface{}
	for _, value := range treats {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2010",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("TreatWaitList: %s", data)
	_ = c.Conn.Write(1,data)
}

func (c *WebsocketController) TreatWaitListHandler(request *models.WSRequest)  {
	if !strings.EqualFold(request.Message, " getWaitList ") {
		return
	}
	c.GetTreatWaitList()
}

//2012号业务处理
//医生主动挂断语音
func (c *WebsocketController) DoctorHangOutHandler(request *models.WSRequest) {
	//更新医生状态
	doctor := c.Service.SearchDoctorClientInfo(ws_account)
	doctor.OnlineStatus = constants.Status_online
	_ = c.Service.UpdateDoctorClientInfo(doctor)
	//推送新的列表到客服
	c.PushOnlineDoctorList()
}

//2013号业务处理
//医生拒绝接听电话
func (c *WebsocketController) DoctorRejectCallHandler(request *models.WSRequest)  {
	//根据医生获取对应的未完成的问诊单,删除掉责任医生,把问诊状态重置为未处理
	treat := c.Service.SearchNotCompleteTreatInfo(constants.ClientType_doctor, ws_account)
	treat.HandleDoctor = "-"
	treat.Status = constants.Status_treat_new
	c.Service.UpdateTreatInfoStatus(treat)
	//更新客服的问诊列表
	c.PushTreatWaitList()
	//更新医生状态为空闲
	doctor := c.Service.SearchDoctorClientInfo(ws_account)
	doctor.OnlineStatus = constants.Status_online
	_ = c.Service.UpdateDoctorClientInfo(doctor)
	//推送新的列表到客服
	c.PushOnlineDoctorList()
}
