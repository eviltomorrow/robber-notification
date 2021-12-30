package command

import (
	"context"
	"log"
	"time"

	"github.com/eviltomorrow/robber-notification/pkg/client"
	"github.com/eviltomorrow/robber-notification/pkg/pb"
	"github.com/spf13/cobra"
)

var testSendCmd = &cobra.Command{
	Use:   "test",
	Short: "test send a email to someone",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		stub, close, err := client.NewClientForNotification()
		if err != nil {
			log.Fatalf("Create client for notification failure, nest error: %v\r\n", err)
		}
		defer close()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if _, err := stub.SendEmail(ctx, &pb.Mail{
			To: []*pb.Contact{
				{Address: toSend},
			},
			Subject: "Test email",
			Body:    "<h4>This is a test email</h4>",
		}); err != nil {
			log.Fatalf("[Failure] Send email[%v] is bad, nest error: %v\r\n", toSend, err)
		} else {
			log.Printf("[Success] Send email[%v] is ok!\r\n", toSend)
		}
	},
}

var (
	toSend = "eviltomorrow@163.com"
)

func init() {
	testSendCmd.Flags().StringVar(&toSend, "to", "eviltomorrow@163.com", "eviltomorrow to somebody")
	rootCmd.AddCommand(testSendCmd)
}
