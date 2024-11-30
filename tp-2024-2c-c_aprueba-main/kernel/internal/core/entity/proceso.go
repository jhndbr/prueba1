package entity

import "github.com/sisoputnfrba/tp-golang/utils/list"

// P es la interfaz a utilizar para el manejo de procesos.
type P interface {
	New(filePath string, size uint32, priority uint32) *PCB
}

// PCB será la estructura base que utilizaremos dentro del Kernel para administrar los procesos.
// El mismo deberá contener como mínimo los datos definidos a continuación,
// que representan la información administrativa.
type PCB struct {
	PID      uint32                 // Identificador único del proceso
	TIDs     list.ArrayList[uint32] // Lista de los identificadores de los hilos asociados al proceso
	MUTEXs   list.ArrayList[Mutex]  // Mutex list
	Priority uint32                 // Priority of Main Thread
	Size     uint32                 // Size of Process
	FilePath string                 // FilePath of the instructions
}

// Global Variable
var pid uint32 = 0

// New Crea la estructura
func (p PCB) Crear(filePath string, size uint32, priority uint32) *PCB {
	return &PCB{
		PID:      p.PID,
		TIDs:     list.ArrayList[uint32]{},
		MUTEXs:   list.ArrayList[Mutex]{},
		Priority: priority,
		FilePath: filePath,
	}
}

func getNextId() uint32 {
	pid = pid + 1
	return pid
}

// CrearPCB crea la estructura
func CrearPCB(filePath string, size uint32, priority uint32) *PCB {
	return &PCB{
		PID:      getNextId(),
		TIDs:     list.ArrayList[uint32]{},
		MUTEXs:   list.ArrayList[Mutex]{},
		Size:     size,
		Priority: priority,
		FilePath: filePath,
	}
}
