package memory

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/repository"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
)

type Response struct {
	PID uint32 `json:"pid"`
}

func crearHilo(pcb *entity.PCB) (*entity.TCB, error) {

	var prioridad = pcb.Priority
	tcb := entity.CrearTCB(pcb, prioridad)
	repository.AddPrioridadEnOrden(prioridad)

	//+EnviarNuevoHilo(tcb,"")

	repository.AgregarHiloEnREADYSegunAlgoritmo(tcb)
	repository.AgregarAListaGlobal(tcb)
	return tcb, nil

}

func EnviarPCBHaciaMemoria(pcb *entity.PCB) string {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/process", conf.IPMemory, conf.PortMemory)

	request := dto.Process{
		PID:  pcb.PID,
		Size: pcb.Size,
	}

	body, _ := json.Marshal(request)
	resp, err := rest.Send(url, "POST", body)

	if err != nil || resp.StatusCode != http.StatusOK {
		//pasa el proceso a exit
		if resp.StatusCode == http.StatusBadRequest {
			return "No hay espacio"
		} else if resp.StatusCode == http.StatusNotFound {
			return "hay hueco"
		}
		return "Error desconocido"
	}

	//if err != nil {
	//	// TODO: Podria fallar porque requiere Compactacion o Porque no hay memoria suficiente
	//	slog.Error("Error al enviar el PID %s a Memoria. Error", "PID:", pcb.PID, "ERROR:", err.Error())
	//	return err
	//}

	var tcb *entity.TCB
	tcb, err = crearHilo(repository.PopProcesoEnNEW())
	if err != nil {
		return "Error al crear hilo"
	}
	EnviarNuevoHilo(tcb, pcb.FilePath)
	return ""
}

func EnviarNuevoHilo(tcb *entity.TCB, file string) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/hilo", conf.IPMemory, conf.PortMemory)

	request := dto.Process{
		PID:      tcb.PID,
		FilePath: file,
	}
	body, _ := json.Marshal(request)
	_, err := rest.Send(url, "POST", body)
	if err != nil {
		return
	}
	// callbackPCB(response, err)
}

func SolicitarFinalizacionHilo(tcb *entity.TCB) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/finalizar-hilo", conf.IPMemory, conf.PortMemory)
	body, _ := json.Marshal(tcb)
	_, err := rest.Send(url, "POST", body)
	if err != nil {
		slog.Error("Error al finalizar el hilo en memoria")
		return
	}

}

func SolicitarFinalizacionProceso(pcb *entity.PCB) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/finalizar-proceso", conf.IPMemory, conf.PortMemory)
	body, _ := json.Marshal(pcb)
	_, err := rest.Send(url, "POST", body)
	if err != nil {
		slog.Error("Error al finalizar el hilo en memoria")
		return
	}
}

func SolicitarCompactacion() error {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/compactarMemoria", conf.IPMemory, conf.PortMemory)
	body, _ := json.Marshal("compactar")
	_, err := rest.Send(url, "POST", body)
	if err != nil {
		slog.Error("Error al solicitar compactacion")
		return err
	}
	return nil
}
