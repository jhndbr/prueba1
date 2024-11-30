package rest

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/usecase"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
)

func AtenderInterrupcion(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
			return
		}
	}(r.Body)

	var interrupcion *dto.Interrupt

	// Deserealizo el JSON a la estructura
	err = json.Unmarshal(body, &interrupcion)
	if err != nil {
		http.Error(w, "Error al deserializar el JSON", http.StatusBadRequest)
		slog.Error("Error al deserializar el JSON", "Error", err)
		return
	}

	// LOG OBLIGATORIO
	//Interrupción Recibida: “## Llega interrupcion al puerto Interrupt”.
	slog.Info("## Llega interrupcion al puerto Interrupt")

	go usecase.CheckInterrupt(interrupcion) // mando el TID a checkInterrupt para ser interrumpido

	// Responder con un mensaje de éxito
	// bloqueo hasta que se habilite para el response

	//<-entity.SEM_INTERRUPCION

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(rest.INTERRUPT_KERNEL)

	w.WriteHeader(http.StatusOK)
}
