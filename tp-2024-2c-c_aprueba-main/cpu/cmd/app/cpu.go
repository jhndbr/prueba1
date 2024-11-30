package main

import (
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/http"
	"log/slog"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	slog.SetLogLoggerLevel(slog.LevelInfo)
	conf := config.GetInstance()
	if conf == nil {
		slog.Error("Failed to get config instance")
		return
	}



	wg.Add(1) // Añade una goroutine al WaitGroup antes de lanzar la goroutine
	go func(port int) {
		err := http.Server("127.0.0.1", port)
		if err != nil {
			slog.Error("CPU", slog.String("error", err.Error()))
			return
		}
		defer wg.Done()
	}(conf.Port)

	wg.Wait()

}
