package command

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eviltomorrow/robber-core/pkg/system"
	"github.com/eviltomorrow/robber-core/pkg/zlog"
	"github.com/eviltomorrow/robber-core/pkg/znet"
	"github.com/eviltomorrow/robber-notification/internal/config"
	"github.com/eviltomorrow/robber-notification/internal/server"
	"github.com/eviltomorrow/robber-notification/internal/service"
	"github.com/eviltomorrow/robber-notification/pkg/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

var rootCmd = &cobra.Command{
	Use:   "robber-notification",
	Short: "",
	Long:  "  \r\nrobber-notification server running",
	Run: func(cmd *cobra.Command, args []string) {
		if pprofMode {
			go func() {
				port, err := znet.GetFreePort()
				if err != nil {
					log.Fatalf("[Fatal] Get free port failure, nest error: %v\r\n", err)
				}

				log.Printf("[Debug] Debug mode is started, pprof service is listend on http://%s:%d/debug/pprof", system.IP, port)
				if err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
					log.Fatalf("[Fatal] ListenAndServe pprof server failure, nest error: %v\r\n", err)
				}
			}()
		}

		setupCfg()
		setupVars()
		smtp, err := service.LoadSMTPFromFile(cfg.SMTP.Path)
		if err != nil {
			zlog.Fatal("LoadSMTPFromFile failure", zap.Error(err))
		}
		server.SMTP = smtp

		if err := server.StartupGRPC(); err != nil {
			zlog.Fatal("Startup GRPC service failure", zap.Error(err))
		}
		registerCleanFuncs()
		blockingUntilTermination()
	},
}

var (
	cleanFuncs []func() error
	cfg        = config.GlobalConfig
	cfgPath    = ""
	pprofMode  bool
)

func init() {
	rootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: true,
	}
	rootCmd.Flags().StringVarP(&cfgPath, "config", "c", "config.toml", "robber-notification's config file")
	rootCmd.Flags().BoolVarP(&pprofMode, "pprof", "p", false, "robber-notification's pprof mode")
	rootCmd.MarkFlagRequired("config")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func blockingUntilTermination() {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-ch {
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	case syscall.SIGUSR1:
	case syscall.SIGUSR2:
	default:
	}

	for _, f := range cleanFuncs {
		f()
	}
}

func registerCleanFuncs() {
	cleanFuncs = append(cleanFuncs, server.RevokeEtcdConn)
	cleanFuncs = append(cleanFuncs, server.ShutdownGRPC)
}

func setupCfg() {
	if err := cfg.Load(cfgPath, nil); err != nil {
		log.Fatalf("[Fatal] Load config file failure, nest error: %v\r\n", err)
	}

	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            cfg.Log.Level,
		Format:           cfg.Log.Format,
		DisableTimestamp: cfg.Log.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename:   cfg.Log.FileName,
			MaxSize:    cfg.Log.MaxSize,
			MaxDays:    30,
			MaxBackups: 30,
		},
		DisableStacktrace: true,
	})
	if err != nil {
		log.Fatalf("[Fatal] Setup log config failure, nest error: %v\r\n", err)
	}
	zlog.ReplaceGlobals(global, prop)
	zlog.Info("Global Config info", zap.String("cfg", cfg.String()))
}

func setupVars() {
	server.Host = cfg.Server.Host
	server.Port = cfg.Server.Port
	server.Endpoints = cfg.Etcd.Endpoints

	client.EtcdEndpoints = cfg.Etcd.Endpoints
}
