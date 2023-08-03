package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	MethodGET = iota
	MethodPOST
	MethodAll
)

type action struct {
	Path     string
	Method   int
	Handlers []gin.HandlerFunc
}

var actionMaps = make([]*action, 0)

func registerGinHttpAction(path string, method int, handlers ...gin.HandlerFunc) {
	if len(handlers) == 0 {
		panic("action no handlers!")
	}
	actionMaps = append(actionMaps, &action{
		Path:     path,
		Method:   method,
		Handlers: handlers,
	})
}

func init() {
	registerGinHttpAction("/metrics", MethodGET, func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	registerGinHttpAction("/status.php", MethodGET, func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok\n")
	})

}
