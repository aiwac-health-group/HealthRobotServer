package controllers

import (
	"HealthRobotServer/manager"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"HealthRobotServer/util"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"math"
)

type DoctorController struct {
	Ctx iris.Context
	Service services.DoctorService
	WsManager manager.WSManager
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/userList","GetUserList")
	b.Handle("POST","/postReport","SaveReport")
	b.Handle("GET","/getSkinInfo","GetSkinInfo")
}

//获取待填写健康报告的用户列表
//获取用户总数和医生总数，相除，向上取整得每个医生负责的用户数N，从而获得每个医生负责的用户ID的范围即：医生ID * N ~ （医生ID + 1）* N
//注意ID是在数据库表中的序号，不是Account
//根据ID范围以及用户的HealthStatus状态来返回列表给医生
func (c *DoctorController) GetUserList() {
	//从token和数据库中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	account := claims["Account"].(string)

	//获取该医生在数据库表中的ID
	doctor := c.Service.SearchDoctorClientInfo(account)
	doctorID := doctor.ID

	//获取注册用户和医生的数目
	RobotCount := c.Service.CountTotalRobotClient()
	DoctorCount := c.Service.CountTotalDoctorClient()

	//计算医生负载
	payLoad := int64(math.Ceil(float64(RobotCount)/float64(DoctorCount)))
	//计算该医生的负责用户ID范围
	DownBoundary := doctorID * payLoad
	UpBoundary := (doctorID + 1) * payLoad

	//根据ID范围获取需要处理的用户列表,并发送给医生
	robots := c.Service.GetRobotListForHealthReport(UpBoundary, DownBoundary)
	var items []interface{}
	for _, value := range robots {
		items = append(items, value)
	}
	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Data:models.List{
			Items:items,
		},
	})
}

// web端医生上传健康报告
func (c *DoctorController) SaveReport() {

	//解析用户上传的表单文件
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		log.Println("Uploading file Error: ", err)
		return
	}
	form := c.Ctx.Request().MultipartForm
	robotAccount := form.Value["clientID"][0]

	//保存用户上传的健康报告到本地目录
	fileUrl, err := util.SaveFileUploaded("/uploads/healthReport/", info.Filename, &file)
	if  err != nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"上传文件失败",
		})
	}

	report := models.HealthReport{
		UserAccount:robotAccount,
		Report: fileUrl,
	}

	//存储健康报告
	if err := c.Service.CreateReportInfo(&report); err != nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"上传健康报告失败",
		})
		return
	}

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"上传健康报告成功",
	})

	//医生上传完健康报告，就删除该robot用户所有测肤数据
	_ = c.Service.DeleteSkinInfoByAccount(robotAccount)

	//更新RobotInfo中的HealthStatus
	robot := models.RobotInfo{
		ClientAccount:robotAccount,
		HealthStatus: "2",
	}
	_ = c.Service.UpdateRobotClientInfo(&robot)

	//把健康报告上传消息推送给用户
	reportID := c.Service.SearchReportInfo(robotAccount).ID
	message := models.MessageNotice{
		MessageType: "3",
		MessageID: reportID,
		UserAccount:robotAccount,
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "0018",
		RobotResponse:models.RobotResponse{
			MessageNotice:message,
		},
	})
	RobotConn := c.WsManager.GetWSConnection(robotAccount)
	if RobotConn == nil {
		println("用户下线，挂号结果反馈失败")
		return
	}
	_ = (*RobotConn).Write(1,data)
	return
}

// web端医生获取待上传健康报告用户图片、测肤结果
func (c *DoctorController) GetSkinInfo() {
	var request models.GetSkinRequest
	err := c.Ctx.ReadJSON(&request)
	if err != nil {
		log.Println("fail to encode request")
		return
	}

	// 查询该用户测肤数据
	skinInfos := c.Service.GetSkinInfoByAccount(request.Account)

	if len(skinInfos) == 0{
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"该用户无测肤数据",
		})
		return
	}

	var items []interface{}
	for _, value := range skinInfos {
		items = append(items, iris.Map{
			"jsonString": value.SkinDesc,
			"url": value.FaceURL,
		})
	}
	_, _ = c.Ctx.JSON(iris.Map{
		"status": "2000",
		"data" : items,
	})
}