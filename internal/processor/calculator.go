package processor

import "log"

func CalculateCPUPercent(v *StatsJSON) float64 {
	// Diferencia en el uso del contenedor
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
	// Diferencia en el uso del sistema
	systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)

	cpuPercent := 0.0
	if cpuDelta > 0.0 && systemDelta > 0.0 {
		onlineCPUs := float64(v.CPUStats.OnlineCPUs)
		if onlineCPUs == 0.0 {
			onlineCPUs = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
		}
		if onlineCPUs == 0.0 {
			onlineCPUs = 1.0 // Fallback to 1 core if unknown to avoid zero
		}

		// Calculamos basándonos en el número de núcleos activos
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	// Debug log
	log.Printf("CPU Debug - Delta: %f, SystemDelta: %f, Percent: %f, Cores (Percpu): %d, OnlineCPUs: %d", cpuDelta, systemDelta, cpuPercent, len(v.CPUStats.CPUUsage.PercpuUsage), v.CPUStats.OnlineCPUs)
	return cpuPercent
}

func CalculateMemUsage(v *StatsJSON) float64 {
	// Convertimos bytes a MiB para que sea legible en tu eww o dashboard
	return float64(v.MemoryStats.Usage) / (1024 * 1024)
}
