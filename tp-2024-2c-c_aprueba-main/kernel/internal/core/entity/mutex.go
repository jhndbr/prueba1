package entity

import (
	"github.com/sisoputnfrba/tp-golang/utils/list"
)

type MutexState int

// Definición de los estados del proceso utilizando iota
const (
	UNLOCKED MutexState = iota
	LOCKED
)

type Mutex struct {
	Nombre          string
	Estado          MutexState
	PID             uint32
	HilosBloqueados list.ArrayList[*TCB]
}

func CrearMutex(mutexID string, pid uint32) *Mutex {
	return &Mutex{
		Nombre:          mutexID,
		Estado:          UNLOCKED,
		PID:             pid,
		HilosBloqueados: list.ArrayList[*TCB]{},
	}
}
