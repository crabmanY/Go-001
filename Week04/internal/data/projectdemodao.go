package data

import (
	"log"
	"week04/internal/service"
)

//语法糖
var _ service.ProjectDemoDao = new(projectDemoDao)

type projectDemoDao struct {
}

func (p projectDemoDao) Save(message *service.Message) string {
	//打印存储信息
	log.Printf("projectDemoDao save message is %s", message.Info)
	return message.Info
}

func NewProjectDao() service.ProjectDemoDao {

	return &projectDemoDao{}
}
