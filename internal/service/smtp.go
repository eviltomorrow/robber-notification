package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
)

var (
	SMTPConfigPath string
)

func FindSMTPJSONFile() (string, error) {
	possibleSMTPJSON := []string{
		SMTPConfigPath,
		"./smtp.json",
		"/etc/smtp.json",
	}
	for _, path := range possibleSMTPJSON {
		if _, err := os.Stat(path); err == nil {
			absFile, err := filepath.Abs(path)
			if err == nil {
				return absFile, nil
			}
			return path, nil
		}
	}
	return "", fmt.Errorf("fail to find smtp.json")
}

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
