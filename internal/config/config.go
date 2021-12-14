package config

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
)

type Config struct {
	SMTP   SMTP   `json:"smtp" toml:"smtp"`
	Log    Log    `json:"log" toml:"log"`
	Etcd   Etcd   `json:"etcd" toml:"etcd"`
	Server Server `json:"server" toml:"server"`
}

type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	FileName         string `json:"filename" toml:"filename"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
}

type Etcd struct {
	Endpoints []string `json:"endpoints" toml:"endpoints"`
}

type Server struct {
	Host string `json:"host" toml:"host"`
	Port int    `json:"port" toml:"port"`
}

type SMTP struct {
	Path string `json:"path" toml:"path"`
}

func (c *Config) Load(path string, override func(cfg *Config)) error {
	if path == "" {
		return nil
	}

	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}
	return nil
}

func (cg *Config) String() string {
	buf, _ := json.Marshal(cg)
	return string(buf)
}

var GlobalConfig = &Config{
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		FileName:         "/tmp/robber-notification/data.log",
		MaxSize:          20,
	},
	Etcd: Etcd{
		Endpoints: []string{
			"localhost:2379",
		},
	},
	Server: Server{
		Host: "0.0.0.0",
		Port: 19091,
	},
	SMTP: SMTP{
		Path: "/etc/smtp.json",
	},
}
