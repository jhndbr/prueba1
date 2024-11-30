package planificador

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/repository"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/rest/cpu"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/gateway/rest/memory"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"github.com/sisoputnfrba/tp-golang/utils/list"
)

//la idea de este map es que el index es el TID y PID del hilo que se espera a que termina,
//y en su contenido hay una lista de hilos bloqueados por ello

// AtenderProcesosListos atiende los procesos listos para ejecutarse
func AtenderProcesosListos() {
	algoritmo = config.GetInstance().SchedulerAlgorithm

	if algoritmo == "FIFO" {
		planificarFIFO()
	} else if algoritmo == "PRIORIDADES" {
		planificarPrioridades()
	} else if algoritmo == "CMN" {
		planificarColasMultinivel()
	} else {
		slog.Error("El algoritmo %s no se logro identificar", algoritmo, "Error")
	}
}

func planificarFIFO() {
	_ = semReady.Wait()
	//asigna el hilo y proceso en ejecucion
	hiloEnEjecucion = repository.PopHiloREADY()
	procesoEnEjecucion = hiloEnEjecucion.Proceso
	slog.Info("Se corre el hilo", "PID:", hiloEnEjecucion.PID, "TID", hiloEnEjecucion.TID)
	for {

		if detenerPlanificacion == true {
			//compactacion
			_ = semPlanificacionDetenida.Post()
			_ = semCompactacionFinalizada.Wait()
		}

		if replanificar == true {
			_ = semReady.Wait()

			hiloEnEjecucion = repository.PopHiloREADY()
			if hiloEnEjecucion != nil {
				repository.AsignarHiloEXEC(hiloEnEjecucion)
				procesoEnEjecucion = hiloEnEjecucion.Proceso
				slog.Info("Se corre el hilo", "PID:", hiloEnEjecucion.PID, "TID:", hiloEnEjecucion.TID)
				replanificar = false
			} else {
				slog.Info("No hay hilos en Ready - Planificación Pausada")
				cantidad, _ := semReady.GetValue()
				for i := 0; i < cantidad; i++ {

					_ = semReady.Wait()

				}
				continue
			}

		}
		//envia el hilo a ejecutar en cpu
		cpu.EnviarIDsACPU(hiloEnEjecucion)

		_ = semPlanificador.Wait()
	}

}

func planificarPrioridades() {
	_ = semReady.Wait()
	hiloEnEjecucion = repository.PopHiloREADY()
	procesoEnEjecucion = hiloEnEjecucion.Proceso
	prioridadActual = hiloEnEjecucion.Priority
	slog.Info("Se corre el hilo", "PID:", hiloEnEjecucion.PID, "TID", hiloEnEjecucion.TID)

	for {
		hilos := repository.GetHilosBLOCKED()
		hilos0 := repository.GetHilosREADY()
		slog.Debug("hilos", "blocked:", hilos, "ready:", hilos0)
		//Esto se usa para compactacion
		if detenerPlanificacion == true {
			//indica que se detuvo la planificacion
			_ = semPlanificacionDetenida.Post()
			//espera a que termine la compactacion para renaudar
			_ = semCompactacionFinalizada.Wait()
		}

		//replanificar indica si hace falta tomar un hilo nuevo o si continua el anterior
		if replanificar == true {
			_ = semReady.Wait()

			//ordena los hilos por prioridad y obtiene el primero
			repository.OrdenarHilosREADYPorPrioridad()
			hiloEnEjecucion = repository.PopHiloREADY()
			hilos := repository.GetHilosBLOCKED()
			hilos0 := repository.GetHilosREADY()
			slog.Debug("hilos", "blocked:", hilos, "ready0:", hilos0)
			if hiloEnEjecucion != nil {

				//asignamos a EXEC el hilo y proceso a ejecutarse
				repository.AsignarHiloEXEC(hiloEnEjecucion)
				procesoEnEjecucion = hiloEnEjecucion.Proceso
				slog.Info("Se corre el hilo", "PID:", hiloEnEjecucion.PID, "TID:", hiloEnEjecucion.TID)

				//actualizamos la prioridad del hilo en ejecucion para manejar interrupciones
				prioridadActual = hiloEnEjecucion.Priority
				replanificar = false

			} else {
				slog.Info("No hay hilos en Ready - Planificación Pausada")
				//al no tener hilos que ejecutar, se reinicia a 0 el valor del semaforo, en caso de desincronizacion
				var cantidad, _ = semReady.GetValue()
				for i := 0; i < cantidad; i++ {
					_ = semReady.Wait()
				}
				continue
			}
		}
		//envia el pid y tid a cpu
		cpu.EnviarIDsACPU(hiloEnEjecucion)

		//espera a que se atienda alguna syscall(o detener por compactacion)
		_ = semPlanificador.Wait()

	}
}

func planificarColasMultinivel() {
	_ = semReady.Wait()
	hiloEnEjecucion = repository.PopHiloColaPrioridadX(prioridadActual)
	procesoEnEjecucion = hiloEnEjecucion.Proceso

	slog.Info("Se corre el hilo", "PID:", hiloEnEjecucion.PID, "TID", hiloEnEjecucion.TID)
	for {
		repository.FiltrarHilosDuplicadosREADY()

		if detenerPlanificacion == true {
			//compactacion
			_ = semPlanificacionDetenida.Post()
			_ = semCompactacionFinalizada.Wait()
		}

		if replanificar == true {
			_ = semReady.Wait()
			prioridadActual, err := repository.GetPrioridad(0)
			if err == "" {

				//busca un hilo en la cola de prioridad actual
				hiloEnEjecucion = repository.PopHiloColaPrioridadX(prioridadActual)
				if hiloEnEjecucion == nil {
					//si no encuentra un hilo, sube la prioridad y continua buscando
					repository.RemovePrioridad(0)

					//hace un post para que no se bloquee buscando el hilo
					_ = semReady.Post()
					continue
				}
			} else if err == "error" {
				slog.Info("No hay hilos en Ready - Planificación Pausada")
				//al no tener hilos que ejecutar, se reinicia a 0 el valor del semaforo, en caso de desincronizacion
				var cantidad, _ = semReady.GetValue()
				for i := 0; i < cantidad; i++ {
					_ = semReady.Wait()
				}
				continue
			}
			//asigna el hilo y el proceso EXEC
			repository.AsignarHiloEXEC(hiloEnEjecucion)
			procesoEnEjecucion = hiloEnEjecucion.Proceso

			slog.Info("Se corre el hilo", "PID:", hiloEnEjecucion.PID, "TID:", hiloEnEjecucion.TID)

			replanificar = false
		}
		//inicia el timer con el quantum del config
		cpu.EnviarIDsACPU(hiloEnEjecucion)
		go timer(hiloEnEjecucion)
		_ = semPlanificador.Wait()

		//hilos0 := repository.GetHilosCola0()
		//for i := 0; i < hilos0.Size(); i++ {
		//	hi, _ := hilos0.Get(i)
		//	fmt.Printf("COLA0: Posicion:%d, PID=%d, TID=%d\n", i, hi.PID, hi.TID)
		//}
		//hilos1 := repository.GetHilosCola1()
		//for i := 0; i < hilos1.Size(); i++ {
		//	hi, _ := hilos1.Get(i)
		//	fmt.Printf("COLA1: Posicion:%d, PID=%d, TID=%d\n", i, hi.PID, hi.TID)
		//}
		//hilosBLOCKED := repository.GetHilosBLOCKED()
		//for i := 0; i < hilosBLOCKED.Size(); i++ {
		//	hi, _ := hilosBLOCKED.Get(i)
		//	fmt.Printf("COLABLOCKED: Posicion:%d, PID=%d, TID=%d\n", i, hi.PID, hi.TID)
		//}
		//hilosEXIT := repository.GetHilosEXIT()
		//for i := 0; i < hilosEXIT.Size(); i++ {
		//	hi, _ := hilosEXIT.Get(i)
		//	fmt.Printf("COLABLOCKED: Posicion:%d, PID=%d, TID=%d\n", i, hi.PID, hi.TID)
		//}
	}
}

func AtenderMutexCreate(mutexID string) {
	//crea el nuevo mutex
	var mutex = *entity.CrearMutex(mutexID, procesoEnEjecucion.PID)
	repository.AddMutex(&mutex)
	//agrega el mutex al proceso en ejecucion
	procesoEnEjecucion.MUTEXs.Add(mutex)

	slog.Info("Nuevo Mutex", "Creado:", mutexID, "- PID:", procesoEnEjecucion.PID)

}

func AtenderMutexLock(mutexID string) {
	var mutex *entity.Mutex
	mutex, _ = repository.ObtenerMutex(mutexID, procesoEnEjecucion)
	if mutex.Estado == 0 {
		repository.LockMutex(mutexID, procesoEnEjecucion)
		slog.Info("Mutex en uso", "Mutex:", mutex.Nombre, "TID:", hiloEnEjecucion.TID)
	} else {
		slog.Info("#", "PID:", hiloEnEjecucion.PID, "TID:", hiloEnEjecucion.TID, "Bloqueado por:", mutex.Nombre)
		mutex.HilosBloqueados.Add(hiloEnEjecucion)
		repository.AgregarHiloEnBLOCKED(hiloEnEjecucion)

		replanificar = true

	}

}

func AtenderMutexUnlock(mutexID string) {
	var mutex *entity.Mutex
	mutex, _ = repository.ObtenerMutex(mutexID, procesoEnEjecucion)

	if mutex.Estado == 0 {
		slog.Error("Unlock innecesario", "Mutex:", mutex.Nombre)
		return
	}

	slog.Info("Mutex desbloqueado", "Mutex_id:", mutexID)
	repository.UnlockMutex(mutexID, procesoEnEjecucion)

	//desbloquea todos los hilos bloqueados por el mutex
	repository.DesbloquearHilosSegunAlgoritmo(mutex.HilosBloqueados)

	for i := 0; i < mutex.HilosBloqueados.Size(); i++ {
		hi, _ := mutex.HilosBloqueados.Get(i)
		repository.AddPrioridadEnOrden(hi.Priority)
		//hace un post por cada uno de los hilos bloqueados al planificador
		_ = semReady.Post()
	}
	mutex.HilosBloqueados = list.ArrayList[*entity.TCB]{}
	replanificar = true

}

func AtenderThreadExit(PID uint32, TID uint32) {

	if hiloEnEjecucion.TID != TID || hiloEnEjecucion.PID != PID {

		slog.Error("Se intenta desalojar un hilo distinto al de ejecucion.")
		return
	}
	//pasa el hilo a EXIT
	repository.AgregarHiloEnEXIT(hiloEnEjecucion)

	//solicita a memoria que finalize el hilo
	memory.SolicitarFinalizacionHilo(hiloEnEjecucion)

	//indica al planificador que elija un nuevo hilo
	replanificar = true

	slog.Info("Finaliza el hilo", "PID:", hiloEnEjecucion.PID, "TID:", hiloEnEjecucion.TID)

	ids := IDs{PID: PID, TID: TID}
	//Obtiene la lista de hilos bloqueados del map de ThreadJoin
	var hilos = mapThreadJoin[ids]

	if hilos.Size() > 0 {
		//desbloquea los hilos bloqueados
		repository.DesbloquearHilosSegunAlgoritmo(hilos)
		slog.Info("Hilos de threadjoin desbloqueados")
		for i := 0; i < hilos.Size(); i++ {
			hi, _ := hilos.Get(i)
			repository.AddPrioridadEnOrden(hi.Priority)
			//hace un post por cada uno de los hilos bloqueados al planificador
			_ = semReady.Post()
		}

		//setea la lista de bloqueados como vacia
		mapThreadJoin[ids] = list.ArrayList[*entity.TCB]{}
	}

}

func AtenderThreadCreate(prioridad uint32, file string) {

	var hilo = entity.CrearTCB(procesoEnEjecucion, prioridad)
	repository.AgregarHiloEnREADYSegunAlgoritmo(hilo)
	repository.AgregarAListaGlobal(hilo)
	repository.AddPrioridadEnOrden(prioridad)

	memory.EnviarNuevoHilo(hilo, file)
	if prioridad < prioridadActual && algoritmo != "FIFO" {
		prioridadActual = prioridad
		replanificar = true
		//reagrega el hilo en ejecucion
		repository.AgregarHiloEnREADYSegunAlgoritmo(hiloEnEjecucion)
		_ = semReady.Post()
	}
	slog.Info("Se crea el hilo - Estado: READY", "PID:", hilo.PID, "TID:", hilo.TID)
	_ = semReady.Post()

}

func AtenderThreadJoin(tid uint32) {

	ultimoTID := procesoEnEjecucion.TIDs.Size() - 1

	if uint32(ultimoTID) >= tid && !repository.EstaHiloEnEXIT(procesoEnEjecucion.PID, tid) {
		repository.AgregarHiloEnBLOCKED(hiloEnEjecucion)

		ids := IDs{PID: procesoEnEjecucion.PID, TID: tid}

		slog.Info("Bloqueado por: THREAD_JOIN", "PID:", hiloEnEjecucion.PID, "TID:", hiloEnEjecucion.TID, "tid del thread_join:", tid)
		replanificar = true

		//toma la lista de hilos bloqueados por threadjoin
		var hilos = mapThreadJoin[ids]
		//le agrega el hilo que llamo al TJ
		hilos.Add(hiloEnEjecucion)
		//y actualiza la lista del map
		mapThreadJoin[ids] = hilos

	}

}

func AtenderThreadCancel(tid uint32) {
	hilo, _ := repository.ObtenerHiloPorTIDdeProcesoX(tid, hiloEnEjecucion.PID)
	if hilo != nil {
		repository.AgregarHiloEnEXIT(hilo)
		memory.SolicitarFinalizacionHilo(hilo)
		slog.Info("Finaliza el hilo", "PID:", hilo.PID, "TID:", hilo.TID)

		//Obtiene la lista de hilos bloqueados del map de ThreadJoin
		ids := IDs{PID: hiloEnEjecucion.PID, TID: tid}
		var hilos = mapThreadJoin[ids]

		if hilos.Size() > 0 {
			//desbloquea los hilos bloqueados
			repository.DesbloquearHilosSegunAlgoritmo(hilos)
			for i := 0; i < hilos.Size(); i++ {

				hi, _ := hilos.Get(i)
				repository.AddPrioridadEnOrden(hi.Priority)
				//hace un post por cada uno de los hilos bloqueados al planificador
				_ = semReady.Post()

			}

			//setea la lista de bloqueados como vacia
			mapThreadJoin[ids] = list.ArrayList[*entity.TCB]{}
		}

		//verifica si el hilo finalizado fue el mismo que el que se ejecutaba
		if hilo.TID == hiloEnEjecucion.TID {
			replanificar = true

		}

	}

}

func AtenderProcessCreate(fileName string, memorySize uint32, priority uint32) error {

	pcb := entity.CrearPCB(fileName, memorySize, priority)
	if pcb == nil {
		return errors.New("error en crear el proceso")
	}
	repository.AgregarProcesoEnNEW(pcb)

	slog.Info("Se crea el proceso - Estado: NEW", "PID:", pcb.PID, "TID:", 0)
	_ = semProcesosNew.Post() // Libera un lugar

	return nil
}

func AtenderProcessExit(PID uint32) {

	if procesoEnEjecucion.PID != PID {
		slog.Error("El proceso", "PID:", PID, "Process.PID:", procesoEnEjecucion.PID)
		return
	}

	//pasa el proceso a EXIT
	repository.QuitarProcesoDeREADY(procesoEnEjecucion)
	repository.AgregarProcesoEnEXIT(procesoEnEjecucion)

	//termina todos los hilos del proceso
	repository.TerminarTodosLosHilos(procesoEnEjecucion)

	//enviar a memoria
	memory.SolicitarFinalizacionProceso(procesoEnEjecucion)

	slog.Info("Finaliza el proceso", "PID:", procesoEnEjecucion.PID)
	value, _ := semProcesosExit.GetValue()
	slog.Debug("semaforo", "value:", value)
	if largoPlazoDetenido == true {
		_ = semProcesosExit.Post()
	}
	value, _ = semProcesosExit.GetValue()
	slog.Debug("semaforo", "value:", value)

	replanificar = true

}

func AtenderDump(PID uint32, TID uint32) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/syscall-memoryDump", conf.IPMemory, conf.PortMemory)
	request := MemoryDumpRequest{PID: PID, TID: TID}
	body, _ := json.Marshal(request)

	resp, err := rest.Send(url, "POST", body)
	slog.Debug("Solicitud de dump de memoria enviada", "PID", PID, "TID", TID)

	repository.AgregarHiloEnBLOCKED(hiloEnEjecucion)

	if err != nil || resp.StatusCode != http.StatusOK {
		//pasa el proceso a exit
		AtenderProcessExit(PID)
		slog.Debug("Error en memory dump se mando el proceso a exit", "Error", err)
	} else {
		// Handle success: desbloquear hilo mover a ready
		var hiloEnLista list.ArrayList[*entity.TCB]
		hiloEnLista.Add(hiloEnEjecucion)
		repository.DesbloquearHilosSegunAlgoritmo(hiloEnLista)
		repository.AddPrioridadEnOrden(hiloEnEjecucion.Priority)
		_ = semReady.Post()
		replanificar = true
	}

	// Handle the response and error appropriately
	if err != nil {
		slog.Error("Error en memory dump request", "Error", err)
	} else {
		slog.Debug("Memory dump exitoso", "Status", resp.StatusCode)
	}

}

func AtenderIO(tiempo uint32) {
	// hilo bloqueado por I/O

	repository.AgregarHiloEnBLOCKED(hiloEnEjecucion)

	request := IORequest{Tiempo: tiempo, TCB: hiloEnEjecucion}
	ioQueue <- request

	//indica al planificador que continue y elija un nuevo hilo
	replanificar = true

}

func ProcessIOQueue() {
	for request := range ioQueue {
		slog.Info("Solicitud de I/O recibida", "TCB", request.TCB)
		wg.Add(1)
		go func(req IORequest) {
			defer wg.Done()
			slog.Debug("Solicitud IO recibida", "TCB", req.TCB)
			time.Sleep(time.Duration(req.Tiempo) * time.Millisecond)
			slog.Debug("Solicitud IO completada", "TCB", req.TCB)
			var hilo list.ArrayList[*entity.TCB]
			hilo.Add(req.TCB)
			repository.DesbloquearHilosSegunAlgoritmo(hilo)
			repository.AddPrioridadEnOrden(req.TCB.Priority)
			slog.Info("(<PID>:<TID>) finalizó IO y pasa a READY")
			if req.TCB.Priority < prioridadActual && algoritmo != "FIFO" {
				prioridadActual = req.TCB.Priority
				replanificar = true
				InterrumpirCPU(hiloEnEjecucion)
				repository.AgregarHiloEnREADYSegunAlgoritmo(hiloEnEjecucion)
			}
			_ = semReady.Post()
		}(request)
	}
}

// en el main iria asi
//
//	func main() {
//		go processIOQueue()
//
//		wg.Wait() // Espera a que todas las solicitudes de I/O se completen
//	}
//
//	func desbloquearHilo(hilo *entity.TCB) {
//		repository.QuitarHiloDeBLOCKED(hilo) // desbloquear hilo deberia eliminarlo de la cola de bloqueados
//
//		repository.AgregarHiloEnREADYSegunAlgoritmo(hilo)
//
// }
func timer(hilo *entity.TCB) {

	//convierte el quantum al formato de duration
	var quantum = config.GetInstance().Quantum

	//"duerme" por
	time.Sleep(time.Duration(quantum) * time.Millisecond)

	//if syscallAtendida == true {
	//
	//}
	//si el hilo no termino/se bloqueo durante su quantum, se lo agrega de nuevo a la cola de READY
	if hiloEnEjecucion == hilo {
		VuelveAReady()
		_ = semReady.Post()
		slog.Info("Desalojando por fin de quantum", "PID:", hilo.PID, "TID", hilo.TID)
		replanificar = true
		InterrumpirCPU(hilo)

	}

}

func InterrumpirCPU(hilo *entity.TCB) {
	//mSyscallInt.Lock()
	cpu.EnviarCPUInterrupt(hilo)

	interrupted = true

	_ = semInterrupter.Wait()

	interrupted = false

	SignalPlanificador()
	//mSyscallInt.Unlock()
}

func VuelveAReady() {
	if !repository.EstaHiloEnBLOCKED(hiloEnEjecucion.PID, hiloEnEjecucion.TID) && !repository.EstaHiloEnEXIT(hiloEnEjecucion.PID, hiloEnEjecucion.TID) && !repository.EstaHiloEnREADY(hiloEnEjecucion) {
		repository.AgregarHiloEnREADYSegunAlgoritmo(hiloEnEjecucion)
	}

}

var syscallAtendida = false
