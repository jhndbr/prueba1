package repository

import (
	"errors"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/utils/list"
)

var mutexs list.ArrayList[*entity.Mutex]

func ObtenerMutex(mutexID string, proceso *entity.PCB) (*entity.Mutex, error) {
	for i := 0; i < mutexs.Size(); i++ {
		var mi, _ = mutexs.Get(i)
		if mi.Nombre == mutexID && proceso.PID == mi.PID {
			return mi, nil
		}
	}
	return nil, errors.New("mutex not found")
}

func LockMutex(mutexID string, proceso *entity.PCB) {
	mutex, err := ObtenerMutex(mutexID, proceso)
	if err != nil {
		slog.Error("Mutex no encontrado")
	}
	mutex.Estado = 1
}

func UnlockMutex(mutexID string, proceso *entity.PCB) {
	mutex, err := ObtenerMutex(mutexID, proceso)
	if err != nil {
		slog.Error("Mutex no encontrado")
	}
	mutex.Estado = 0
}

func AddMutex(mutex *entity.Mutex) {
	mutexs.Add(mutex)
}
