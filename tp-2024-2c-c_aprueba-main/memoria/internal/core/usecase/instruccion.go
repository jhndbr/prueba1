package usecase

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
)

func ObtenerInstruccion(pid uint32, tid uint32, pc int) (entity.Instruction, error) {

	key := fmt.Sprintf("%d-%d", pid, tid)

	instruccionesHilo, existe := instruccionesPorHilo[key]
	if !existe {
		return entity.Instruction{}, errors.New("instrucciones no encontradas para el TID especificado")
	}

	// controlamos que el pc este en el rango de instrucciones
	if pc < 0 || pc >= len(instruccionesHilo.Instrucciones) {
		return entity.Instruction{}, errors.New("PC fuera de los límites de las instrucciones")
	}

	//LE PASAMOS DE A UNA A CPU LAS INTRUCCIONES
	return instruccionesHilo.Instrucciones[pc], nil
}
func CargarListaArchivo(pid uint32, filepath string) {
	archivo := entity.Archivo{
		PID:      pid,
		FilePath: filepath,
	}
	entity.ListaArchivos = append(entity.ListaArchivos, archivo)
	slog.Debug("Archivo registrado", "PID", pid, "Filepath", filepath)
}

func BuscarArchivoPorPID(pid uint32) (string, error) {
	for _, archivo := range entity.ListaArchivos {
		if archivo.PID == pid {
			return archivo.FilePath, nil
		}
	}
	return "", errors.New("archivo no encontrado para el PID especificado")
}

func EliminarArchivo(pid uint32) error {
	for i, archivo := range entity.ListaArchivos {
		if archivo.PID == pid {
			entity.ListaArchivos = append(entity.ListaArchivos[:i], entity.ListaArchivos[i+1:]...)
			slog.Debug("Archivo eliminado", "PID", pid)
			return nil
		}
	}
	return errors.New("archivo no encontrado para el PID especificado")
}
