package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/entrypoint/dto"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
)

func DecodeContext(ctx *CtxCpu) entity.ContextoEjecucion {

	return entity.ContextoEjecucion{
		PID:    ctx.PID,
		TID:    ctx.TID,
		AX:     ctx.Register.AX,
		BX:     ctx.Register.BX,
		CX:     ctx.Register.CX,
		DX:     ctx.Register.DX,
		EX:     ctx.Register.EX,
		FX:     ctx.Register.FX,
		GX:     ctx.Register.GX,
		HX:     ctx.Register.HX,
		PC:     int(ctx.PC),
		Base:   ctx.Base,
		Limite: ctx.Limit,
	}

}

// Obtener contexto de ejecución (desde CPU)
func ObtenerContexto(w http.ResponseWriter, r *http.Request) {
	// Obtener PID y TID de los parámetros de la URL
	pidStr := r.URL.Query().Get("pid")
	tidStr := r.URL.Query().Get("tid")

	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid < 0 {
		http.Error(w, "PID inválido", http.StatusBadRequest)
		return
	}
	tidStr = strings.TrimSpace(tidStr)
	tidStr = strings.ReplaceAll(tidStr, "\n", "")

	tid, err := strconv.Atoi(tidStr)
	if err != nil || tid < 0 {
		http.Error(w, "TID inválido", http.StatusBadRequest)
		return
	}

	// Convertir pid y tid a uint32
	pidUint32 := uint32(pid)
	tidUint32 := uint32(tid)

	contexto, err := usecase.GetContexto(pidUint32, tidUint32)
	if err != nil {
		http.Error(w, "Contexto no encontrado", http.StatusNotFound)
		return
	}
	slog.Debug("Contexto encontrado", "contexto", contexto)
	slog.Info(fmt.Sprintf("## Contexto Solicitado - PID: %d - TID: %d", pid, tid))
	// envio en formato JSON
	w.Header().Set("Content-Type", "application/json")

	ctx := dto.NewContext(contexto)
	slog.Info("ctx a ENVIAR", "ctx", ctx)
	err = json.NewEncoder(w).Encode(ctx)
	if err != nil {
		slog.Error("Error=", err.Error())
		return
	}
}

type CtxCpu struct {
	PID      uint32       `json:"PID"`
	TID      uint32       `json:"tid"`
	Register dto.Register `json:"register"`
	PC       uint32       `json:"pc"`
	Base     uint32       `json:"base"`
	Limit    uint32       `json:"limit"`
}

func ActualizarContexto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method has not allowed.", http.StatusMethodNotAllowed)
		return
	}
	slog.Debug("recibo Actualizando contexto")
	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request is Empty or has an erro", http.StatusInternalServerError)
		return
	}
	var ctx CtxCpu
	err = json.Unmarshal(body, &ctx)
	nuevoContexto := DecodeContext(&ctx)

	// actualizo
	err = usecase.UpdateContexto(nuevoContexto)
	if err != nil {
		http.Error(w, "Error actualizando el contexto", http.StatusInternalServerError)
		return
	}
	slog.Debug("Contexto actualizado correctamente")
	slog.Info(fmt.Sprintf("## Contexto Actualizado - PID: %d - TID: %d", nuevoContexto.PID, nuevoContexto.TID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"OK"}`))
	//pidStr := r.URL.Query().Get("pid")
	//tidStr := r.URL.Query().Get("tid")
	//
	//pid, err := strconv.Atoi(pidStr)
	//if err != nil || pid < 0 {
	//	http.Error(w, "PID inválido", http.StatusBadRequest)
	//	return
	//}
	//tidStr = strings.TrimSpace(tidStr)
	//tidStr = strings.ReplaceAll(tidStr, "\n", "")
	//tid, err := strconv.Atoi(tidStr)
	//if err != nil || tid < 0 {
	//	http.Error(w, "TID inválido", http.StatusBadRequest)
	//	return
	//}
	//
	//// Log para depuración
	//slog.Debug("PID:", pid)
	//slog.Debug("TID:", tid)
	// Decodificar el contexto nuevo desde el cuerpo de la solicitud
	//err = json.NewDecoder(r.Body).Decode(&nuevoContexto)
	//if err != nil {
	//	http.Error(w, "Error procesando el contexto", http.StatusBadRequest)
	//	return
	//}
}
