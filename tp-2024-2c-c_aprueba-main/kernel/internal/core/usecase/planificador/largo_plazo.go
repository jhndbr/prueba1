package planificador

import (
	"errors"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/repository"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/rest/memory"
)

// AtenderProcesosNuevos rutina destina a crear los procesos en Memoria
func AtenderProcesosNuevos() {

	for {
		_ = semProcesosNew.Wait()                         // Toma un lugar
		nuevoProceso := repository.ObtenerProcesoEnNEW(0) // Utiliza FIFO

		err := memory.EnviarPCBHaciaMemoria(nuevoProceso)
		slog.Info("SE ENVIO EL PROCESO A MEMORIA", "proceso:", nuevoProceso)
		if err != "" {
			switch err {
			case "hay hueco":
				AtenderCompactacion()
				_ = semProcesosNew.Post()
				continue

			case "No hay espacio":
				slog.Info("No hay espacio en memoria; espero PROCESS_EXIT")
				//setea esto para que se le informe cuando terminan los procesos de aca en adelante
				largoPlazoDetenido = true
				value, _ := semProcesosExit.GetValue()
				slog.Debug("semaforo", "value:", value)
				_ = semProcesosExit.Wait()

				//termino un proceso
				largoPlazoDetenido = false
				_ = semProcesosNew.Post()
				//pasa al siguiente loop del for para intentar de nuevo
				continue
			default:
				slog.Error("Error al enviar proceso a memoria")
				continue
			}

		}

		repository.AddPrioridadEnOrden(nuevoProceso.Priority)

		//si el algoritmo es PRIORIDADES o CMN, setea la prioridad menor como la actual
		if nuevoProceso.Priority < prioridadActual && algoritmo != "FIFO" {
			prioridadActual = nuevoProceso.Priority

			//si hiloEnEjecucion es nil, no hay que interrumpir cpu
			if hiloEnEjecucion != nil {
				replanificar = true
				repository.AgregarHiloEnREADYSegunAlgoritmo(hiloEnEjecucion)
				InterrumpirCPU(hiloEnEjecucion)
			}

		}
		repository.AgregarProcesoEnREADY(nuevoProceso)
		_ = semReady.Post()

	}
}

func AtenderCompactacion() {
	slog.Debug("Detengo planificacion y pido compactacion")
	detenerPlanificacion = true

	_ = semPlanificacionDetenida.Wait()

	err := memory.SolicitarCompactacion()
	if err != nil {
		slog.Error("Error en la compactacion")
		return
	}
	detenerPlanificacion = false
	_ = semCompactacionFinalizada.Post()

}

// + Me parece que estos hay que borrar, estan obsoletos
// CrearProcesoNEW rutina destinada a crear las estructuras en la cola de NEW
func CrearProcesoNEW(fileName string, size uint32, priority uint32) error {
	pcb := entity.CrearPCB(fileName, size, priority)
	if pcb == nil {
		return errors.New("error en crear el proceso")
	}
	repository.AgregarProcesoEnNEW(pcb)
	slog.Info("Se crea el proceso - Estado: NEW -", "PID:", pcb.PID)
	_ = semProcesosNew.Post() // Unlock AtenderProcesosNew
	return nil
}
