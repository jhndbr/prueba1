package main

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/http"
)

func main() {
	var wg sync.WaitGroup
	slog.SetLogLoggerLevel(slog.LevelInfo)
	conf := config.GetInstance()

	if conf == nil {
		slog.Error("Error: la configuración no se cargó correctamente.")
		return
	}

	slog.Debug("->", "puerto", conf.Port)
	slog.Debug("->", "Search:", conf.MemorySize)
	slog.Debug("->", "Algorithm:", conf.SearchAlgorithm)

	entity.InicializarMemoriaUsuario(conf.MemorySize)
	entity.InicializarMemoriaSistema(conf.MemorySize)

	if conf.Scheme == "FIJAS" {
		usecase.InicializarParticion(conf.Partitions)
	} else if conf.Scheme == "DINAMICAS" {
		slog.Debug("Inicializando memoria dinámica")
		usecase.InicializarMemoriaDinamica(conf.MemorySize)

	}

	addr := fmt.Sprintf("localhost:%d", conf.Port)
	slog.Debug("Starting Up", "Address", addr)

	wg.Add(1)
	go func() {
		err := http.Server("localhost", conf.Port)
		if err != nil {
			slog.Error("Error to start Server", "Port", conf.Port)
		}
	}()
	wg.Wait()
}
