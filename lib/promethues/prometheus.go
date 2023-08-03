package prometheus

import (
	"errors"
	"net/http"
	"runtime/debug"
	"sync"

	"prometheus-test/lib/logger"
	"prometheus-test/lib/util"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Type int

const (
	TypeQPS Type = iota
	TypeTotal
	TypeSummary
)

var DefaultBuckets = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 5000, 10000, 50000, 100000, 500000}

type prometheusInner struct {
	serverName   string
	idc          string
	ip           string
	qpsVec       map[string]*prometheus.CounterVec   //sum count
	totalVec     map[string]*prometheus.HistogramVec //sum count bulk
	summaryVec   map[string]*prometheus.SummaryVec   //summary vec
	qpsMutex     sync.Mutex
	totalMutex   sync.Mutex
	summaryMutex sync.Mutex
}

var inner *prometheusInner

func newPrometheusInner(name string, idcName string) *prometheusInner {
	ipStr, err := util.GetInternalIP()
	if err != nil {
		ipStr = "0.0.0.0"
	}
	ins := &prometheusInner{
		serverName: name,
		idc:        idcName,
		ip:         ipStr,
		qpsVec:     make(map[string]*prometheus.CounterVec),
		totalVec:   make(map[string]*prometheus.HistogramVec),
		summaryVec: make(map[string]*prometheus.SummaryVec),
	}
	return ins
}

func (pI *prometheusInner) registeQps(name string, labels []string) {
	pI.qpsMutex.Lock()
	defer pI.qpsMutex.Unlock()
	if _, ok := pI.qpsVec[name]; !ok {
		pI.qpsVec[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "dispatcher",
			Subsystem: pI.serverName,
			Name:      name,
			Help:      "dispatcher qps",
		}, append(labels, "idc", "ip"))
		prometheus.MustRegister(pI.qpsVec[name])
	}

}

func (pI *prometheusInner) registeTotal(name string, labels []string, bulks []float64) {
	pI.totalMutex.Lock()
	defer pI.totalMutex.Unlock()
	if _, ok := pI.totalVec[name]; !ok {

		opts := prometheus.HistogramOpts{
			Namespace: "dispatcher",
			Subsystem: pI.serverName,
			Name:      name,
			Help:      "dispatcher total",
		}
		if bulks != nil {
			opts.Buckets = bulks
		} else {
			opts.Buckets = DefaultBuckets
		}
		pI.totalVec[name] = prometheus.NewHistogramVec(opts, append(labels, "idc", "ip"))
		prometheus.MustRegister(pI.totalVec[name])
	}

}

func (pI *prometheusInner) registerSummary(name string, labels []string, bulks []float64) {
	pI.totalMutex.Lock()
	defer pI.totalMutex.Unlock()
	if _, ok := pI.summaryVec[name]; !ok {

		opts := prometheus.SummaryOpts{
			Namespace:  "dispatcher",
			Subsystem:  pI.serverName,
			Name:       name,
			Objectives: map[float64]float64{0.5: 0.5, 0.9: 0.1, 0.99: 0.01},
		}

		pI.summaryVec[name] = prometheus.NewSummaryVec(opts, append(labels, "idc", "ip"))
		prometheus.MustRegister(pI.summaryVec[name])
	}

}

func (pI *prometheusInner) incQps(key string, kv map[string]string) (err error) {
	pI.qpsMutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			_ = errors.New("check labels")
			logger.NotCtxErrorf("stack:%v", string(debug.Stack()))
		}
	}()
	defer pI.qpsMutex.Unlock()
	if qpsP, ok := pI.qpsVec[key]; ok {
		kv["ip"] = pI.ip
		kv["idc"] = pI.idc
		qpsP.With(kv).Inc()
	} else {
		err = errors.New("not corret name,please check")
	}
	return
}

func (pI *prometheusInner) updateTotal(key string, kv map[string]string, value float64) (err error) {
	pI.totalMutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			_ = errors.New("check labels")
			logger.NotCtxErrorf("stack:%v", string(debug.Stack()))
		}
	}()
	defer pI.totalMutex.Unlock()
	if totalP, ok := pI.totalVec[key]; ok {
		kv["ip"] = pI.ip
		kv["idc"] = pI.idc
		totalP.With(kv).Observe(value)
	} else {
		err = errors.New("not corret name,please check")
	}
	return
}

func (pI *prometheusInner) updateSummary(key string, kv map[string]string, value float64) (err error) {
	pI.summaryMutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			_ = errors.New("check labels")
			logger.NotCtxErrorf("stack:%v", string(debug.Stack()))
		}
	}()
	defer pI.summaryMutex.Unlock()
	if summaryIP, ok := pI.summaryVec[key]; ok {
		kv["ip"] = pI.ip
		kv["idc"] = pI.idc
		summaryIP.With(kv).Observe(value)
	} else {
		err = errors.New("not corret name,please check")
	}
	return
}

func Init(serverName string, idc string) {
	inner = newPrometheusInner(serverName, idc)
}

func Registe(pType Type, name string, labels []string, bulks []float64) {
	switch pType {
	case TypeQPS:
		inner.registeQps(name, labels)
	case TypeTotal:
		inner.registeTotal(name, labels, bulks)
	case TypeSummary:
		inner.registerSummary(name, labels, bulks)
	}
}

func Update(pType Type, name string, kv map[string]string, value float64) error {
	switch pType {
	case TypeQPS:
		return inner.incQps(name, kv)
	case TypeTotal:
		return inner.updateTotal(name, kv, value)
	case TypeSummary:
		return inner.updateSummary(name, kv, value)
	}
	return nil
}

func Inc(pType Type, name string, kv map[string]string) error {
	switch pType {
	case TypeQPS:
		return inner.incQps(name, kv)
	case TypeTotal:
		return inner.updateTotal(name, kv, 1)
	case TypeSummary:
		return inner.updateSummary(name, kv, 1) //
	}
	return nil
}

func NewHttpHander() http.Handler {
	return promhttp.Handler()
}
