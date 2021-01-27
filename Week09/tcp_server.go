package main

import (
	"bufio"
	context "context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 20000,
		Zone: "",
	})
	if err != nil {
		panic(err)
	}
	//通知连接关闭
	context, cancelFunc := context.WithCancel(context.Background())

	// 监听signal信号
	signalChan := make(chan os.Signal, 1)
	go func() {
		signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	}()

	//监听signal信号,收到signal信号通知其他goroutine关闭连接
	go func() {
		for {
			select {
			case <-signalChan:
				log.Println("received os signal, ready cancel other conn")
				cancelFunc()
			}
		}
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Errorf("accept failed,err is %+v", err)
			continue
		}
		go process(conn, context)
	}

}

func process(conn net.Conn, parentContext context.Context) {
	cancel, _ := context.WithCancel(parentContext)

	var connChan = make(chan []byte, 1)
	go read(conn, connChan, cancel)
	go write(conn, connChan, cancel)
}

//读取数据
func read(conn net.Conn, dataChan chan []byte, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("close read")
			close(dataChan)
			conn.Close()
			return
		default:
			reader := bufio.NewReader(conn)
			var buf [1024]byte
			n, err := reader.Read(buf[:])
			if err != nil {
				fmt.Errorf("read failed")
				break
			}
			dataChan <- buf[:n]

		}
	}
}

//操作数据
func write(conn net.Conn, dataChan chan []byte, ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("close write")
			return
		default:

			dataByte := <-dataChan
			//do something
			time.Sleep(time.Second * 5)

			_, err := conn.Write(dataByte[:])
			fmt.Println("send message：" + string(dataByte[:]))
			if err != nil {
				break
				fmt.Errorf("write message failed,error is %+v", err)
			}

		}
	}

}
