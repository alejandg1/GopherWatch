package collector

import (
	"context"
	"log"

	"github.com/actiometa/gopherwatch/internal/processor"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

func StartEventWatcher(ctx context.Context, cli *client.Client, engine *processor.Engine) {
	options := events.ListOptions{}

	msgs, errs := cli.Events(ctx, options)

	// Goroutine para manejar eventos de contenedores
	go func() {
		for {
			// Escuchar eventos de goroutine
			select {
			case err := <-errs:
				if err != nil {
					log.Printf("Error en el stream de eventos: %v", err)
					return
				}
			case msg := <-msgs:
				handleEvent(ctx, cli, engine, msg)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func handleEvent(ctx context.Context, cli *client.Client, engine *processor.Engine, msg events.Message) {
	if msg.Type != "container" {
		return
	}

	containerID := msg.Actor.ID
	containerName := msg.Actor.Attributes["name"]

	switch msg.Action {
	case "start":
		log.Printf("Nuevo contenedor detectado: %s", containerName)
		go engine.ProcessContainerMetrics(ctx, cli, containerID, containerName)
	case "die", "stop":
		log.Printf("Contenedor detenido: %s", containerName)
		engine.RemoveContainer(containerID)
	}
}

// pide el stream de datos al socket
func GetContainerStats(ctx context.Context, cli *client.Client, containerID string) (container.StatsResponseReader, error) {
	// Esto devuelve un stream que se actualiza cada segundo por defecto
	return cli.ContainerStats(ctx, containerID, true)
}
