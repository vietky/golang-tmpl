package echoprometheus

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"strconv"
	"time"
)

var defaultMetricPath = "/metrics"

// Prometheus contains the metrics gathered by the instance and its path
type Prometheus struct {
	reqCnt *prometheus.CounterVec
	reqDur *prometheus.HistogramVec

	MetricsPath string
}

// NewPrometheus generates a new set of metrics with a certain subsystem name
func NewPrometheus(subsystem string) *Prometheus {
	p := &Prometheus{
		MetricsPath: defaultMetricPath,
	}

	p.registerMetrics(subsystem)

	return p
}

func (p *Prometheus) registerMetrics(subsystem string) {

	// p.reqCnt = prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Subsystem: subsystem,
	// 		Name:      "requests_total",
	// 		Help:      "How many HTTP requests processed, partitioned by status code and HTTP path.",
	// 	},
	// 	[]string{"code", "path"},
	// )
	// prometheus.MustRegister(p.reqCnt)

	p.reqDur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
			Buckets:   []float64{.005, .01, .02, 0.04, .06, 0.08, .1, 0.15, .25, 0.4, .6, .8, 1, 1.5, 2, 3, 5},
		},
		[]string{"code", "path"},
	)
	prometheus.MustRegister(p.reqDur)

}

// Use adds the middleware to a echo enechoe.
func (p *Prometheus) Use(e *echo.Echo) {
	e.Use(p.middlewareFunc())
	e.GET(p.MetricsPath, prometheusHandler())
}

func (p *Prometheus) middlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			if c.Request().URL.String() == p.MetricsPath {
				return next(c)
			}

			start := time.Now()

			ctxPath := c.Path()

			if err = next(c); err != nil {
				c.Error(err)
				ctxPath = "/404"
			}

			status := strconv.Itoa(c.Response().Status)
			elapsed := float64(time.Since(start)) / float64(time.Second)

			path := c.Request().Method + "_" + ctxPath
			p.reqDur.WithLabelValues(status, path).Observe(elapsed)

			//p.reqCnt.WithLabelValues(status, path).Inc()
			c.Response().Header().Set("X-Response-Time", strconv.Itoa(int(elapsed*1000)))

			return
		}
	}
}

func prometheusHandler() echo.HandlerFunc {
	h := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			// EnableOpenMetrics: true,
		},
	)
	return echo.WrapHandler(h)
}
