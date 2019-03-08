package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
	"math/rand"
	"time"
)

type SMSService interface {
	SendIdentifyCodeToPhone(phoneNumber string, identifyCode string) error
	GenerateIdentifyCode() int
}

const (
	accessKeyId  = "LTAI2UKC3wrs5Gle"
	accessKeySecret = "bcJzQFwMgCKLWlr11utrDOpbb9BTuf"
)

func NewSMSService() SMSService {
	client, err := sdk.NewClientWithAccessKey("default", accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatal("open SMS service fail")
	}
	return &smsService{Client: client}
}

type smsService struct {
	Client *sdk.Client
}

func (s *smsService) SendIdentifyCodeToPhone(phoneNumber string, identifyCode string) error {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	//手机号
	request.QueryParams["PhoneNumbers"] = phoneNumber
	//控制台里面的身份验证模版
	request.QueryParams["TemplateCode"] = "SMS_159730046"
	request.QueryParams["SignName"] = "aiwac"
	//验证码
	request.QueryParams["TemplateParam"] = "{\"code\":\"" + identifyCode + "\"}"
	_, err := s.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Printf("send IdentifyCode to %s fail: %v", phoneNumber, err)
	}
	log.Printf("send IdentifyCode to %s sucessfully: %v", phoneNumber, err)
	return err
}

func (s *smsService) GenerateIdentifyCode() int {
	rand.Seed(time.Now().Unix())
	identifyCode := rand.Intn(999999)
	return identifyCode
}