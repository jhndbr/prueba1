package entity

import (
	"time"
)

// TCB será la estructura base que utilizaremos dentro del Kernel
// para administrar los hilos de los diferentes procesos.
type TCB struct {
	TID               uint32 // Lista de los identificadores de los hilos asociados al proceso
	PID               uint32 // ID de su proceso padre
	Proceso           *PCB
	Priority          uint32        // Es la prioridad del hilo dentro del sistema.
	CPUTimeUsed       time.Duration // Tiempo total de CPU usado por el proceso
	CreationTime      time.Time     // Momento de creación del proceso
	LastExecutionTime time.Time     // Momento de la última ejecución
}

func CrearTCB(pcb *PCB, priority uint32) *TCB {
	var tid = uint32(pcb.TIDs.Size())
	pcb.TIDs.Add(tid)
	return &TCB{
		TID:               tid,
		PID:               pcb.PID,
		Proceso:           pcb,
		Priority:          priority,
		CPUTimeUsed:       0,
		CreationTime:      time.Now(),
		LastExecutionTime: time.Now(),
	}
}
