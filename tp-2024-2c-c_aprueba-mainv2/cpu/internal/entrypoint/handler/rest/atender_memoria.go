package rest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type ValorMemoria struct {
	Valor uint32
}

func AtenderMemoria(w http.ResponseWriter, r *http.Request) {

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

	var valorMemoria ValorMemoria

	// Deserealizo el JSON a la estructura
	err = json.Unmarshal(body, &valorMemoria)
	if err != nil {
		http.Error(w, "Error al deserializar el JSON", http.StatusBadRequest)
		slog.Error("Error al deserializar el JSON", "Error", err)
		return
	}

	slog.Info("Se recibió el valor de memoria correspsondiente a la direccion fisica ", "valor:", valorMemoria.Valor)

	//usecase. = &valorMemoria.Valor // guardo el valor de memoria en variable global

	// DESBLOQUEAR SEMAFORO DE READMEM
	//usecase.SemValorRecibido <- "RECIBI EL VALOR DE MEMORIA"

	w.WriteHeader(http.StatusOK)
}
