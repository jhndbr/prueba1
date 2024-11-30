package kernel

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"log/slog"
)

// EnviarSolicitudCrearProceso Solicitar a Kernel crear Proceso.
func EnviarSolicitudCrearProceso(fileName string, tamanioProceso uint32, prioridad uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.ProcessCreate{
		Code:              "PROCESS_CREATE",
		FileName:          fileName,
		ProcessSize:       tamanioProceso,
		PrioridadHiloMain: prioridad,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudFinalizarProceso(pid uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.ProcessExit{
		Code: "PROCESS_EXIT",
		PID:  pid,
	}
	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudCrearHilo(fileName string, prioridad uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.ThreadCreate{
		Code:     "THREAD_CREATE",
		FileName: fileName,
		Priority: prioridad,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudThreadJoin(tid uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	body, _ := json.Marshal(dto.ThreadJoin{
		Code: "THREAD_JOIN",
		TID:  tid,
	})

	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudThreadCancel(tid uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	body, _ := json.Marshal(dto.ThreadCancel{
		Code: "THREAD_CANCEL",
		TID:  tid,
	})

	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudFinalizarHilo(pid uint32, tid uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.ThreadExit{
		Code: "THREAD_EXIT",
		TID:  tid,
		PID:  pid,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudMutexCreate(mutexID string) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.Mutex{
		Code:    "MUTEX_CREATE",
		MutexID: mutexID,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudMutexLock(mutexID string) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.Mutex{
		Code:    "MUTEX_LOCK",
		MutexID: mutexID,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudMutexUnlock(mutexID string) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.Mutex{
		Code:    "MUTEX_UNLOCK",
		MutexID: mutexID,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

// EnviarSolicitudIO se envía entrada y salida a kernel.
func EnviarSolicitudIO(tiempo uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.IO{
		Code:   "IO",
		Tiempo: tiempo,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}

func EnviarSolicitudDUMP(pid uint32, tid uint32) {
	slog.Info("Enviando SYSCALL a kernel")
	request := dto.Dump{
		Code: "DUMP",
		TID:  tid,
		PID:  pid,
	}

	body, _ := json.Marshal(request)
	_, err := rest.Send(urlSyscall, "POST", body)
	if err != nil {
		return
	}

}
