package config

import (
	"context"
	"fmt"
	"gim/pkg/grpclib/picker"
	_ "gim/pkg/grpclib/resolver/addrs"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type defaultBuilder struct{}

func (*defaultBuilder) Build() Configuration {
	logger.Level = zap.DebugLevel
	logger.Target = logger.Console

	return Configuration{
		MySQL:                "hedali:123123@tcp(192.168.3.6:3306)/gim?charset=utf8&parseTime=true",
		RedisHost:            "192.168.88.130:6379",
		RedisPassword:        "",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		ConnectLocalAddr:     "192.168.88.130:8000",
		ConnectRPCListenAddr: "192.168.88.130:8000",
		ConnectTCPListenAddr: "192.168.88.130:8001",
		ConnectWSListenAddr:  "192.168.88.130:8002",

		LogicRPCListenAddr:    "192.168.3.6:8010",
		BusinessRPCListenAddr: "192.168.3.6:8020",
		FileHTTPListenAddr:    "192.168.3.6:8030",

		ConnectIntClientBuilder: func() pb.ConnectIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///192.168.88.130:8000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, picker.AddrPickerName)))
			if err != nil {
				panic(err)
			}
			return pb.NewConnectIntClient(conn)
		},
		LogicIntClientBuilder: func() pb.LogicIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///192.168.3.6:8010", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicIntClient(conn)
		},
		BusinessIntClientBuilder: func() pb.BusinessIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///192.168.3.6:8020", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewBusinessIntClient(conn)
		},
	}
}
