package processor

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/docker/docker/client"
)

// escucha el stream del socket y actualiza el Engine
func (e *Engine) ProcessContainerMetrics(ctx context.Context, cli *client.Client, containerID string, name string) {
	statsStream, err := cli.ContainerStats(ctx, containerID, true)
	if err != nil {
		log.Printf("Error obteniendo stats para %s: %v", name, err)
		return
	}
	defer statsStream.Body.Close()

	decoder := json.NewDecoder(statsStream.Body)
	var v StatsJSON

	for {
		select {
		// manejo reactivo en lugar de if bloqueante
		case <-ctx.Done():
			return
		default:
			// Decodificamos el JSON que viene del stream del socket
			if err := decoder.Decode(&v); err != nil {
				if err == io.EOF {
					return
				}
				continue
			}

			s := ContainerStats{
				ID:       containerID,
				Name:     name,
				Status:   "running",
				CPUUsage: CalculateCPUPercent(&v),
				MemUsage: CalculateMemUsage(&v),
				MemLimit: float64(v.MemoryStats.Limit) / (1024 * 1024),
			}
			s.MemPerc = (s.MemUsage / s.MemLimit) * 100

			e.UpdateContainer(s)
		}
	}
}
