package main

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/usecase/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/infra/http"
)

func gateway(conf *config.Config, wg *sync.WaitGroup) {
	defer wg.Done() // Marca como completado al finalizar
	err := http.Server("127.0.0.1", conf.Port)
	if err != nil {
		slog.Error("Server error", "error", err)
		return
	}
}

func console() (filePath, size string) {

	fmt.Println("Indique el numero de prueba a hacer:")

	fmt.Println()
	fmt.Println("1:Planificacion")
	fmt.Println("2:Race Condition")
	fmt.Println("3:Particiones Fijas")
	fmt.Println("4:Particiones Dinamicas")
	fmt.Println("5:Fibonacci")
	fmt.Println("6:Stress")

	var prueba string

	_, err1 := fmt.Scan(&prueba)

	switch prueba {
	case "1":
		filePath = "PLANI_PROC"
		size = "32"
	case "2":
		filePath = "RECURSOS_MUTEX_PROC"
		size = "32"
	case "3":
		filePath = "MEM_FIJA_BASE"
		size = "12"
	case "4":
		filePath = "MEM_DINAMICA_BASE"
		size = "128"
	case "5":
		filePath = "PRUEBA_FS"
		size = "0"
	case "6":
		filePath = "THE_EMPTINESS_MACHINE"
		size = "16"
	}

	if err1 != nil {
		fmt.Println("Error:", err1)
		return "", ""
	}

	return filePath, size

}

func main() {
	var wg sync.WaitGroup
	slog.SetLogLoggerLevel(slog.LevelInfo)
	var filePath, size string
	filePath = "RECURSOS_MUTEX_PROC"
	size = "32"
	//filePath, size = console()

	planificador.InicializarSemaforos()
	planificador.VerificarSemaforos()

	slog.Info("Loading configuration", "file:", filePath)
	slog.Info("Loading configuration", "size:", size)

	// Inicio Planificador de Large Plazo.
	go planificador.AtenderProcesosNuevos()

	// Inicio Planificador de Corto Plazo.
	go planificador.AtenderProcesosListos()

	// Inicio Cola de IO
	go planificador.ProcessIOQueue()

	//wg.Wait() // Espera a que todas las solicitudes de I/O se completen

	// Cargo Proceso por defecto
	sizeOfProcess, _ := strconv.Atoi(size)
	err := planificador.CrearProcesoNEW(filePath, uint32(sizeOfProcess), 0)
	if err != nil {
		return
	}

	// Inicio el Gateway de entrada
	conf := config.GetInstance()

	wg.Add(1) // Añade una goroutine al WaitGroup antes de lanzar la goroutine
	go gateway(conf, &wg)

	// Espero a que todas las goroutines terminen
	wg.Wait()

	planificador.CerrarSemaforos()
}
