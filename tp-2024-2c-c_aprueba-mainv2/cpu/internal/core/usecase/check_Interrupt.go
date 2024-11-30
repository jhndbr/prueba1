package usecase

import (
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"log/slog"
)

func CheckInterrupt(interrupt *dto.Interrupt) {
	if interrupt.TID == dto.CtxEnEjecucion.TID {

		dto.Interrumpir = true

	} else {
		slog.Info("CHECK_INTERRUPT: El TID a interrumpir no coincide con el TID en ejecución y no se detiene la ejecucion")
	}
}
