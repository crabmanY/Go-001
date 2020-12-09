package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//用于监听signal信号
	signalChan := make(chan os.Signal, 1)
	// 用于通知所有服务都已关闭
	stop := make(chan struct{})

	// server1
	server1 := http.Server{
		Addr: ":8088",
	}

	// server2
	server2 := http.Server{
		Addr: ":8089",
	}

	// 创建带有cancel的父context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建errgroup
	group, _ := errgroup.WithContext(ctx)

	//使用errgroup启动server1 和server 2
	group.Go(func() error {
		return server1.ListenAndServe()
	})
	group.Go(func() error {
		return server2.ListenAndServe()
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
				cancel()
			case <-ctx.Done():
				log.Println(ctx.Err())
				serverCancelHandler(stop, &server1, &server2)
			}
		}
	}()

	if err := group.Wait(); err != nil {
		// 收到第一个错误后，开始关闭全部server流程
		cancel()
		log.Println(err)
	}
	//阻塞等待所有服务完成之后,进行退出
	<-stop
}

// 处理context cancel,关闭所有的http服务,通知main goroutine退出
func serverCancelHandler(stop chan struct{}, servers ...*http.Server) {
	// 开始优雅关闭
	go func() {
		success := 0
		error := 0
		for _, server := range servers {
			if err := server.Shutdown(context.Background()); err != nil {
				log.Printf("端口%s shutdown failed, err: %v\n", server.Addr, err)
				error++
				continue
			}
			success++
			log.Printf("端口%s shutdown succrss", server.Addr)
		}

		log.Printf("server  graceful shutdown completed, success total is %d,fail total is %d", success, error)
		close(stop)
		return
	}()

	// 超时强制退出
	<-time.After(time.Minute * 5)
	log.Println("graceful shutdown timeout")
	close(stop)
	return
}
