package handler

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
)

type MemoryDumpRequest struct {
	PID uint32 `json:"pid"`
	TID uint32 `json:"tid"`
}

func MemoryDump(w http.ResponseWriter, r *http.Request) {

	slog.Debug("Recibida solicitud de dump de memoria")

	if r.Method != http.MethodPost {
		http.Error(w, "Method has not allowed.", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request is Empty or has an erro", http.StatusInternalServerError)
		return
	}
	// var context *dto.Context
	var request MemoryDumpRequest

	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Request has not been deserialize", http.StatusBadRequest)
		return
	}

	// Obtener las particiones del proceso
	var memoryDump []byte
	for _, particion := range entity.Particiones {
		slog.Debug("Partición", "Base", particion.Base, "Limite", particion.Limite, "Libre", particion.Libre, "PID", particion.PID)
		if particion.PID == request.PID {
			if config.GetInstance().Scheme == "FIJAS" {
				start := particion.Base
				end := particion.Base + particion.Limite
				memoryDump = append(memoryDump, entity.MemoriaUsuario[start:end]...)
			} else {
				start := particion.Base
				end := particion.Limite
				memoryDump = append(memoryDump, entity.MemoriaUsuario[start:end]...)

			}
		}

	}
	// Crear el objeto que se enviará a Filesystem
	dumpData := struct {
		PID        uint32 `json:"pid"`
		TID        uint32 `json:"tid"`
		Size       int    `json:"size"`
		MemoryDump []byte `json:"contenido"`
	}{
		PID:        request.PID,
		TID:        request.TID,
		Size:       len(memoryDump), // verificar si es correcto o es el tamanio que pasa kernel
		MemoryDump: memoryDump,
	}
	slog.Debug("Envio a Filesystem", "dumpData", dumpData)
	// Codificar el dumpData en JSON
	dumpDataJSON, err := json.Marshal(dumpData)
	if err != nil {
		http.Error(w, "Error al codificar el dump de memoria", http.StatusInternalServerError)
		return
	}
	slog.Debug("Dump de memoria codificado en JSON envio a filesistem")
	// Enviar el dumpData a Filesystem
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/filesystem/crear", conf.IPFilesystem, conf.PortFilesystem)
	resp, err := rest.Send(url, "POST", dumpDataJSON)

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(dumpDataJSON))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Error al enviar el dump de memoria a Filesystem", http.StatusInternalServerError)
		slog.Error("Error al enviar el dump de memoria a Filesystem", "Error", err)
		return
	}
	slog.Debug("Dump de memoria enviado a Filesystem")
	slog.Info(fmt.Sprintf("## Memory Dump solicitado - PID: %d - TID: %d", dumpData.PID, dumpData.TID))
	// Responder OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(`{"status":"OK"}`))
}
