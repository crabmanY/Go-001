package service

import (
	"context"
	v1 "week04/api/projectdemo/v1"
)

type Message struct {
	Info string
}

//定义数据操作先关接口
type ProjectDemoDao interface {
	Save(message *Message) string
}

//创建service结构体
type ProjectDemoService struct {
	projectDemoDao ProjectDemoDao
	v1.UnimplementedProjectDemoServer
}

//依赖Dao 创建 demo_controller
func NewProjectDemoService(projectDemoDao ProjectDemoDao) *ProjectDemoService {
	return &ProjectDemoService{projectDemoDao: projectDemoDao}
}

//service 保存消息
func (service *ProjectDemoService) Save(message *Message) {
	service.projectDemoDao.Save(message)
}

//实现 rpc接口
func (service *ProjectDemoService) HealthCheck(ctx context.Context, in *v1.ProjectDemoRequest) (*v1.ProjectDemoResponse, error) {

	//调用service保存数据
	service.Save(&Message{
		Info: in.GetName(),
	})

	return &v1.ProjectDemoResponse{Message: in.GetName()}, nil

}
