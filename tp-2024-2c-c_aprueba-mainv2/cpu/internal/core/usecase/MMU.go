package usecase


import (
	"log/slog"

)

func MMU(offset uint32, base uint32, limite uint32) {
	slog.Info("Calculando direccion fisica","offset",offset,"base",base,"limite",limite)
	
	if offset > limite {
		huboSegmentation = true
		return
	}
	huboSegmentation = false

	DireccionFisica = base + offset
}
