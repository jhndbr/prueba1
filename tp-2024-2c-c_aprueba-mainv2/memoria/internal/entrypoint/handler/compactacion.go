package handler

import (
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
)

func CompactarMemoriaHandler(w http.ResponseWriter, r *http.Request) {

	// Verificar que el mensaje es "compactar"
	slog.Debug("CompactarMemoriaHandler")
	usecase.Compactar()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
