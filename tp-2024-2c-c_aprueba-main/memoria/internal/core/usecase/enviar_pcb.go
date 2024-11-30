package usecase

/*
import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"log/slog"
)

func EnviarProceso(pcb *entity.PCB) error {

	conf := config.GetInstance()
	urlCpu := fmt.Sprintf("http://%s:%d/cpu/execute", conf.IPCPU, conf.PortCPU)

	body, err := json.Marshal(pcb)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	status, err := rest.Send(urlCpu, body)

	if err != nil {
		slog.Error("No se pudo enviar PCB a CPU.")
		return err
	}

	slog.Info("Proceso enviado correctamente. Estado:", status)

	return nil
}

*/
