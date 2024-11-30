package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
)

type Process struct {
	PID  uint32 `json:"pid"`
	Size uint32 `json:"size"`
}

func CrearProcesoHandler(w http.ResponseWriter, r *http.Request) {
	var request Process
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	slog.Debug("CrearProcesoHandler", slog.String("PID", strconv.Itoa(int(request.PID))), slog.Int("Tamano", int(request.Size)))

	conf := config.GetInstance()
	estrategia := conf.Scheme

	if estrategia == "DINAMICAS" {

		if espacioLibreTotal() >= int(request.Size) {
			if !hayhuecoynoescontinuo(int(request.Size)) {
				http.Error(w, "hay hueco", http.StatusNotFound)
			} else {
				particionEstado, asignarError := usecase.AsignarParticion(uint32(request.PID), int(request.Size), estrategia)
				if asignarError != nil {
					http.Error(w, "Error al asignar partición", http.StatusBadRequest)
					return
				}
				slog.Info(fmt.Sprintf("## Proceso Creado - PID: %d - Tamaño: %d", request.PID, request.Size))
				if particionEstado != -1 {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
				}
			}
		} else {
			slog.Debug("No hay espacio suficiente para el proceso")
			http.Error(w, "No hay espacio", http.StatusBadRequest)
			return
		}
	} else {

		particionEstado, asignarError := usecase.AsignarParticion(uint32(request.PID), int(request.Size), estrategia)
		if asignarError != nil {
			http.Error(w, "No hay espacio", http.StatusBadRequest)
			return
		}
		slog.Info(fmt.Sprintf("## Proceso Creado - PID: %d - Tamaño: %d", request.PID, request.Size))
		if particionEstado != -1 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}
	}

}
func hayhuecoynoescontinuo(size int) bool {
	for _, p := range entity.Particiones {
		if p.Libre && (p.Limite-p.Base) >= size {
			return true
		}
	}
	return false
}
func espacioLibreTotal() int {
	total := 0
	for _, particion := range entity.Particiones {
		if particion.Libre {
			total += int(math.Abs(float64(particion.Base - particion.Limite)))
		}
	}
	return total
}

func avisarCompactacionAlKernel() {
	//TODO FALTA IMPLEMETAR ESTA LOGICA
	slog.Debug("Se notifica al Kernel para realizar compactación")
}

func FinalizarProcesoHandler(w http.ResponseWriter, r *http.Request) {
	// Estructura para parsear el body
	var requestBody struct {
		PID int `json:"pid"`
	}

	// Decodificar el cuerpo de la solicitud en JSON
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Cuerpo de la solicitud inválido", http.StatusBadRequest)
		return
	}
	pid := requestBody.PID
	slog.Debug("se realiza finalizar proceso: ", slog.Int("PID", pid))
	// Validar el PID
	if pid <= 0 {
		http.Error(w, "PID inválido", http.StatusBadRequest)
		return
	}

	var hilosAEliminar []uint32
	for i, contexto := range entity.MemoriaSistema {
		if entity.MemoriaSistema[i].PID == uint32(pid) {
			hilosAEliminar = append(hilosAEliminar, contexto.TID)
		}
	}

	for _, tid := range hilosAEliminar {
		usecase.EliminarHilo(uint32(pid), tid)
	}

	usecase.LiberarParticion(uint32(pid))

	//FALTA ELIMINAR ESTRUCTURAS DEL SISTEMA
	//slog.Info("Memoria compactada y particiones libres consolidadas.")
	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	//falta el tamnio log minimo y obligatorio
	slog.Info(fmt.Sprintf("## Proceso Destruido - PID: %d ", pid))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	slog.Debug("Proceso finalizar finalizado")
	return

	http.Error(w, "Proceso no encontrado", http.StatusNotFound)
}

func hayEspacio(tamano int) bool {
	for _, particion := range entity.Particiones {
		// Verifica si la partición está libre y si su tamaño es suficiente
		if particion.Libre && particion.Limite >= tamano {
			return true
		}
	}
	return false
}

//FALTA LOGICA PARA TERMINARLO
/*
func CrearProcesoHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error en el formato del request", http.StatusBadRequest)
		return
	}

	conf := config.GetInstance()
	base, err := AsignarParticion(req.PID, req.Size, conf.Scheme)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"base": base})
}

func EliminarProcesoHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error en el formato del request", http.StatusBadRequest)
		return
	}

	LiberarParticion(req.PID)
	w.WriteHeader(http.StatusOK)
}
*/
