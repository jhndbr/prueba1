package usecase

import (
	"errors"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
)


func GetContexto(pid uint32, tid uint32) (*entity.ContextoEjecucion, error) {
	for _, contexto := range entity.MemoriaSistema {
		if contexto.PID == pid && contexto.TID == tid {
			return &contexto, nil
		}
	}
	return nil, errors.New("Contexto no encontrado")
}

func UpdateContexto(nuevoContexto entity.ContextoEjecucion) error {
	for i, contexto := range entity.MemoriaSistema {
		if contexto.PID == nuevoContexto.PID && contexto.TID == nuevoContexto.TID {
			entity.MemoriaSistema[i] = nuevoContexto
			return nil
		}
	}

	return errors.New("Contexto no encontrado para actualizar")
}
func CrearContexto(pid uint32, tid uint32) {
	contexto := entity.ContextoEjecucion{
		PID:    pid,
		TID:    tid,
		AX:     0,
		BX:     0,
		CX:     0,
		DX:     0,
		EX:     0,
		FX:     0,
		GX:     0,
		HX:     0,
		PC:     0,
		Base:   0,
		Limite: 0,
	}
	entity.MemoriaSistema = append(entity.MemoriaSistema, contexto)

}
