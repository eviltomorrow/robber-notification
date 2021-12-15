package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendWithoutSSL(t *testing.T) {
	_assert := assert.New(t)

	smtp, err := LoadSMTPFromFile("../../scripts/smtp.json")
	_assert.Nil(err)
	fmt.Println(smtp.String())
	err = SendWithSSL(smtp.Server, smtp.Username, smtp.Password, &Message{
		From: Contact{Address: smtp.Username},
		To: []Contact{
			{Address: "eviltomorrow@163.com"},
		},
		Subject:     "This is text",
		Body:        "test",
		ContentType: TextHTML,
	})
	_assert.Nil(err)
}
