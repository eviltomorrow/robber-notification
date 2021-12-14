package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

func LoadSMTPFromFile(path string) (SMTP, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return SMTP{}, err
	}

	var data = bytes.TrimSpace(buf)
	var s = SMTP{}
	if err := json.Unmarshal(data, &s); err != nil {
		return SMTP{}, err
	}
	return s, nil
}

type SMTP struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

func (m *SMTP) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}
