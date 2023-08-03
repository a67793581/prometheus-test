package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"prometheus-test/infrastructure/config"
	"prometheus-test/infrastructure/drivers"
	"prometheus-test/infrastructure/environment"
	"prometheus-test/infrastructure/http_client/trace_http"
	"prometheus-test/infrastructure/metrics"
	"prometheus-test/infrastructure/recycle"
	"prometheus-test/lib/logger"
	"prometheus-test/server/httpserver"
	"runtime/debug"
	"syscall"
)

type Level int

const (
	Info  = Level(0)
	Fatal = Level(2)
)

var (
	Ctx      context.Context
	confPath = flag.String("c", "./conf/common.dev.toml", "config path")
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			DoubleOutput(Fatal, "[DS]PanicError panic",
				"panic", r, "trace", string(debug.Stack()))
		}
	}()
	flag.Parse()
	Ctx = signalHandler()
	environment.InitEnv()
	err := config.InitConfig(*confPath)
	if err != nil {
		log.Fatalf("[DS]InitConfig failed,err=%v", err)
	}
	err = setCrashLog(config.Cfg.CommonConf.CrashLogPath)
	if err != nil {
		log.Fatalf("[DS]setCrashLog failed,err=%v", err)
	}
	if logger.Init(config.Cfg.Log) != nil {
		return
	}
	defer logger.Close()
	trace_http.Init()
	InitMetrics()
	if InitMysql() {
		return
	}
	if StartHttpServer() {
		return
	}
	DoubleOutput(Info, "[DS]Start http server successfully")
	<-Ctx.Done()
	DoubleOutput(Info, "[DS]Ready to close Server")
	recycle.ReleaseResources()

}
func StartHttpServer() bool {
	DoubleOutput(Info, "[DS]Ready to start http server !")
	if err := httpserver.Start(); err != nil {
		logger.NotCtxFatal("[DS]Start HttpServer Failed.", "error", err)
		return true
	}
	return false
}

func InitMysql() bool {
	if err := drivers.InitMysql(); err != nil {
		DoubleOutput(Fatal, "[IOT]Init mysql  failed err %v", err)
		return true
	}
	DoubleOutput(Info, "[DS]Init mysql success!")
	return false
}

func InitMetrics() {
	if err := metrics.Init(config.Cfg.CommonConf.ServerName); err != nil {
		DoubleOutput(Fatal, "[DS]Init monitor  failed ,err=%v", err)
	}
	DoubleOutput(Info, "[DS]Init metrics success!")
}

func signalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT,
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
		s := <-c
		logger.NotCtxInfof("getting signal for quit siginal %s", s)
		logger.NotCtxInfo("getting signal for quit", "siginal", s)
		cancel()
	}()
	return ctx
}

func setCrashLog(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	return syscall.Dup2(int(f.Fd()), 2)
}

func DoubleOutput(level Level, msg string, args ...any) {
	if args != nil {
		log.Printf(msg, args)
	} else {
		log.Printf(msg)
	}
	switch level {
	case Info:
		{
			if args != nil {
				logger.NotCtxInfof(msg, args)
			} else {
				logger.NotCtxInfof(msg)
			}
			return
		}
	case Fatal:
		{
			if args != nil {
				logger.NotCtxFatalf(msg, args)
			} else {
				logger.NotCtxFatalf(msg)
			}
			return
		}
	default:
		{
			if args != nil {
				logger.NotCtxInfof(msg, args)
			} else {
				logger.NotCtxInfof(msg)
			}
		}
	}
}
