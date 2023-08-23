package main

import (
	"flag"
	"fmt"
	"icego/etcd/echo"
	"icego/etcd/echo-client/client"
	"icego/etcd/echo-server/server"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 9000, "the port on which the scanning starts")
	addr = flag.String("addr", "localhost:50051", "the hostname to test")
)

func main() {
	// go 程序采集程序,重新布置项目目录
	// myapp.MyApp()

	go func() {
		grpcServer()
	}()

	time.Sleep(2 * time.Second)
	grpcClient()

}
func grpcServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &server.EchoService{})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func grpcClient() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := echo.NewEchoClient(conn)
	client.CallUnaryEcho(c)
}
