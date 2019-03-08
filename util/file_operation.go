package util

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	_, _ = io.Copy(out, *file)
	//fileUrl := "http://localhost:8080" + filePath
	fileUrl := "E:/Code/GoProjects/src/HealthRobotServer0307" + filePath
	return fileUrl, err
}
