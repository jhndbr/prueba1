package usecase

import (
	"errors"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
)

// lee 4 bytes a partir de una dirección específica

func LeerMemoria(direccion int) ([]byte, error) {
	entity.MemMutex.Lock()
	defer entity.MemMutex.Unlock()

	if direccion < 0 || direccion+4 > len(entity.MemoriaUsuario) {
		return nil, errors.New("Segmentation fault: Dirección fuera de los límites")
	}
	return entity.MemoriaUsuario[direccion : direccion+4], nil
}

// escribe un valor en una dirección específica
func EscribirMemoria(direccion int, valor []byte) error {
	//slog.Debug("MEMORIA DE USUARIO ANTES DE ESCRIBIR: ","memoriaUsuario:", entity.MemoriaUsuario)
	entity.MemMutex.Lock()
	slog.Debug("Escribiendo en memoria")
	defer entity.MemMutex.Unlock()

	if direccion < 0 || direccion+4 > len(entity.MemoriaUsuario) {
		slog.Debug("Segmentation fault: Dirección fuera de los límites")
		return errors.New("Segmentation fault: Dirección fuera de los límites")
	}
	if len(valor) != 4 {
		slog.Debug("Error: Se deben escribir exactamente 4 bytes")
		return errors.New("Error: Se deben escribir exactamente 4 bytes")
	}
	copy(entity.MemoriaUsuario[direccion:direccion+4], valor)
	//slog.Debug("MEMORIA DE USUARIO DESPUES DE ESCRIBIR: ","memoriaUsuario:", entity.MemoriaUsuario)
	return nil
}
