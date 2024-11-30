package usecase

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
)

var instruccionesPorHilo = make(map[string]entity.InstruccionesPorHilo)

func CrearHilo(pid uint32, filepath string) error {

	// TODO: Cuando estamos creando un nuevo HILO va a buscarlo, y no lo encuentra.
	// Entonces siempre pisa el contexto en ejection.
	slog.Debug("Creando hilo", "PID:", pid, "Filepath:", filepath)
	tid := obtenerSiguienteTID(pid) // Revisar esta funcion.
	nuevoContexto := entity.ContextoEjecucion{
		PID: pid,
		TID: tid,
		AX:  0, BX: 0, CX: 0, DX: 0, EX: 0, FX: 0, GX: 0, HX: 0, PC: 0,
		Base:   uint32(obtenerBaseParticion(pid)),
		Limite: uint32(obtenerLimiteParticion(pid)),
	}
	slog.Debug("Nuevo contexto", "Contexto:", nuevoContexto)
	if err := cargarInstruccionesEnMemoria(pid, tid, filepath); err != nil {
		return err
	}

	err := GuardarContextoEnMemoria(nuevoContexto)
	slog.Info(fmt.Sprintf("## Hilo Creado - PID: %d - TID: %d", nuevoContexto.PID, nuevoContexto.TID))
	if err != nil {
		return err
	}
	slog.Info("Se crea el hilo - Estado: READY", "PID:", nuevoContexto)
	return nil
}

func EliminarHilo(pid uint32, tid uint32) error {

	entity.MemMutex.Lock()
	defer entity.MemMutex.Unlock()

	indexToDelete := -1
	for i, contexto := range entity.MemoriaSistema {
		if contexto.PID == pid && contexto.TID == tid {
			indexToDelete = i
			break
		}
	}
	if indexToDelete != -1 {
		entity.MemoriaSistema = append(entity.MemoriaSistema[:indexToDelete], entity.MemoriaSistema[indexToDelete+1:]...)
	} else {
		return errors.New("Hilo no encontrado en MemoriaSistema")
	}

	// eliminamos las intrucciones asociadas al hilo
	key := fmt.Sprintf("%d-%d", pid, tid)
	slog.Info(fmt.Sprintf("## Hilo Destruido - PID: %d - TID: %d", pid, tid))
	if _, exists := instruccionesPorHilo[key]; exists {
		delete(instruccionesPorHilo, key)
	} else {
		return errors.New("Instrucciones del hilo no encontradas")
	}

	return nil
}

func cargarInstruccionesEnMemoria(pid, tid uint32, pathArchivo string) error {
	slog.Debug("Cargando instrucciones en memoria", "PID:", pid, "TID:", tid, "Path:", pathArchivo)
	config := config.GetInstance()
	var pathFile = config.InstructionPath + pathArchivo
	archivo, err := os.Open(pathFile)
	if err != nil {
		return errors.New("No se pudo abrir el archivo de pseudocódigo")
	}
	defer archivo.Close()

	var instrucciones []entity.Instruction
	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() {
		partes := strings.Fields(scanner.Text())
		if len(partes) == 0 {
			continue
		}
		instrucciones = append(instrucciones, entity.Instruction{
			Code: partes[0],
			Args: partes[1:],
		})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if instruccionesPorHilo == nil {
		instruccionesPorHilo = make(map[string]entity.InstruccionesPorHilo)
	}

	// aca guardamos las instrucciones en el mapa con la clave de PID-TID
	key := fmt.Sprintf("%d-%d", pid, tid)
	instruccionesPorHilo[key] = entity.InstruccionesPorHilo{Instrucciones: instrucciones}
	slog.Debug("Instrucciones cargadas en memoria", "PID:", pid, "TID:", tid, "Instrucciones:", instrucciones)
	return nil

}

// guardamos el contexto en memoria de sistema
func GuardarContextoEnMemoria(contexto entity.ContextoEjecucion) error {
	//entity.MemMutex.Lock()
	//entity.MemMutex.Unlock()
	entity.MemoriaSistema = append(entity.MemoriaSistema, contexto)
	return nil
}

func obtenerBaseParticion(pid uint32) int {
	for _, particion := range entity.Particiones {
		if particion.PID == pid && !particion.Libre {
			return particion.Base
		}
	}
	return 0
}

func obtenerLimiteParticion(pid uint32) int {
	for _, particion := range entity.Particiones {
		if particion.PID == pid && !particion.Libre {
			return particion.Limite
		}
	}
	slog.Debug("No se encontro la particion", "PID:", pid)
	return 0
}
func obtenerSiguienteTID(pid uint32) uint32 {

	tid := uint32(0)

	for _, contexto := range entity.MemoriaSistema {
		if contexto.PID == pid {

			if contexto.TID >= tid {
				tid = contexto.TID + 1
			}
		}
	}
	return tid
}
