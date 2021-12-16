package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/robber-notification/pkg/client"
	"github.com/eviltomorrow/robber-notification/pkg/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestVersion(t *testing.T) {
	stub, close, err := client.NewClientForNotification()
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repley, err := stub.Version(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Version error: %v", err)
	}
	fmt.Println(repley.Value)
}

func TestSendEmail(t *testing.T) {
	stub, close, err := client.NewClientForNotification()
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repley, err := stub.SendEmail(ctx, &pb.Mail{
		To: []*pb.Contact{
			{Name: "shepard", Address: "eviltomorrow@163.com"},
		},
		Subject:     "This is one test",
		Body:        "<h1>Hello world</h1>",
		ContentType: pb.Mail_TEXT_HTML,
	})
	if err != nil {
		log.Fatalf("Version error: %v", err)
	}
	fmt.Println(repley.String())
}
