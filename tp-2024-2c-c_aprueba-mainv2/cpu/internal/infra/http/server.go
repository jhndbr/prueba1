package http

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/entrypoint/handler/rest"
	"log/slog"
	"net/http"
)

func Server(ip string, port int) error {

	// Asociamos las rutas con las funciones correspondientes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Establecemos el tipo de contenido como "application/json"
		w.Header().Set("Content-Type", "application/json")
		// Escribimos la respuesta JSON en el cuerpo de la respuesta
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/cpu/dispatcher", rest.AtenderDispatcher)
	http.HandleFunc("/cpu/interrupter", rest.AtenderInterrupcion)

	addr := fmt.Sprintf("%s:%d", ip, port)
	url := fmt.Sprintf("http://%s:%d", ip, port)

	slog.Info("Starting Up!", "url", url)

	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("Cannot start server. Error=", err.Error())
		return err
	}

	return nil
}
