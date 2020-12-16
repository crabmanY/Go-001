package biz

import (
	"week04/internal/service"
)

//project demo handler struct
type ProjectDemoHandler struct {
	ProjectDemoService *service.ProjectDemoService
}

//依赖service创建handler
func NewProjectDemoHandler(projectDemoService *service.ProjectDemoService) *ProjectDemoHandler {
	return &ProjectDemoHandler{ProjectDemoService: projectDemoService}
}

//handler 保存信息
func (handler *ProjectDemoHandler) SaveHealthCheckInfo(info string) {
	handler.ProjectDemoService.Save(&service.Message{
		Info: info,
	})
}
