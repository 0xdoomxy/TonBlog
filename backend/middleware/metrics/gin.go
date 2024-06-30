package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func BindMetrics(engine *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(engine)
}

func AddMetric(metric *ginmetrics.Metric) error {
	return ginmetrics.GetMonitor().AddMetric(metric)
}

func GetMetric(name string) *ginmetrics.Metric {
	return ginmetrics.GetMonitor().GetMetric(name)
}
