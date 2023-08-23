package server

import (
	"context"
	"fmt"
	"icego/etcd/echo"
)

type EchoService struct {
	echo.UnimplementedEchoServer
}

func (EchoService) UnaryEcho(ctx context.Context, in *echo.EchoMessage) (*echo.EchoMessage, error) {
	fmt.Println("ServerRecv <<<<<<<<<<<", in.Message)
	return &echo.EchoMessage{
		Message: "ServerSend >>>>>>>>>>" + " hello client",
	}, nil
}
