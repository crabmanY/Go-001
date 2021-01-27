package main

import (
	"fmt"
	rand2 "math/rand"
	"net"
	"testing"
)

func BenchmarkServer(b *testing.B) {
	for i := 0; i < 10; i++ {
		go sendMsg()
	}
}

func sendMsg() bool {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err :", err)
		return true
	}
	defer conn.Close() // 关闭连接
	for {

		_, err = conn.Write([]byte(string(rand2.Int31()))) // 发送数据
		if err != nil {
			return true
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return true
		}
		fmt.Println(string(buf[:n]))
	}
	return false
}
