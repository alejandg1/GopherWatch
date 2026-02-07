package processor

import (
	"sync"
)

// Engine gestiona el estado actual de todos los contenedores
type Engine struct {
	// escritura concurrente segura con sync.Map
	Containers sync.Map
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) UpdateContainer(stats ContainerStats) {
	e.Containers.Store(stats.ID, stats)
}

func (e *Engine) RemoveContainer(id string) {
	e.Containers.Delete(id)
}
