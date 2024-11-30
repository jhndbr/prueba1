package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/usecase/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/entrypoint/dto"
)

func validarSyscall(syscall dto.Syscall) error {

	switch syscall.Code {
	case "PROCESS_CREATE":
		prioridad := syscall.PrioridadHiloMain
		slog.Info("", "prioridad", prioridad)
		err := planificador.AtenderProcessCreate(syscall.FileName, syscall.ProcessSize, syscall.PrioridadHiloMain)
		if err != nil {
			return err
		}

	case "PROCESS_EXIT":
		planificador.AtenderProcessExit(syscall.PID)

	case "THREAD_CREATE":
		planificador.AtenderThreadCreate(syscall.Priority, syscall.FileName)

	case "THREAD_JOIN":
		planificador.AtenderThreadJoin(syscall.TID)

	case "THREAD_CANCEL":
		planificador.AtenderThreadCancel(syscall.TID)

	case "THREAD_EXIT":
		planificador.AtenderThreadExit(syscall.PID, syscall.TID)

	case "MUTEX_CREATE":
		planificador.AtenderMutexCreate(syscall.MutexID)

	case "MUTEX_LOCK":
		planificador.AtenderMutexLock(syscall.MutexID)

	case "MUTEX_UNLOCK":
		planificador.AtenderMutexUnlock(syscall.MutexID)

	case "DUMP":
		planificador.AtenderDump(syscall.PID, syscall.TID)

	case "IO":
		planificador.AtenderIO(syscall.Tiempo)
	default:
		return fmt.Errorf("invalid syscall code: %s", syscall.Code)
	}

	//planificador.LockMutexSyscall()
	if planificador.GetInterrupted() == true {
		planificador.SignalInterrupt()

	} else {
		planificador.SignalPlanificador()
	}
	//planificador.UnlockMutexSyscall()

	return nil
}

// AtenderSyscall atiende las syscall de CPU
func AtenderSyscall(w http.ResponseWriter, r *http.Request) {

	// Establecemos el tipo de contenido como "application/json"
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var syscall dto.Syscall
		byteValue, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(byteValue, &syscall)
		//slog.Info("", "syscall:", syscall)
		if err != nil {
			slog.Error("Error to deserializer Syscall.")
			http.Error(w, "Error to validate Syscall.", http.StatusBadRequest)
			return
		}
		errSyscall := validarSyscall(syscall)
		if errSyscall != nil {
			slog.Error("Error to validate Syscall.")
			http.Error(w, "Error to validate Syscall.", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusBadRequest)
	}

}
