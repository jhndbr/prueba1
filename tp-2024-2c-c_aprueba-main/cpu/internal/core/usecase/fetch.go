package usecase

import (
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/kernel"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/memoria"
	"log/slog"
)

// Fetch va a buscar a Memoria el siguiente IR
func Fetch(context *dto.Context) {

	slog.Info("##", "TID", context.TID, "FETCH - Program counter:", context.PC)

	instruction, err := memoria.EnviarSolicitudInstruccion(context)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	if instruction == nil {
		err := memoria.ActualizarContextoDeEjecucion(context)
		if err != nil {
			return
		}
		kernel.EnviarMotivoInterrupcion(context.PID, context.TID, "EXIT")
		return
	}

	// Referencia a la próxima Instrucción a Ejecutar
	context.IR = instruction
	Decode(context)

}
