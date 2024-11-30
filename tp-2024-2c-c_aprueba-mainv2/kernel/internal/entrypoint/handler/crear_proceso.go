package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/repository"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/usecase/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/entrypoint/dto"
)

func CrearProceso(w http.ResponseWriter, r *http.Request) {

	// Establecemos el tipo de contenido como "application/json"
	w.Header().Set("Content-Type", "application/json")

	go func() {
		err := planificador.CrearProcesoNEW("instructions.txt", 512, 0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}()
	w.WriteHeader(http.StatusCreated)
}

func ListarProcesosHandle(w http.ResponseWriter, r *http.Request) {

	// Creamos una instancia de la estructura Status
	procesos := dto.ProcesosDTO{
		List: repository.ObtenerProcesosNEW(),
	}

	// Establecemos el tipo de contenido como "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Convertimos la estructura a JSON
	jsonResponse, err := json.Marshal(procesos)
	if err != nil {
		// Si hay un error al convertir a JSON, devolvemos un error interno del servidor
		http.Error(w, "Error al generar la respuesta JSON", http.StatusInternalServerError)
		return
	}

	// Escribimos la respuesta JSON en el cuerpo de la respuesta
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Error al reforzar la respuesta JSON", http.StatusInternalServerError)
		return
	}

}
