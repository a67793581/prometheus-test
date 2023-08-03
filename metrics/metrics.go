package metrics

import (
	"fmt"
	"time"

	"prometheus-test/lib/gomonitor"
	"prometheus-test/lib/logger"
	prometheus "prometheus-test/lib/promethues"
	"prometheus-test/lib/util"
)

const (
	MonitorNameInterface          = "interface"
	MonitorNameInterfaceQps       = "interface_qps"
	MonitorNameInterfaceErrorCode = "interface_code"
	MonitorNameDependence         = "dependence"
	MonitorNameDependenceQps      = "dependence_qps"
	MonitorNameStatistics         = "statistics"

	MonitorNameDb    = "DB"
	MonitorNameDbQps = "DB_qps"
)

func Init(srvName string) error {
	prometheus.Init(srvName, util.IdcName())
	register()
	go monitor()
	return nil

}

func register() {
	prometheus.Registe(prometheus.TypeSummary, MonitorNameInterface, []string{"interface", "status"}, nil)
	prometheus.Registe(prometheus.TypeQPS, MonitorNameInterfaceQps, []string{"interface", "status"}, nil)
	prometheus.Registe(prometheus.TypeQPS, MonitorNameInterfaceErrorCode, []string{"interface", "status"}, nil)

	prometheus.Registe(prometheus.TypeSummary, MonitorNameDependence, []string{"dependence_service", "function", "status"}, nil)
	prometheus.Registe(prometheus.TypeQPS, MonitorNameDependenceQps, []string{"dependence_service", "function", "status"}, nil)

	prometheus.Registe(prometheus.TypeTotal, MonitorNameStatistics, []string{"module"}, nil)

	prometheus.Registe(prometheus.TypeSummary, MonitorNameDb, []string{"table", "function", "status"}, nil)
	prometheus.Registe(prometheus.TypeQPS, MonitorNameDbQps, []string{"table", "function", "status"}, nil)
}

func monitor() {

	var interval int64 = 10
	for {
		stat := gomonitor.GetState()
		logger.NotCtxInfo("[Monitor]", "MEMStat", util.StructToJson(stat))
		UpdateStatistics("GO_GCNum", int64(stat.GCNum))
		UpdateStatistics("GO_GCPause", int64(stat.GCPause))
		UpdateStatistics("GO_MemStack", int64(stat.MemStack))
		UpdateStatistics("GO_MemMallocs", int64(stat.MemMallocs))
		UpdateStatistics("GO_MemAllocated", int64(stat.MemAllocated))
		UpdateStatistics("GO_MemObjects", int64(stat.MemObjects))
		UpdateStatistics("GO_MemHeap", int64(stat.MemHeap))
		UpdateStatistics("GO_GoroutineNum", int64(stat.GoroutineNum))
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
func UpdateInterface(method string, status int, value int64) {
	labels := map[string]string{
		"interface": method,
		"status":    fmt.Sprintf("%v", status),
	}
	err := prometheus.Update(prometheus.TypeSummary, MonitorNameInterface, labels, float64(value))
	if err != nil {
		logger.NotCtxInfof("prometheus.Update UpdateInterface failed,err=%v", err)
	}
}

func UpdateInterfaceQPS(method string, status int, value int64) {
	labels := map[string]string{
		"interface": method,
		"status":    fmt.Sprintf("%v", status),
	}
	err := prometheus.Update(prometheus.TypeQPS, MonitorNameInterfaceQps, labels, float64(value))
	if err != nil {
		logger.NotCtxInfof("prometheus.Update UpdateInterfaceQPS failed,err=%v", err)
	}
}

func UpdateInterfaceQPSByErrorCode(method string, status int, value int64) {
	labels := map[string]string{
		"interface": method,
		"status":    fmt.Sprintf("%v", status),
	}
	err := prometheus.Update(prometheus.TypeQPS, MonitorNameInterfaceErrorCode, labels, float64(value))
	if err != nil {
		logger.NotCtxInfof("prometheus.Update UpdateInterfaceQPSByErrorCode failed,err=%v", err)
	}
}

func UpdateDependence(service, function string, value int64, err error) {
	labels := map[string]string{
		"function":           function,
		"dependence_service": service,
	}
	if err != nil {
		labels["status"] = "1"
	} else {
		labels["status"] = "0"
	}
	err2 := prometheus.Update(prometheus.TypeSummary, MonitorNameDependence, labels, float64(value))
	if err2 != nil {
		logger.NotCtxInfof("prometheus.Update UpdateDependence failed,err=%v", err2)
	}
}

func UpdateDependenceQPS(service, method string, status int, value int64) {
	labels := map[string]string{
		"function":           method,
		"dependence_service": service,
		"status":             fmt.Sprintf("%v", status),
	}
	err := prometheus.Update(prometheus.TypeQPS, MonitorNameDependenceQps, labels, float64(value))
	if err != nil {
		logger.NotCtxInfof("prometheus.Update UpdateDependenceQPS failed,err=%v", err)
	}
}

func UpdateStatistics(module string, value int64) {
	labels := map[string]string{
		"module": module,
	}
	err := prometheus.Update(prometheus.TypeTotal, MonitorNameStatistics, labels, float64(value))
	if err != nil {
		logger.NotCtxInfof("prometheus.Update UpdateStatistics failed,err=%v,%+v,%+v", err, module, value)
	}
}

func UpdateDB(table, method string, value int64, err error) {
	labels := map[string]string{
		"function": method,
		"table":    table,
	}
	if err != nil {
		labels["status"] = "1"
	} else {
		labels["status"] = "0"
	}
	err2 := prometheus.Update(prometheus.TypeSummary, MonitorNameDb, labels, float64(value))
	if err2 != nil {
		logger.NotCtxInfof("prometheus.Update UpdateDB failed,err=%v", err2)
	}
}

func UpdateDBQPS(table, method string, status error, value int64) {
	labels := map[string]string{
		"function": method,
		"table":    table,
		"status":   fmt.Sprintf("%v", status == nil),
	}
	err := prometheus.Update(prometheus.TypeQPS, MonitorNameDbQps, labels, float64(value))
	if err != nil {
		logger.NotCtxInfof("prometheus.Update UpdateDBQPS failed,err=%v", err)
	}
}
