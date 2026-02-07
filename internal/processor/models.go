package processor

type ContainerStats struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	CPUUsage float64 `json:"cpu_percentage"`
	MemUsage float64 `json:"mem_usage_mb"`
	MemLimit float64 `json:"mem_limit_mb"`
	MemPerc  float64 `json:"mem_percentage"`
}
