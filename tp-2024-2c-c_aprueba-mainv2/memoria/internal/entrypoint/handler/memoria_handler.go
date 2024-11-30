package handler

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
)

type WriteMemRequest struct {
	Direccion int    `json:"direction"`
	Valor     []byte `json:"value"`
	Pid       uint32 `json:"pid"`
	Tid       uint32 `json:"tid"`
}

func VerValorMemoria(valor []byte) (uint32, error) {
	if len(valor) != 4 {
		return 0, fmt.Errorf("longitud de valor incorrecta: se esperaban 4 bytes, se recibieron %d", len(valor))
	}
	valorUint32 := binary.BigEndian.Uint32(valor)
	return valorUint32, nil
}
func LeerMemoriaHandler(w http.ResponseWriter, r *http.Request) {
	//direccionStr := r.URL.Query().Get("direccion")
	//direccion, err := strconv.Atoi(direccionStr)
	//if err != nil {
	//	http.Error(w, "Dirección inválida", http.StatusBadRequest)
	//	return
	//}
	// Obtener la dirección de la memoria
	direccionStr := r.URL.Query().Get("direccion")
	direccion, err := strconv.Atoi(direccionStr)
	if err != nil {
		http.Error(w, "Dirección inválida", http.StatusBadRequest)
		return
	}

	// Obtener el PID
	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "PID inválido", http.StatusBadRequest)
		return
	}

	// Obtener el TID
	tidStr := r.URL.Query().Get("tid")
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		http.Error(w, "TID inválido", http.StatusBadRequest)
		return
	}

	valor, err := usecase.LeerMemoria(direccion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	valorEnUint32, err := VerValorMemoria(valor)
	slog.Debug("Valor leido de memoria", "direccion", direccion, "valor", valorEnUint32)
	slog.Info("## Lectura (PID:TID)",
		"PID", pid,
		"TID", tid,
		"direccion", direccion,
		"valor", valorEnUint32)
	w.Header().Set("Content-Type", "application/json")
	// Arregle como se envia el parametro valor, no se si esta bien
	err = json.NewEncoder(w).Encode(map[string][]byte{"valor": valor})
	if err != nil {
		slog.Error("Error", "error", err.Error())
		return
	}
}

func EscribirMemoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method has not allowed.", http.StatusMethodNotAllowed)
		return
	}
	slog.Debug("recibo Solicitud Write_MEM")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request is Empty or has an erro", http.StatusInternalServerError)
		return
	}
	// var ctx CtxCpu
	var req WriteMemRequest
	err = json.Unmarshal(body, &req)
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	http.Error(w, "Error en el formato del request", http.StatusBadRequest)
	// 	return
	// }
	slog.Debug("Escribiendo en memoria", "direccion", req.Direccion, "valor", req.Valor)
	slog.Info("## Leer (PID:TID)",
		"PID", req.Pid,
		"TID", req.Tid,
		"direccion", req.Direccion,
		"valor", req.Valor)
	// Arregle como se envia el parametro reqvalor, no se si esta bien
	err = usecase.EscribirMemoria(req.Direccion, req.Valor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
