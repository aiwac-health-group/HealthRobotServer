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
	"encoding/json"
	"bytes"
	"log"
	"github.com/dgrijalva/jwt-go"
)
type LectureController struct {
	Ctx iris.Context
	Service services.LectureService
}

type Responseinfo struct {
	status     int
	message   string
}

func (c *LectureController) BeforeActivation(b mvc.BeforeActivation)  {
	//指定LectureController TestLecture功能的访问Url
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/uploadLecturetext","UploadLecturetext")
	b.Handle("POST","/uploadLectureaudio","UploadLectureaudio")
	b.Handle("POST","/uploadLecturevideo","UploadLecturevideo")
}



func (c *LectureController) UploadLecturetext() {

	 lectureinfo := c.GetJsonText(1)
	 status :=c.Service.Insert(lectureinfo)
	 if status !=nil{
		responseinfo := Responseinfo{
			status:     2001,
			message:   "insert error",
		}
	 jsonfile,err:=json.Marshal(responseinfo)
   if err != nil {
		log.Fatal("error:", err)
     }
	  c.Ctx.JSON(jsonfile)
   }

}

func (c *LectureController) UploadLectureaudio() {
	lectureinfo := c.GetJsonFile(2)
	cover :=c.GetAudioFile(lectureinfo)
	lectureinfo.Cover = cover
	status :=c.Service.Insert(lectureinfo)
	if status !=nil{
		responseinfo := Responseinfo{
			status:     2001,
			message:   "insert error",
		}
	jsonfile,err:=json.Marshal(responseinfo)
	if err != nil {
		log.Fatal("error:", err)
     }
	c.Ctx.JSON(jsonfile)
  }
}

func (c *LectureController) UploadLecturevideo() {
	lectureinfo := c.GetJsonFile(3)
	cover := c.GetVideoFile(lectureinfo)
	lectureinfo.Cover = cover
	status :=c.Service.Insert(lectureinfo)
	if status !=nil{
		responseinfo := Responseinfo{
			status:     2001,
			message:   "insert error",
		}
	jsonfile,err:=json.Marshal(responseinfo)
	if err != nil {
		log.Fatal("error:", err)
     }
	c.Ctx.JSON(jsonfile)
  }
}
//获取文件版Lecture model 的json字符信息并格式化
func (c *LectureController)GetJsonFile(filetype int) *models.LectureInfo{
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)
	info := &models.LectureInfo{}
	jsonfile :=&models.JsonFileInfo{}
    if err := c.Ctx.ReadJSON(jsonfile); err != nil{
		    return info
		    // panic(err.Error())
        }else{
			info.Title = jsonfile.Title
			info.Abstract = jsonfile.Blief
			info.Filetype = filetype
			info.Filename = jsonfile.Filename 
			info.HandleService = service
			return info
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
		    // panic(err.Error())
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
func (c *LectureController)GetVideoFile(lecture *models.LectureInfo)string{
	const maxSize = 50 << 20 // 50MB
	 file, info, err := c.Ctx.FormFile("file")
	 if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return "error"
	}
	defer file.Close()
	  //上传文件至/uplods文件夹
	  info.Filename = lecture.Filename
	  fname := info.Filename
	  out, err := os.OpenFile("../uploads/"+fname,
		  os.O_WRONLY|os.O_CREATE, 0666)
	  if err != nil {
		  c.Ctx.StatusCode(iris.StatusInternalServerError)
		  c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		  return "error"
	  }
	  defer out.Close()
	  io.Copy(out, file)
	  filename :="../uploads/"+fname	  
	  return GetBase64Frame(filename)
}

//获取上传音频文件
func (c *LectureController)GetAudioFile(lecture *models.LectureInfo)string{
	const maxSize = 50 << 20 // 50MB
	 file, info, err := c.Ctx.FormFile("file")
	 if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return "error"
	}
	defer file.Close()
	  //上传文件至/uplods文件夹
	  info.Filename = lecture.Filename
	  fname := info.Filename
	  out, err := os.OpenFile("../uploads/"+fname,
		  os.O_WRONLY|os.O_CREATE, 0666)
	  if err != nil {
		  c.Ctx.StatusCode(iris.StatusInternalServerError)
		  c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		  return "error"
	  }
	  defer out.Close()
	  io.Copy(out, file)
		//音频文件默认展示图片logo.jpg
		//该操作为转换默认展示图片的base64码，类型为string
	  filename :="../uploads/logo.jpg"
	  picture, err := os.Open(filename)
	  if err != nil {
		log.Println(err)
		return "error"
    	}
	  defer picture.Close()
	
	  fileinfo, err := picture.Stat()
 	  if err != nil {
		log.Println(err)
		return "error"
	    } 
	
	  filesize := fileinfo.Size()
	  buffer := make([]byte, filesize)
	
	  picture.Read(buffer)
	  if err != nil {
		log.Println(err)
		return "error"
	  }
	  encodeString := base64.StdEncoding.EncodeToString(buffer)
	  return encodeString
}

//该函数为Base64转码，需要注意的是因为go没有处理视频的工具包，所以只能使用cmd命令中的ffmpeg来进行处理
//所以需要在测试时注意文件执行位置，传入filename为.../uploads/filename
func GetBase64Frame(filename string) string {
    width := 2752
	height := 2208
	//这里需要注意下文件位置问题，未测试需要测试
    cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")

    buf := new(bytes.Buffer)
	
    cmd.Stdout = buf
    
    if cmd.Run() != nil {
        panic("could not generate frame")
	}
	input := buf.Bytes()
	encodeString := base64.StdEncoding.EncodeToString(input)
    return encodeString
}

