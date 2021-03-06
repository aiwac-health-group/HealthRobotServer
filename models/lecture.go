package models
import "time"
//健康讲座表的字段
//健康讲座分为文字版，语音版，视频版三种类型

type LectureInfo struct {
	 Base  
	 Title       string     `gorm:"column:title;type:varchar(100)"`    //健康讲座题目
	 Abstract    string 	`gorm:"column:abstract;type:varchar(500)"` //健康讲座摘要
	 Content     string     `gorm:"column:content;type:text"`          //健康讲座文字内容
	 Filename    string     `gorm:"column:filename;type:varchar(100)"` //音频文件名
	 Filetype    int        `gorm:"column:filetype;type:int(11)"`      //健康讲座内容类型  1为文本类型，2为音频类型，3为视频类型
	 Duration    string     `gorm:"column:duration;type:varchar(100)"` //健康讲座时长，默认为Null
	 Cover       string     `gorm:"column:cover;type:text"`    //健康讲座缩略图
	 HandleService string   `gorm:"column:handle_service;type:varchar(11)"`//健康讲座处理客服工号
}

type TextAbstract struct{
	 ID   int        `gorm:"column:ID;type:int(11)" json:"lectureID"`
	 Title        string     `gorm:"column:title;type:varchar(100)" json:"name"`
	 Updated_at  time.Time  `gorm:"column:updated_at"    json:"updateTime"`
	 Abstract    string 	`gorm:"column:abstract;type:varchar(500)" json:"description"`
}

type FileAbstract struct{
	 ID   int        `gorm:"column:ID;type:int(11)" json:"lectureID"`
	 Title        string     `gorm:"column:title;type:varchar(100)" json:"name"`
	 Updated_at  time.Time  `gorm:"column:updated_at"    json:"updateTime"`
	 Description  string    `gorm:"column:abstract;type:varchar(500)"   json:"description"`
	 Cover        string    `gorm:"column:cover;type:varchar(100)"   json:"cover"`
	 Duration    string     `gorm:"column:duration;type:varchar(100)"   json:"duration"`
}

type TextContent struct{
	Content  string    `gorm:"column:content;type:varchar(500)"   json:"lectureContext"`
}

type FileContent struct{
	Filename    string     `gorm:"column:filename;type:varchar(100)" json:"link"` 
}