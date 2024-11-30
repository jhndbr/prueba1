package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
)

// ObtenerInstruccionHandler maneja la solicitud para obtener una instrucción
func ObtenerInstruccionHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.GetInstance()
	time.Sleep(time.Duration(conf.ResponseDelay) * time.Millisecond) // Retardo de peticiones

	// Obtenemos el PID, TID y PC de los parámetros de la URL
	pidStr := r.URL.Query().Get("pid")
	tidStr := r.URL.Query().Get("tid")
	pcStr := r.URL.Query().Get("pc")

	// Validar PID
	pid64, err := strconv.ParseUint(pidStr, 10, 32) // Convierte a uint64 temporalmente
	if err != nil || pid64 > uint64(^uint32(0)) {   // Asegura que está dentro del rango de uint32
		http.Error(w, "PID inválido", http.StatusBadRequest)
		return
	}
	pid := uint32(pid64)

	// Validar TID
	tid64, err := strconv.ParseUint(tidStr, 10, 32) // Convierte a uint64 temporalmente
	if err != nil || tid64 > uint64(^uint32(0)) {   // Asegura que está dentro del rango de uint32
		http.Error(w, "TID inválido", http.StatusBadRequest)
		return
	}
	tid := uint32(tid64)

	// Validar PC
	pc, err := strconv.Atoi(pcStr)
	if err != nil || pc < 0 {
		http.Error(w, "PC inválido", http.StatusBadRequest)
		return
	}

	// Obtener la instrucción
	instruccion, err := usecase.ObtenerInstruccion(pid, tid, pc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf("## Obtener instrucción - PID: %d - TID: %d", pid, tid))

	// Configurar el encabezado de respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instruccion) // Envía la instrucción en formato JSON
}
