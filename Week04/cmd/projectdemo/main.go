package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	v1 "week04/api/projectdemo/v1"
	owngrpc "week04/internal/pkg/grpc"
)

func main() {
	//用于监听signal信号
	signalChan := make(chan os.Signal, 1)

	//初始化资源,获取handler
	handler := InitProjectDemo()

	//注册服务
	rpcServer := owngrpc.NewProjectDemoServer(":8888")
	v1.RegisterProjectDemoServer(rpcServer.Server, handler.ProjectDemoService)

	// 创建带有cancel的父context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建errgroup
	group, _ := errgroup.WithContext(ctx)

	//使用errgroup启动rpc服务
	group.Go(func() error {
		return rpcServer.Run(ctx)
	})

	// 监听signal信号
	go func() {
		signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	}()

	//监听signal信号,收到signal信号通知其他http服务退出
	go func() {
		for {
			select {
			case <-signalChan:
				log.Println("received os signal, ready cancel other running server")
				//调用rpc接口,模拟其他事情的处理
				healthCheck(":8888")
				cancel()
			case <-ctx.Done():
				//优雅关闭
				rpcServer.Server.GracefulStop()
				log.Printf("grpc server gracefull stop")
			}
		}
	}()

	if err := group.Wait(); err != nil {
		// 收到第一个错误后，开始关闭全部server流程
		cancel()
		log.Println(err)
	}
}

//进行rpc接口调用,这里为了模拟数据的处理
func healthCheck(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("connect is failed,erros is : %v", err)
	}
	defer conn.Close()

	// grpc client
	projecrDemoClient := v1.NewProjectDemoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// health check
	response, err := projecrDemoClient.HealthCheck(ctx, &v1.ProjectDemoRequest{Name: "埋点检查"})
	if err != nil {
		log.Fatalf("healthcheck failed,error is: %v", err)
	}
	log.Printf("healthcheck success: %s", response.GetMessage())

}
