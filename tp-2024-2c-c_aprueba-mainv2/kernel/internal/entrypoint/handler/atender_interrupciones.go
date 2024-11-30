package handler

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/usecase/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/entrypoint/dto"
)

// AtenderInterrupciones atiende las interrupciones de CPU
func AtenderInterrupciones(w http.ResponseWriter, r *http.Request) {

	// Establecemos el tipo de contenido como "application/json"
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var interrupter *dto.Interrupt
		byteValue, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(byteValue, &interrupter)

		if err != nil {
			slog.Error("Error to deserializer Interrupter.")
			http.Error(w, "Error to validate Syscall.", http.StatusBadRequest)
			return
		}
		slog.Info("ERROR", "Interrupcion de CPU:", interrupter.Motivo)
		w.WriteHeader(http.StatusCreated)
		switch interrupter.Motivo {
		case "SEGMENTATION FAULT":
			planificador.AtenderProcessExit(interrupter.PID)
			planificador.SignalPlanificador()
		case "DESALOJO":
			slog.Info("Hilo correctamente desalojado de CPU")

			planificador.SignalInterrupt()

		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusBadRequest)
	}

}
