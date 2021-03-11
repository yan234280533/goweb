package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
)

var responseLabels = prometheus.Labels{
	"resource": "resource1",
	"group":    "group1",
	"warning":  "2.5",
	"critical": "2.8",
}

var (
	appcount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "app",
		Name:      "request_total",
		Help:      "The total number of processed request",
	})

	appsummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "app",
			Name:       "request_summary",
			Help:       "This is my summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"service"})

	appResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "app",
			Name:      "response_time",
			Help:      "Finance Services http response time average over 1 minute",
		}, []string{"normal"})
)

func init() {

	//prometheus.MustRegister(appcount)
	prometheus.MustRegister(appsummary)
	prometheus.MustRegister(appResponseTime)

	// Add Go module build info.
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func myWeb(w http.ResponseWriter, r *http.Request) {
	//appcount.Inc()

	cost := float64(rand.Intn(100)*1.0) / 100.0
	fmt.Fprintf(w, fmt.Sprintf("%.2f\n", cost))

	appResponseTime.WithLabelValues("normal").Set(float64(cost))
	appsummary.WithLabelValues("service").Observe(float64(cost))

	fmt.Fprintf(w, fmt.Sprintf("这%.2f\n", cost))
}

func main() {
	http.HandleFunc("/", myWeb)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("服务器即将开启，访问地址 http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器开启错误: ", err)
	}
}
