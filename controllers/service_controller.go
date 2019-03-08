package controllers

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/manager"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"HealthRobotServer/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//service controller处理客服发出的http请求
type ServiceController struct {
	Ctx iris.Context
	Service services.ServiceService
	WsManager manager.WSManager
}

func (c *ServiceController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/changeDoctor","ModifyDoctorProfile")
	b.Handle("POST","/setWebChat","AllocateDoctorForTreat")
	b.Handle("POST","/ownRegist","OwnRegist")
	b.Handle("POST","/getRegist","GetRegist")
	b.Handle("POST","/postRegist","PostRegistResult")
	b.Handle("POST","/postRegist","PostRegistResult")
	b.Handle("POST","/uploadLecturetext","UploadLecturetext")
	b.Handle("POST","/uploadLectureaudio","UploadLectureaudio")
	b.Handle("POST","/uploadLecturevideo","UploadLecturevideo")
}

//客服修改医生信息
func (c *ServiceController) ModifyDoctorProfile() {
	var request models.DoctorProfileModifyRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	//判断医生是否已经存在
	if client := c.Service.SearchDoctorClientInfo(request.Account); client == nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Account doesn't exist",
		})
		return
	}

	//添加详细信息
	var profile = models.DoctorInfo{
		ClientAccount:request.Account,
		ClientName:request.Name,
		ClientType:"doctor",
		Department:request.Department,
		Brief:request.Brief,
	}

	if err := c.Service.UpdateDoctorClientInfo(&profile); err != nil {
		_, _ =c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"failed to modify doctor profile",
		})
		return
	}

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"successfully",
	})

}

//客服分配医生给用户
func (c *ServiceController) AllocateDoctorForTreat() {
	//从token和数据库中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)

	var allocation models.TreatAllocation

	if err := c.Ctx.ReadJSON(&allocation); err != nil {
		fmt.Println("fail to encode request")
		return
	}

	//验证分配的医生账号是否存在
	doctor := c.Service.SearchDoctorClientInfo(allocation.Doctor)
	if doctor == nil || !strings.EqualFold(doctor.ClientType, constants.ClientType_doctor) { //账号不存在，或者分配的账号不是医生账号
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:  "2001",
			Message: "Wrong Account",
		})
	}

	//获取医生和客服的websocket连接，保证两边连接有效时，才可以发送roomID
	DoctorConn := c.WsManager.GetWSConnection(allocation.Doctor)
	PatientConn := c.WsManager.GetWSConnection(allocation.Patient)
	if DoctorConn == nil || PatientConn == nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:  "2001",
			Message: "医生或用户已经下线，请更新问诊列表",
		})
		return
	}

	//更新该treat中责任医生的账号,以及问诊单状态为正在处理
	treat := c.Service.SearchNewTreatInfo(allocation.Patient)
	treat.HandleDoctor = allocation.Doctor
	treat.Status = constants.Status_treat_onHandle
	c.Service.UpdateTreatInfoHandleDoctor(treat)

	//分配成功后，根据当前时间的纳秒值返回一个roomID给对应医生和用户
	RoomID := strconv.FormatInt(time.Now().UnixNano(),10)

	//将语音连接请求发送给医生
	dataForDoctor, _ := json.Marshal(models.WebsocketResponse{
		Code:   "2011",
		WebResponse:models.WebResponse{
			TreatResponse:models.TreatResponse{
				RoomID: RoomID ,
			},
		},
	})
	if err := (*DoctorConn).Write(1, dataForDoctor); err != nil {
		log.Println("Fail to Send Call to Doctor")
		return
	}
	//更新医生状态
	doctor.OnlineStatus = constants.Status_onbusy
	if err := c.Service.UpdateDoctorClientInfo(doctor); err != nil {
		log.Println("Fail to Update Doctor status")
		return
	}
	//把新的在线医生列表推送给客服
	CurrentConn := c.WsManager.GetWSConnection(service)
	c.SendOnlineDoctorList(*CurrentConn)

	//把语音通话房间号返回给机器人
	dataForPatient, _ := json.Marshal(models.WebsocketResponse{
		Status:"2000",
		Message:"",
		RobotResponse:models.RobotResponse{
			Account:allocation.Patient,
			UniqueID:"",
			ClientType:"robot",
			TreatResponse:models.TreatResponse{
				RoomID:RoomID,
			},
		},
	})
	_ = (*PatientConn).Write(1,dataForPatient)

}

//获取在线医生列表
func (c *ServiceController) SendOnlineDoctorList(conn websocket.Connection) {
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
	_ = conn.To("service").EmitMessage(data)
}

// 客服抢挂号单,查询该客服是否能够抢挂号单
func (c *ServiceController) OwnRegist() {
	var req models.RegistProcessRequest
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		log.Println("fail to encode request, ", err)
		return
	}

	idle, mes := c.Service.QueryRegistrationIdle(req.ID, req.HandleSevice)

	if !idle {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:mes,
		})

	} else {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:  "2000",
			Message: mes,
		})
	}
}

// web端：客服待完成挂号单
func (c *ServiceController) GetRegist() {
	log.Println("GetRegist() ")

	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	ServiceAccount := claims["Account"].(string)

	registration := c.Service.QueryRegistrationProcessing(ServiceAccount)

	if registration.ID == 0 {
		log.Println("无待完成挂号单")
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"无待完成挂号单",
		})
	} else {
		_, _ = c.Ctx.JSON(iris.Map{
			"status":      "2000",
			"id":          registration.ID,
			"userName":    "",
			"userAccount": registration.UserAccount,
			"class":       registration.Regist.Province + registration.Regist.City + registration.Regist.Hospital + registration.Regist.Department,
			"others":      "",
		})
	}
}

// 客服挂号反馈
func (c *ServiceController) PostRegistResult() {
	var result models.RegistProcessResponse
	err := c.Ctx.ReadJSON(&result)
	if err != nil && !iris.IsErrPath(err) {
		log.Println("fail to encode request")
		return
	}

	if err := c.Service.SaveRegistResult(&result); err != nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"账号不存在",
		})
		return
	}
	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"成功",
	})

	//把挂号结果反馈给用户
	robotAccount := c.Service.SearchRegistByID(result.ID).UserAccount
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0018",
		RobotResponse:models.RobotResponse{
			MessageNotice: models.MessageNotice {
				MessageType: "2",
				MessageID: result.ID,
				UserAccount:robotAccount,
			},
		},
	})
	RobotConn := c.WsManager.GetWSConnection(robotAccount)
	if RobotConn == nil {
		println("用户下线，挂号结果反馈失败")
		return
	}
	_ = (*RobotConn).Write(1,data)
}

// 根据提交的表单创建体检推荐信息并添加到数据库
func (c *ServiceController) SaveRecommendation() {
	log.Println("Service SaveRecommendation()")

	//解析用户上传的表单文件
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		log.Println("Uploading file Error: ", err)
		return
	}
	form := c.Ctx.Request().MultipartForm

	//保存用户上传的推荐到本地目录
	fileUrl, err := util.SaveFileUploaded("/uploads/", info.Filename, &file)
	log.Println(fileUrl)
	//Base64编码
	cover :=util.GetBase64Frame(fileUrl)

	// 校验函数，校验提交的数据是否满足要求
	if len(form.Value["title"][0]) == 0 || len(form.Value["text"][0]) == 0 {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"题目和具体内容不能为空",
		})
		return
	}

	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	ServiceAccount := claims["Account"].(string)

	recommend := &models.PhysicalExamine{
		Title: form.Value["title"][0],
		Abstract: form.Value["blief"][0],
		Infor: form.Value["text"][0],
		Cover: cover,
		HandleSevice: ServiceAccount,
	}

	log.Printf("recommend is : %v", recommend)

	// 数据库操作
	if err := c.Service.CreatePhysicalExamine(recommend); err != nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"创建体检推荐失败",
		})
		return
	}

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"",
	})

	//将体检推荐推送到所有的用户
	message := models.MessageNotice{
		MessageType: "0",
		MessageID: c.Service.GetTheNewExamine().ID,
		UserAccount:"",
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0018",
		RobotResponse:models.RobotResponse{
			MessageNotice:message,
		},
	})
	Conn := c.WsManager.GetWSConnection(ServiceAccount)
	_ = (*Conn).To("robot").EmitMessage(data)
	return

}

//客服上传文本健康讲座
func (c *ServiceController) UploadLecturetext() {
	var request models.TextLectureUploadRequest
	if err := c.Ctx.ReadJSON(request); err != nil{
		log.Println("encode request fail, ", err)
		return
	}

	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)

	var lecture = models.LectureInfo {
		Title:request.Title,
		Abstract:request.Blief,
		Content:request.Text,
		Filetype:1,
		HandleService:service,
	}

	err := c.Service.InsertLecture(&lecture)
	if err  != nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2001",
			Message: "insert error",
		})
	}else{
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2000",
			Message: "success",
		})
	}
}

//客服上传音频版健康讲座
//解析并存储上传的音频文件
func (c *ServiceController)GetAudioFile(filetype int)*models.LectureInfo{
	const maxSize = 50 << 20 // 50MB

	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)

	var lecture models.LectureInfo

	file, info, err := c.Ctx.FormFile("file")
	form := c.Ctx.Request().MultipartForm
	title := form.Value["title"][0]
	blief := form.Value["blief"][0]

	lecture.Filetype = filetype
	lecture.Title = title
	lecture.Abstract = blief
	lecture.HandleService = service
	lecture.Filename, err = util.SaveFileUploaded("/uploads/lecture/audio/",info.Filename, &file)

	//音频文件默认展示图片logo.jpg
	//该操作为转换默认展示图片的base64码，类型为string
	filename :="http://localhost:8080/uploads/audio/logo.jpg"
	picture, err := os.Open(filename)
	if err != nil {
		log.Println("获取logo图片失败, ", err)
	}
	defer picture.Close()

	fileInfo, err := picture.Stat()
	if err != nil {
		log.Println("解析logo图片失败, ", err)
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	if _, err := picture.Read(buffer); err != nil {
		log.Println(err)
	}
	lecture.Cover = base64.StdEncoding.EncodeToString(buffer)

	return &lecture
}

func (c *ServiceController) UploadLectureaudio() {
	lecture :=c.GetAudioFile(2)
	err :=c.Service.InsertLecture(lecture)
	if err !=nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2001 ",
			Message:   "insert error",
		})
	}else{
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2000 ",
			Message: "success",
		})

	}
}

//客服上传视频版健康讲座
func (c *ServiceController) UploadLecturevideo() {
	lecture := c.GetVideoFile(3)
	status :=c.Service.InsertLecture(lecture)
	if status !=nil{
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2001 ",
			Message:   "insert error",
		})
	}else{
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2000 ",
			Message: "success",
		})

	}
}
//解析并存储上传的视频文件至uploads文件夹中
func (c *ServiceController)GetVideoFile(filetype int) *models.LectureInfo{
	const maxSize = 50 << 20 // 50MB

	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)

	lecture := &models.LectureInfo{}

	file, info, err := c.Ctx.FormFile("file")
	form := c.Ctx.Request().MultipartForm
	// 获取其他参数
	title := form.Value["title"][0]
	blief := form.Value["blief"][0]
	lecture.Filetype = filetype
	lecture.Title = title
	lecture.Abstract = blief
	lecture.HandleService = service
	lecture.Filename, err = util.SaveFileUploaded("/uploads/lecture/video/",info.Filename,&file)

	log.Println("保存视频讲座失败, ", err)

	lecture.Cover = util.GetBase64Frame(lecture.Filename)
	return lecture
}




