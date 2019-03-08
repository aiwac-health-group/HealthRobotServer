package services
import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
	"HealthRobotServer/constants"
)

type LectureService interface{
	InsertLecture(info *models.LectureInfo) error
	GetLectureTextContent(LectureID string) *models.TextContent
	GetLectureTextAbstractList() []models.TextAbstract
	GetLectureFileAbstractList(Filetype int)  []models.FileAbstract
	GetLectureFileContent(LectureID string)  *models.FileContent
 }

func NewLectureService()  LectureService {
	return &lectureService{
		dao:dao.NewDao(datasource.Instance()),
	}
}


type lectureService struct {
	dao *dao.Dao
}

func (service *lectureService) InsertLecture(info *models.LectureInfo) error {
	if err := service.dao.Insert(constants.Table_Lecture, info); err != nil {
		log.Println("lecture_service.go Insert lecture info err")
		return err}
	return nil
}

func (service *lectureService) GetLectureTextContent(LectureID string) *models.TextContent {
	var info  models.TextContent
	service.dao.Engine.Table(constants.Table_Lecture).Select("lecture.content").Where("lecture.ID = ? ", LectureID).Find(&info)
    return &info
}

func (service *lectureService)  GetLectureTextAbstractList() []models.TextAbstract {
	var info []models.TextAbstract
	service.dao.Engine.Table(constants.Table_Lecture).Select("lecture.ID,lecture.title,lecture.abstract,lecture.updated_at").Find(&info)
    return info
}

func (service *lectureService) GetLectureFileAbstractList(Filetype int)  []models.FileAbstract {
	var info []models.FileAbstract
	service.dao.Engine.Table(constants.Table_Lecture).Select("lecture.ID,lecture.title,lecture.updated_at,lecture.abstract,lecture.cover,lecture.duration").Where("lecture.filetype = ? ", Filetype).Find(&info)
    return info
}

func (service *lectureService) GetLectureFileContent(LectureID string)  *models.FileContent {
	var info  models.FileContent
	service.dao.Engine.Table(constants.Table_Lecture).Select("lecture.filename,lecture.abstract").Where("lecture.ID = ? ", LectureID).Find(&info)
	return &info
}