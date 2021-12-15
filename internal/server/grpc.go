package server

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/eviltomorrow/robber-core/pkg/grpclb"
	"github.com/eviltomorrow/robber-core/pkg/system"
	"github.com/eviltomorrow/robber-core/pkg/zlog"
	"github.com/eviltomorrow/robber-core/pkg/znet"
	"github.com/eviltomorrow/robber-notification/internal/middleware"
	"github.com/eviltomorrow/robber-notification/internal/service"
	"github.com/eviltomorrow/robber-notification/pkg/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	Host           = "0.0.0.0"
	Port           = 19091
	Endpoints      = []string{}
	RevokeEtcdConn func() error
	Key            = "grpclb/service/notification"

	server *grpc.Server

	SMTP service.SMTP
)

type GRPC struct {
	pb.UnimplementedNotificationServer
}

func (g *GRPC) Version(ctx context.Context, _ *emptypb.Empty) (*wrapperspb.StringValue, error) {
	var buf bytes.Buffer
	buf.WriteString("Server: \r\n")
	buf.WriteString(fmt.Sprintf("   Robber-notification Version (Current): %s\r\n", system.MainVersion))
	buf.WriteString(fmt.Sprintf("   Go Version: %v\r\n", system.GoVersion))
	buf.WriteString(fmt.Sprintf("   Go OS/Arch: %v\r\n", system.GoOSArch))
	buf.WriteString(fmt.Sprintf("   Git Sha: %v\r\n", system.GitSha))
	buf.WriteString(fmt.Sprintf("   Git Tag: %v\r\n", system.GitTag))
	buf.WriteString(fmt.Sprintf("   Git Branch: %v\r\n", system.GitBranch))
	buf.WriteString(fmt.Sprintf("   Build Time: %v\r\n", system.BuildTime))
	buf.WriteString(fmt.Sprintf("   HostName: %v\r\n", system.HostName))
	buf.WriteString(fmt.Sprintf("   IP: %v\r\n", system.IP))
	buf.WriteString(fmt.Sprintf("   Running Time: %v\r\n", system.RunningTime()))
	return &wrapperspb.StringValue{Value: buf.String()}, nil
}

func (g *GRPC) SendEmail(ctx context.Context, req *pb.Mail) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("illegal parameter, nest error: mail is nil")
	}
	if len(req.To) == 0 {
		return nil, fmt.Errorf("illegal parameter, nest error: to is nil")
	}

	var contentType = service.TextHTML
	switch req.ContentType {
	case pb.Mail_TEXT_PLAIN:
		contentType = service.TextPlain
	default:
	}
	var message = &service.Message{
		From: service.Contact{
			Name:    SMTP.Alias,
			Address: SMTP.Username,
		},
		Subject:     req.Subject,
		Body:        req.Body,
		ContentType: contentType,
	}

	var to = make([]service.Contact, 0, len(req.To))
	for _, c := range req.To {
		if c != nil {
			to = append(to, service.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.To = to

	var cc = make([]service.Contact, 0, len(req.Cc))
	for _, c := range req.Cc {
		if c != nil {
			cc = append(cc, service.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.Cc = cc

	var bcc = make([]service.Contact, 0, len(req.Bcc))
	for _, c := range req.Bcc {
		if c != nil {
			bcc = append(bcc, service.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.Bcc = bcc

	if err := service.SendWithSSL(SMTP.Server, SMTP.Username, SMTP.Password, message); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func StartupGRPC() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil {
		return err
	}

	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryServerRecoveryInterceptor,
			middleware.UnaryServerLogInterceptor,
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamServerRecoveryInterceptor,
			middleware.StreamServerLogInterceptor,
		),
	)

	reflection.Register(server)
	pb.RegisterNotificationServer(server, &GRPC{})

	localIp, err := znet.GetLocalIP2()
	if err != nil {
		return fmt.Errorf("get local ip failure, nest error: %v", err)
	}

	close, err := grpclb.Register(Key, localIp, Port, Endpoints, 10)
	if err != nil {
		return fmt.Errorf("register service to etcd failure, nest error: %v", err)
	}
	RevokeEtcdConn = func() error {
		close()
		return nil
	}

	go func() {
		if err := server.Serve(listen); err != nil {
			zlog.Fatal("GRPC Server startup failure", zap.Error(err))
		}
	}()
	return nil
}

func ShutdownGRPC() error {
	if server == nil {
		return nil
	}
	server.Stop()
	return nil
}
