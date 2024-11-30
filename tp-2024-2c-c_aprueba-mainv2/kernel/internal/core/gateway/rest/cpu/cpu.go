package cpu

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
)

type Request struct {
	PID uint32 `json:"pid"`
	TID uint32 `json:"tid"`
}

// EnviarIDsACPU falta terminar, queria tener la forma aproximada nomas
func EnviarIDsACPU(tcb *entity.TCB) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/cpu/dispatcher", conf.IPCPU, conf.PortCPU)
	request := Request{
		PID: tcb.PID,
		TID: tcb.TID,
	}
	body, _ := json.Marshal(request)
	response, err := rest.Send(url, "POST", body)
	if err != nil {
		slog.Error("Error=", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		slog.Error("", "StatusCode=", response.StatusCode)
	}
	slog.Debug("HILO ENVIADO A CPU", "PID:", tcb.PID, "TID:", tcb.TID)
}

func EnviarCPUInterrupt(tcb *entity.TCB) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/cpu/interrupter", conf.IPCPU, conf.PortCPU)
	request := Request{
		PID: tcb.PID,
		TID: tcb.TID,
	}
	body, _ := json.Marshal(request)
	response, err := rest.Send(url, "POST", body)
	if err != nil {
		slog.Error("Error=", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		slog.Error("", "StatusCode=", response.StatusCode)
	}
	slog.Debug("Solicitada interrupcion de CPU")

}
