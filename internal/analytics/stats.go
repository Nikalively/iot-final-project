package analytics

import (
	"math"
	"sync"
	"github.com/Nikalively/iot-final-project/internal/models"
)

const WindowSize = 50
const ZThreshold = 2.0

type Analyzer struct {
	window []float64
	mu     sync.RWMutex
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{window: make([]float64, 0, WindowSize)}
}

func (a *Analyzer) ProcessMetric(metric models.Metric, ch chan<- models.AnalyticsResult) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.window = append(a.window, metric.RPS)
	if len(a.window) > WindowSize {
		a.window = a.window[1:]
	}

	smoothed := a.calculateRollingAvg()

	anomalyCount := a.detectAnomalies()

	ch <- models.AnalyticsResult{
		SmoothedLoad: smoothed,
		AnomalyCount: anomalyCount,
	}
}

func (a *Analyzer) calculateRollingAvg() float64 {
	if len(a.window) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range a.window {
		sum += v
	}
	return sum / float64(len(a.window))
}

func (a *Analyzer) detectAnomalies() int {
	if len(a.window) < 2 {
		return 0
	}
	mean, std := a.calculateMeanStd()
	count := 0
	for _, v := range a.window {
		z := math.Abs((v - mean) / std)
		if z > ZThreshold {
			count++
		}
	}
	return count
}

func (a *Analyzer) calculateMeanStd() (float64, float64) {
	sum := 0.0
	for _, v := range a.window {
		sum += v
	}
	mean := sum / float64(len(a.window))
	variance := 0.0
	for _, v := range a.window {
		variance += math.Pow(v-mean, 2)
	}
	std := math.Sqrt(variance / float64(len(a.window)))
	return mean, std
}

func (a *Analyzer) GetSmoothedLoad() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.calculateRollingAvg()
}

func (a *Analyzer) GetAnomalyCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.detectAnomalies()
}
