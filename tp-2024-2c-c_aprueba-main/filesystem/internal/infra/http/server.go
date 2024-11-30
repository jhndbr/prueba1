package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/filesystem/internal/entrypoint/handler"
)

// func handle() func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodGet:
// 		case http.MethodPost:
// 			handler.Crear_archivo(w, r)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	}
// }

func Server(ip string, port int) error {
	// Asociamos las rutas con las funciones correspondientes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	// Asociamos ruta
	http.HandleFunc("/filesystem/crear", handler.Crear_archivo)

	addr := fmt.Sprintf("%s:%d", ip, port)
	url := fmt.Sprintf("http://%s:%d", ip, port)

	slog.Info("Starting Up!", "url", url)
	fmt.Println()
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
