package protobuf

import (
	"fmt"
	"google.golang.org/grpc"
)

func NewGrpcClient(ip, port string) SendAnalysisResultClient {
	host := ip + ":" + port
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	c := NewSendAnalysisResultClient(conn)
	return c
}
