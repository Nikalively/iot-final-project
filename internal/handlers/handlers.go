package handlers

import (
	"encoding/json"
	"github.com/Nikalively/iot-final-project/internal/analytics"
	"github.com/Nikalively/iot-final-project/internal/config"
	"github.com/Nikalively/iot-final-project/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
	"time"
)

var (
	rpsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "iot_rps_total",
		Help: "Total RPS processed",
	})
	anomalyCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "iot_anomalies_total",
		Help: "Total anomalies detected",
	})
	latencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "iot_request_duration_seconds",
		Help: "Request latency",
	})
)

func init() {
	prometheus.MustRegister(rpsCounter, anomalyCounter, latencyHistogram)
}

type Handler struct {
	analyzer *analytics.Analyzer
	redis    *redis.Client
	mu       sync.Mutex
}

func NewHandler(cfg *config.Config) *Handler {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	return &Handler{
		analyzer: analytics.NewAnalyzer(),
		redis:    rdb,
	}
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(latencyHistogram)
	defer timer.ObserveDuration()

	var metric models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rpsCounter.Inc()
	
	key := "metric:" + metric.Timestamp.Format(time.RFC3339)
	data, _ := json.Marshal(metric)
	h.redis.Set(r.Context(), key, data, 5*time.Minute)

	ch := make(chan models.AnalyticsResult, 1)
	go h.analyzer.ProcessMetric(metric, ch)
	result := <-ch
	if result.AnomalyCount > 0 {
		anomalyCounter.Add(float64(result.AnomalyCount))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	smoothed := h.analyzer.GetSmoothedLoad()
	anomalies := h.analyzer.GetAnomalyCount()
	json.NewEncoder(w).Encode(models.AnalyticsResult{
		SmoothedLoad: smoothed,
		AnomalyCount: anomalies,
	})
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func SetupRoutes(cfg *config.Config) *mux.Router {
	h := NewHandler(cfg)
	r := mux.NewRouter()
	r.HandleFunc("/metrics", h.MetricsHandler).Methods("POST")
	r.HandleFunc("/analyze", h.AnalyzeHandler).Methods("GET")
	r.HandleFunc("/health", h.HealthHandler).Methods("GET")
	r.Handle("/metrics/prometheus", promhttp.Handler())
	return r
}
