//+build wireinject

package main

import (
	"github.com/google/wire"
	"week04/internal/biz"
	"week04/internal/data"
	"week04/internal/service"
)

//利用wire构建依赖
func InitProjectDemo() *biz.ProjectDemoHandler {

	wire.Build(data.NewProjectDao, service.NewProjectDemoService, biz.NewProjectDemoHandler)
	return &biz.ProjectDemoHandler{}
}
