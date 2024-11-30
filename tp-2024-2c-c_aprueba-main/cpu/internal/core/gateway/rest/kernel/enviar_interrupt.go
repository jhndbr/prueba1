package kernel

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
)

// EnviarMotivoInterrupcion al Kernel
func EnviarMotivoInterrupcion(PID uint32, TID uint32, motivo string) {

	request := dto.Interrupt{
		PID:    PID,
		TID:    TID,
		Motivo: motivo,
	}
	body, _ := json.Marshal(request)
	_, err := rest.Send(urlInterrupt, "POST", body)
	if err != nil {
		return
	}
}
