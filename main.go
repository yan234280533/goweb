package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"syscall"
	"time"
)

var responseLabels = prometheus.Labels{
	"resource": "resource1",
	"group":    "group1",
	"warning":  "2.5",
	"critical": "2.8",
}

const MIN = 0.000000001

func IsEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < MIN
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

	app2summary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "app",
			Name:       "request2_summary",
			Help:       "This is my2 summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"service"})

	/*httpsummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "happ",
		Name: "happ_response_time_summary",
		Help: "Duration of HTTP requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01,0.95: 0.005,0.99: 0.001}, //返回五分数， 九分数， 九九分数
	}, []string{"service"})*/

	appResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "app",
			Name:      "response_time",
			Help:      "Finance Services http response time average over 1 minute",
		}, []string{"normal"})
)

func doCpuWork(cpuSize int, count int) {
	var wg sync.WaitGroup
	var timeNow uint64
	for c := 0; c < cpuSize; c++ {
		wg.Add(1)
		go func() {
			for i := 0; i < count; i++ {
				time.Sleep(1 * time.Nanosecond)
				timeNow++
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fd, _, err := syscall.RawSyscall(syscall.SYS_SOCKET, syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
}

func init() {
	//prometheus.MustRegister(appcount)
	prometheus.MustRegister(appsummary)
	prometheus.MustRegister(app2summary)
	prometheus.MustRegister(appResponseTime)

	// Add Go module build info.
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func myWeb(w http.ResponseWriter, r *http.Request) {
	//appcount.Inc()
	cost := myfunc()
	fmt.Fprintf(w, fmt.Sprintf("%.2f\n", cost))
}

func mySqrt(w http.ResponseWriter, r *http.Request) {
	startT := time.Now()
	doCpuWork(1, 4000000)
	tc := time.Since(startT)
	app2summary.WithLabelValues("service").Observe(float64(tc.Nanoseconds()) / 1000.0 / 1000.0)
}

func myfunc() float64 {
	cost := float64(rand.Intn(100)*1.0) / 100.0
	appResponseTime.WithLabelValues("normal").Set(float64(cost))
	appsummary.WithLabelValues("service").Observe(float64(cost))
	return cost
}

func myQuantile(w http.ResponseWriter, r *http.Request) {
	metric, err := appsummary.MetricVec.GetMetricWith(prometheus.Labels{"service": "service"})
	if err != nil {
		fmt.Fprintf(w, "err %s", err.Error())
		return
	} else {
		var dtoMetric = dto.Metric{}

		err = metric.Write(&dtoMetric)
		if err != nil {
			fmt.Fprintf(w, "err %s", err.Error())
			return
		}

		var value float64
		for _, v := range dtoMetric.Summary.Quantile {
			if !IsEqual(*(v.Quantile), 0.99) {
				continue
			} else {
				value = *(v.Value)
				break
			}
		}

		fmt.Fprintf(w, fmt.Sprintf("%v\n", value))
	}
}

func main() {
	http.HandleFunc("/", myWeb)
	http.HandleFunc("/sqrt", mySqrt)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/quantile", myQuantile)

	go func() {
		fmt.Println("Loop of my func")
		for {
			time.Sleep(10 * time.Millisecond)
			//_ = myfunc()
		}
	}()

	fmt.Println("服务器即将开启，访问地址 http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器开启错误: ", err)
	}
}
