package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/entrypoint/handler"
)

func Server(ip string, port int) error {

	// Asociamos las rutas con las funciones correspondientes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Establecemos el tipo de contenido como "application/json"
		w.Header().Set("Content-Type", "application/json")

		// Escribimos la respuesta JSON en el cuerpo de la respuesta
		w.WriteHeader(http.StatusOK)
	})

	// Asociamos ruta
	http.HandleFunc("/kernel/syscall", handler.AtenderSyscall)
	http.HandleFunc("/kernel/interrupt", handler.AtenderInterrupciones)

	addr := fmt.Sprintf("%s:%d", ip, port)
	url := fmt.Sprintf("http://%s:%d", ip, port)

	slog.Info("Starting Up!", "url", url)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
