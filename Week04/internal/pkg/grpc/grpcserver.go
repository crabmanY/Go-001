package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ProjectServer struct {
	*grpc.Server
	address string
}

//创建grpc服务
func NewProjectDemoServer(address string) *ProjectServer {
	server := grpc.NewServer()
	return &ProjectServer{
		Server:  server,
		address: address,
	}
}

//启动grpc服务
func (server *ProjectServer) Run(context context.Context) error {
	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		return err
	}

	log.Printf("grpc server start ,address is %s", server.address)

	return server.Serve(listener)

}
