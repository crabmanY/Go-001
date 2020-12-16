学习笔记
####本周作业
按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

####作业思路
参照`kratos`创建目录,利用`wire`去初始化`dao`,`service`,`biz`之间的依赖关系,最后在结合`errgroup`和`context`去处理优雅退出问题


####目录结构
```
├── README.md
├── api ------------------------------------------存放grpc pb相关文件
│   └── projectdemo
│       └── v1
│           ├── projectdemo.pb.go
│           ├── projectdemo.proto
│           └── projectdemo_grpc.pb.go
├── cmd ----------------------------------------- wire初始化依赖以及启动函数
│   └── projectdemo
│       ├── main.go -----------------------  启动函数
│       ├── wire.go
│       └── wire_gen.go -------------------- 初始化依赖
├── go.mod
├── go.sum
└── internal ------------------------------------  内部处理
    ├── biz ------------------------------------- biz层(依赖service)
    │   └── projectdemohandler.go
    ├── data ------------------------------------- dao层
    │   └── projectdemodao.go
    ├── pkg  ------------------------------------- grpc先关定义
    │   └── grpc
    │       └── grpcserver.go
    └── service ---------------------------------- service层(依赖dao)
        └── projectdemoservice.go


```



#### 对于signal信号的处理
```
2020/12/16 22:04:13 grpc server start ,address is :8888
2020/12/16 22:04:15 received os signal, ready cancel other running server
2020/12/16 22:04:15 projectDemoDao save message is 埋点检查
2020/12/16 22:04:15 healthcheck success: 埋点检查
2020/12/16 22:04:15 grpc server gracefull stop

```