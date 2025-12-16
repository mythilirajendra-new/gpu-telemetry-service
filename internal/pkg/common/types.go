package common

import "time"

type Telemetry struct {
	ID          string    `json:"id"`
	HostID      string    `json:"host_id"`
	GPUId       string    `json:"gpu_id"`
	Utilization float64   `json:"utilization"`
	MemoryUsed  float64   `json:"memory_used"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}
