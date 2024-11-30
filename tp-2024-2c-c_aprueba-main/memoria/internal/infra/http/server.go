package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/entrypoint/handler"
)

// Función para manejar las solicitudes a la ruta de contexto de ejecución
// func handle() func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodGet:
// 			if r.URL.Path == "/memory/context" {
// 				handler.ObtenerContexto(w, r)
// 			}else if(r.URL.Path == "/memory/leerMemoria"){
// 				handler.LeerMemoriaHandler(w, r)
// 			}
// 			// handler.ObtenerContexto(w, r)
// 			// handler.LeerMemoriaHandler(w, r)
// 		case http.MethodPost:
// 			if r.URL.Path == "/memory/instruccion" {
// 				handler.ObtenerInstruccionHandler(w, r)
// 			} else if r.URL.Path == "/memory/hilo" {
// 				handler.CrearHiloHandler(w, r)
// 			} else if r.URL.Path == "/memory/finalizar-hilo" {
// 				handler.FinalizarHiloHandler(w, r)
// 			} else if r.URL.Path == "/memory/escribirMemoria" {
// 				handler.EscribirMemoriaHandler(w, r)
// 			} else {
// 				handler.ActualizarContexto(w, r)
// 			}
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	}
// }

func Server(ip string, port int) error {
	// Establecemos el tipo de contenido como "application/json" en la raíz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	// Configurar rutas
	http.HandleFunc("/memory/process", handler.CrearProcesoHandler)
	http.HandleFunc("/memory/finalizar-proceso", handler.FinalizarProcesoHandler)
	http.HandleFunc("/memory/compactarMemoria", handler.CompactarMemoriaHandler)
	http.HandleFunc("/memory/hilo", handler.CrearHiloHandler)
	http.HandleFunc("/memory/finalizar-hilo", handler.FinalizarHiloHandler)
	http.HandleFunc("/memory/instruction", handler.ObtenerInstruccionHandler)
	http.HandleFunc("/memory/context", handler.ObtenerContexto)
	http.HandleFunc("/memory/actualizarContexto", handler.ActualizarContexto)
	http.HandleFunc("/memory/escribirMemoria", handler.EscribirMemoriaHandler)
	http.HandleFunc("/memory/leerMemoria", handler.LeerMemoriaHandler)
	http.HandleFunc("/memory/syscall-memoryDump", handler.MemoryDump)

	addr := fmt.Sprintf("%s:%d", ip, port)
	url := fmt.Sprintf("http://%s:%d", ip, port)
	slog.Info("Starting Up!", "url", url)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
