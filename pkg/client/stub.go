package client

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/robber-core/pkg/grpclb"
	"github.com/eviltomorrow/robber-notification/internal/server"
	"github.com/eviltomorrow/robber-notification/pkg/pb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var (
	EtcdEndpoints = []string{
		"127.0.0.1:2379",
	}
)

func init() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   EtcdEndpoints,
		DialTimeout: 5 * time.Second,
		LogConfig: &zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.ErrorLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	})
	if err != nil {
		panic(fmt.Errorf("create etcd client failure, nest error: %v", err))
	}

	builder := &grpclb.Builder{
		Client: cli,
	}
	resolver.Register(builder)

}

func NewClientForNotification() (pb.NotificationClient, func(), error) {
	target := fmt.Sprintf("etcd:///%s", server.Key)
	conn, err := grpc.DialContext(
		context.Background(),
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewNotificationClient(conn), func() { conn.Close() }, nil
}
