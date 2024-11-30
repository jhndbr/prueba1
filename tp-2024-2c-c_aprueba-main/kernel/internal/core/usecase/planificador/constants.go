package planificador

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/utils/list"
	"github.com/sisoputnfrba/tp-golang/utils/semaphore"
)

type MemoryDumpRequest struct {
	PID uint32 `json:"pid"`
	TID uint32 `json:"tid"`
}
type IORequest struct {
	Tiempo uint32      `json:"tiempo"`
	TCB    *entity.TCB `json:"tcb"`
}

type IDs struct {
	PID uint32
	TID uint32
}

// SEMAFOROS
var semProcesosExit semaphore.Semaphore
var semProcesosNew semaphore.Semaphore
var semReady semaphore.Semaphore
var semPlanificador semaphore.Semaphore
var semPlanificacionDetenida semaphore.Semaphore
var semCompactacionFinalizada semaphore.Semaphore
var semInterrupter semaphore.Semaphore
var wg sync.WaitGroup
var mSyscallInt sync.Mutex

// VARIABLES LOCALES
var hiloEnEjecucion *entity.TCB = nil
var procesoEnEjecucion *entity.PCB
var prioridadActual uint32 = 99999999
var algoritmo string
var quantumEnd = false

// BOOLS
var interrupted bool = false
var largoPlazoDetenido = false
var replanificar = false
var detenerPlanificacion = false

// LISTAS Y MAPS
var ioQueue = make(chan IORequest, 100)
var mapThreadJoin = make(map[IDs]list.ArrayList[*entity.TCB])

//funcs

func InicializarSemaforos() {

	_ = semProcesosExit.Open("EXIT", 0644, 0)
	_ = semProcesosNew.Open("NEW", 0644, 0)
	_ = semReady.Open("READY", 0644, 0)
	_ = semPlanificador.Open("PLANIFICADOR", 0644, 0)
	_ = semPlanificacionDetenida.Open("PLANIFICACIONDETENIDA", 0644, 0)
	_ = semCompactacionFinalizada.Open("COMPACTACIONFINALIZADA", 0644, 0)
	_ = semInterrupter.Open("INTERRUPTER", 0644, 0)
}

func CerrarSemaforos() {
	_ = semProcesosExit.Close()
	_ = semProcesosNew.Close()
	_ = semReady.Close()
	_ = semPlanificador.Close()
	_ = semPlanificacionDetenida.Close()
	_ = semCompactacionFinalizada.Close()
	_ = semInterrupter.Close()

	_ = semProcesosExit.Unlink()
	_ = semProcesosNew.Unlink()
	_ = semReady.Unlink()
	_ = semPlanificador.Unlink()
	_ = semPlanificacionDetenida.Unlink()
	_ = semCompactacionFinalizada.Unlink()
	_ = semInterrupter.Unlink()
}

func VerificarSemaforos() {
	valueREADY, _ := semReady.GetValue()
	valuePEXIT, _ := semProcesosExit.GetValue()
	valuePNEW, _ := semProcesosNew.GetValue()
	valuePLAN, _ := semPlanificador.GetValue()
	valuePD, _ := semPlanificacionDetenida.GetValue()
	valueCF, _ := semCompactacionFinalizada.GetValue()
	valueINT, _ := semInterrupter.GetValue()

	if valueREADY != 0 || valuePEXIT != 0 || valuePNEW != 0 || valuePLAN != 0 || valuePD != 0 || valueCF != 0 || valueINT != 0 {
		CerrarSemaforos()
		InicializarSemaforos()
	}
}

func SignalInterrupt() {
	_ = semInterrupter.Post()
}

func SignalPlanificador() {
	value, _ := semPlanificador.GetValue()
	if value == 0 {
		_ = semPlanificador.Post()
	}
}

func SetInterrupted(b bool) {
	interrupted = b
}
func GetInterrupted() bool {
	return interrupted
}

func LockMutexSyscall() {
	mSyscallInt.Lock()
}

func UnlockMutexSyscall() {
	mSyscallInt.Unlock()
}

func IndicarAtencionSyscall(a bool) {
	syscallAtendida = a
}
