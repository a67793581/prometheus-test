package httpserver

import (
	"fmt"
	"net/http"
	"prometheus-test/infrastructure/config"
	middleware2 "prometheus-test/server/httpserver/middleware"
	"time"

	"prometheus-test/lib/logger"
	"prometheus-test/lib/util"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	HttpSvr      *http.Server
}

func newHttpGinServer(port int, rTimeout int, wTimeout int) *HttpServer {
	server := &HttpServer{
		Port:         port,
		ReadTimeout:  rTimeout,
		WriteTimeout: wTimeout,
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(
		util.SetReqId(),
		middleware2.GinLogger())

	registerMetrics(engine)
	registerHealthDetect(engine)
	pprof.Register(engine, "/qnk8avm9pa/debug/pprof")

	registerCommonBizRouters(engine)

	server.HttpSvr = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      engine,
		ReadTimeout:  time.Duration(rTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(wTimeout) * time.Millisecond,
	}
	return server
}

func registerHealthDetect(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

func registerCommonBizRouters(engine *gin.Engine) {
	for _, r := range actionMaps {
		switch r.Method {
		case MethodGET:
			engine.GET(r.Path, r.Handlers...)
		case MethodPOST:
			engine.POST(r.Path, r.Handlers...)
		case MethodAll:
			engine.GET(r.Path, r.Handlers...)
			engine.POST(r.Path, r.Handlers...)
		}
	}
}

func registerMetrics(engine *gin.Engine) {
	engine.Use(middleware2.MonitorHandler())
}

func Start() error {
	httpConf := config.Cfg.ServerConf
	logger.NotCtxInfo("Start http server", "gport", httpConf.GPort)
	ginServer := newHttpGinServer(httpConf.GPort,
		httpConf.RTimeout, httpConf.WTimeout)
	err := gracehttp.Serve(ginServer.HttpSvr)
	return err
}
