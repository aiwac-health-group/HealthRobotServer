package controllers
import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"HealthRobotServer-master/services"
	"HealthRobotServer-master/models"
	"HealthRobotServer-master/middleware"
  "io"
	"os"
	"os/exec"
	"fmt"
	"encoding/base64"
	"bytes"
	"log"
	"github.com/dgrijalva/jwt-go"
	"time"
	"strconv"
	"mime/multipart"
)
type LectureController struct {
	Ctx iris.Context
	Service services.LectureService
}

//GetBase64Frame函数中必须修改好配置的ffmpeg.exe文件的位置否则函数会因无法找到ffmpeg命令而执行错误
//GetAudioFile 函数中logo.jpg 需要存放在/uploads/audio文件夹下
func (c *LectureController) BeforeActivation(b mvc.BeforeActivation)  {
	//指定LectureController TestLecture功能的访问Url
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/uploadLecturetext","UploadLecturetext")
	b.Handle("POST","/uploadLectureaudio","UploadLectureaudio")
	b.Handle("POST","/uploadLecturevideo","UploadLecturevideo")
}



func (c *LectureController) UploadLecturetext() {

	 lectureinfo := c.GetJsonText(1)
	 err := c.Service.Insert(lectureinfo)
	 if err  != nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2001",
			Message:   "insert error",
		})
   }else{
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status: "2000",
			Message: "success",
		})
	 } 
}

func (c *LectureController) UploadLectureaudio() {
	lectureinfo :=c.GetAudioFile(2)
	err :=c.Service.Insert(lectureinfo)
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

func (c *LectureController) UploadLecturevideo() {
	lectureinfo := c.GetVideoFile(3)
	status :=c.Service.Insert(lectureinfo)
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

//获取文字版Lecture model 的json字符信息并格式化
func (c *LectureController)GetJsonText(filetype int) *models.LectureInfo{
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)

	info := &models.LectureInfo{}
	jsontext :=&models.JsonTextInfo{}
    if err := c.Ctx.ReadJSON(jsontext); err != nil{
    	return info
    }else{
    	info.Title = jsontext.Title
    	info.Abstract = jsontext.Blief
    	info.Content = jsontext.Text
    	info.Filetype = filetype
    	info.HandleService = service
    	return info
	}
}

//获取上传视频文件至uploads文件夹中
//针对表单提交FormFile("file")，file应有一个name为file
func (c *LectureController)GetVideoFile(filetype int) *models.LectureInfo{
	const maxSize = 50 << 20 // 50MB
	file, info, err := c.Ctx.FormFile("file")
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)
	lectureinf := &models.LectureInfo{}
	form := c.Ctx.Request().MultipartForm
	// 获取其他参数
	title := form.Value["title"][0]
	blief := form.Value["blief"][0]
	lectureinf.Filetype = filetype
	lectureinf.Title = title
	lectureinf.Abstract = blief
	lectureinf.HandleService = service
	lectureinf.Filename, err = SaveFileUploaded("/uploads/lecture/video",info.Filename,file)
	log.Println(err)
	lectureinf.Cover=GetBase64Frame(lectureinf.Filename)
	return lectureinf
}

//获取上传音频文件
func (c *LectureController)GetAudioFile(filetype int)*models.LectureInfo{
	const maxSize = 50 << 20 // 50MB
	file, info, err := c.Ctx.FormFile("file")
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)
	lectureinf := &models.LectureInfo{}
	form := c.Ctx.Request().MultipartForm
	// 获取其他参数
	title := form.Value["title"][0]
	blief := form.Value["blief"][0]
	lectureinf.Filetype = filetype
	lectureinf.Title = title
	lectureinf.Abstract = blief
	lectureinf.HandleService = service
	lectureinf.Filename, err = SaveFileUploaded("/uploads/lecture/video",info.Filename,file)
	log.Println(err)
	//音频文件默认展示图片logo.jpg
	//该操作为转换默认展示图片的base64码，类型为string
	filename :="http://localhost:8080/uploads/audio/logo.jpg"
  picture, err := os.Open(filename)
	if err != nil {
	log.Println(err)
  }
  defer picture.Close()
	fileinfo, err := picture.Stat()
 	if err != nil {
	 log.Println(err)
	    } 
	
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, _ = picture.Read(buffer)
	if err != nil {
		log.Println(err)
	  }
	lectureinf.Cover=base64.StdEncoding.EncodeToString(buffer)
	return lectureinf
}

//该函数为Base64转码，需要注意的是因为go没有处理视频的工具包，所以只能使用cmd命令中的ffmpeg来进行处理
//所以需要在测试时注意文件执行位置，传入filename为.../uploads/filename
func GetBase64Frame(filename string) string {
    width := 275
    height := 220

    cmd := exec.Command("E:/Go/libs/src/HealthRobotServer-master/libs/bin/ffmpeg.exe", "-i", filename, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")

    buf := new(bytes.Buffer)
	
    cmd.Stdout = buf
    
    if err := cmd.Run(); err != nil {
        log.Println("could not generate frame, ", err)
	}

	input := buf.Bytes()
	encodeString := base64.StdEncoding.EncodeToString(input)
    return encodeString
}

// 上传文件,path为文件夹路径，filename为文件名, account为用户账号
func SaveFileUploaded(path string, filename string, file *multipart.File) (string, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	//解析文件名，在后缀前面增加一个时间戳
	fileNameParts := strings.Split(filename, ".")
	if len(fileNameParts) != 2 {
		return "", errors.New("文件缺少后缀")
	}
	filePath := path + fileNameParts[0] + "_" + timeStamp + "." + fileNameParts[1]
	out, err := os.OpenFile("." + filePath , os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("upload file err: ", err)
		return "", err
	}
	defer out.Close()
	_, _ = io.Copy(out, *filefile)
	fileUrl := "http://localhost:8080" + filePath
	return fileUrl, err
}