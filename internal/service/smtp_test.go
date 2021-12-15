package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSMTPFromFile(t *testing.T) {
	_assert := assert.New(t)

	var path = "../../scripts/smtp.json"
	smtp, err := LoadSMTPFromFile(path)
	_assert.Nil(err)

	t.Logf("smtp: %v\r\n", smtp.String())
}
