package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/actiometa/gopherwatch/internal/collector"
	"github.com/actiometa/gopherwatch/internal/processor"
	"github.com/actiometa/gopherwatch/internal/server"
	"github.com/docker/docker/api/types/container"
)

func main() {
	log.Println("Starting")

	cli, err := collector.NewDockerClient()
	if err != nil {
		log.Fatalf("Error al conectar con el socket: %v", err)
	}
	defer cli.Close()

	info, err := cli.Info(context.Background())
	if err != nil {
		log.Fatalf("No se pudo obtener información: %v", err)
	}
	fmt.Printf("Conectado a: %s (Versión: %s)\n", info.Name, info.ServerVersion)

	//  Inicializar motor y contexto
	engine := processor.NewEngine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Iniciar monitoreo de contenedores existentes
	containers, err := cli.ContainerList(ctx, container.ListOptions{}) // Nota: Cambiado a container.ListOptions
	if err != nil {
		log.Printf("Error listando contenedores: %v", err)
	}

	for _, c := range containers {
		log.Printf("Monitoreando: %s", c.Names[0])
		go engine.ProcessContainerMetrics(ctx, cli, c.ID, c.Names[0])
	}
	// Iniciar monitoreo de eventos (start/stop)
	collector.StartEventWatcher(ctx, cli, engine)

	srv := server.NewApiServer(engine)
	go func() {
		if err := srv.Start("8080"); err != nil {
			log.Printf("Error en el servidor HTTP: %v", err)
		}
	}()

	// Manejo de señales
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Presiona Ctrl+C para salir.")

	<-stop
	log.Println("\nStopping GopherWatch")
}
