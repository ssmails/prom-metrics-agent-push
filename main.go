package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(1 * time.Second)
		}
	}()
}

func main() {
	log.Println("starting API server")

	prometheus.MustRegister(respTime)
	prometheus.MustRegister(reqProcessed)
	prometheus.MustRegister(reqSize)

	recordMetrics()

	router := mux.NewRouter()
	log.Println("creating routes")

	router.HandleFunc("/agentmetrics", AgentMetrics).Methods("POST")
	router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8081", router)
}

////////////////////////////////////

type Metric struct {
	RespTime float64 `json:"resp_time"`
	ReqCnt   int     `json:"req_cnt"`
	ReqSize  int     `json:"req_size"`
}

var (
	//app metrics
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	//agent metrics
	labels = []string{"tenantid", "agentid"}

	respTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "pod_container_resource_requests",
			Name: "resp_time_seconds",
			Help: "The resp time for tenant svc req.",
		},
		labels)

	reqProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "req_total",
			Help: "The total number of processed events",
		},
		labels)

	reqSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "req_size",
			Help: "The req size",
		},
		labels)
)

func AgentMetrics(w http.ResponseWriter, r *http.Request) {
	var metric Metric

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&metric); err != nil {
		log.Printf("%v", err)
		return
	}
	defer r.Body.Close()
	log.Printf("Got req Body:%v", metric)

	metricLabels := prometheus.Labels{
		"tenantid": "tenant1",
		"agentid":  "agent1",
	}
	respTime.With(metricLabels).Set(metric.RespTime)
	//respTime.WithLabelValues("tenant1", "agent1").Set(metric.RespTime)

	reqProcessed.With(metricLabels).Add(float64(metric.ReqCnt))
	//reqProcessed.WithLabelValues("tenant1", "agent1").Add(float64(metric.ReqCnt))

	reqSize.With(metricLabels).Observe(float64(metric.ReqSize))

	w.WriteHeader(http.StatusOK)
}

//https://antonputra.com/monitoring/monitor-golang-with-prometheus/#gauge
