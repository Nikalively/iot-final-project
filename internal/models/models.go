package models

import (
	"time"
)

type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float64   `json:"cpu"`
	RPS       float64   `json:"rps"`
}

type AnalyticsResult struct {
	SmoothedLoad float64 `json:"smoothed_load"`
	AnomalyCount int     `json:"anomaly_count"`
}
