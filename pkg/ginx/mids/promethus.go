package mids

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type PrometheusBuilder struct {
	NameSpace     string
	SubSystem     string
	Name          string
	DurationHelp  string
	ActiveReqHelp string
	InstanceId    string
}

func NewPrometheusBuilder(namespace, subSystem, name string) *PrometheusBuilder {
	return &PrometheusBuilder{
		NameSpace: namespace,
		SubSystem: subSystem,
		Name:      name,
	}
}

func (p *PrometheusBuilder) SetDurationHelp(help string) *PrometheusBuilder {
	p.DurationHelp = help
	return p
}

func (p *PrometheusBuilder) SetActiveReqHelp(help string) *PrometheusBuilder {
	p.ActiveReqHelp = help
	return p
}
func (p *PrometheusBuilder) SetInstanceId(id string) *PrometheusBuilder {
	p.InstanceId = id
	return p
}

func (p *PrometheusBuilder) Build() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: p.NameSpace,
		Subsystem: p.SubSystem,
		Name:      p.Name + "_resp_time",
		Help:      p.DurationHelp,
		ConstLabels: map[string]string{
			"instance_id": p.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(summary)

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: p.NameSpace,
		Subsystem: p.SubSystem,
		Name:      p.Name + "_active_req",
		Help:      p.ActiveReqHelp,
		ConstLabels: map[string]string{
			"instance_id": p.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)

	return func(ctx *gin.Context) {
		start := time.Now()
		gauge.Inc()
		defer func() {
			gauge.Dec()
			// fullPath like /detail/:id
			pattern := ctx.FullPath()
			if pattern == "" {
				pattern = "unknown"
			}
			summary.WithLabelValues(
				ctx.Request.Method,
				pattern,
				strconv.Itoa(ctx.Writer.Status()),
			).Observe(float64(time.Since(start).Milliseconds()))
		}()
		ctx.Next()
	}
}
