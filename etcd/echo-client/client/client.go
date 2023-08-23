package client

import (
	"context"
	"fmt"
	"icego/etcd/echo"
	"log"
)

func CallUnaryEcho(c echo.EchoClient) {
	ctx := context.Background()
	in := &echo.EchoMessage{
		Message: "client say hello",
	}
	res, err := c.UnaryEcho(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("clientRecv <<<<<<<<<<< ", res.Message)
}
