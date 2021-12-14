package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendWithoutSSL(t *testing.T) {
	_assert := assert.New(t)

	var (
		username = "x"
		password = "x"
	)
	err := SendWithoutSSL("x", 465, username, password, &Message{
		From: Contact{Address: username},
		To: []Contact{
			{Address: "x@163.com"},
		},
		Subject:     "This is text",
		Body:        "test",
		ContentType: TextHTML,
	})
	_assert.Nil(err)
}
