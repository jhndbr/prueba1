package repository

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/list"
)

var mREADY sync.Mutex
var mBLOCKED sync.Mutex
var mEXIT sync.Mutex
var hilosREADY list.ArrayList[*entity.TCB]
var hilosBLOCKED list.ArrayList[*entity.TCB]
var hilosEXIT list.ArrayList[*entity.TCB]
var hiloEXEC *entity.TCB
var listaGlobalDeHilos list.ArrayList[*entity.TCB]

func AgregarAListaGlobal(hilo *entity.TCB) {
	listaGlobalDeHilos.Add(hilo)
}

func EstaEnListaGlobal(hilo *entity.TCB) bool {
	for i := 0; i < listaGlobalDeHilos.Size(); i++ {
		var hi = ObtenerHiloDeLista(listaGlobalDeHilos, i)
		if hi.TID == hilo.TID && hi.PID == hilo.PID {
			return true
		}
	}
	return false
}

var prioridadesEnUso list.ArrayList[uint32]
var algoritmo = config.GetInstance().SchedulerAlgorithm

// mapColas de las colas Ready
var mapColas = map[uint32]list.ArrayList[*entity.TCB]{}

func AgregarHiloEnREADY(tcb *entity.TCB) {
	mREADY.Lock()
	hilosREADY.Add(tcb)
	mREADY.Unlock()
}

func AgregarHiloEnREADYSegunAlgoritmo(tcb *entity.TCB) {
	mREADY.Lock()
	switch algoritmo {
	case "FIFO":
		hilosREADY.Add(tcb)
	case "PRIORIDADES":
		hilosREADY.Add(tcb)
	case "CMN":
		AgregarAColaDePrioridadX(tcb.Priority, tcb)
	}
	mREADY.Unlock()
}

func QuitarHiloDeREADY(tcb *entity.TCB) {

	mREADY.Lock()
	if algoritmo == "CMN" {
		QuitarDeColaDePrioridadX(tcb.Priority, tcb)
	} else {
		for i := 0; i < hilosREADY.Size(); i++ {
			hi, _ := hilosREADY.Get(i)
			if hi.TID == tcb.TID && hi.PID == tcb.PID {
				hilosREADY.Remove(i)
			}
		}
	}
	mREADY.Unlock()
}

func AgregarHiloEnBLOCKED(tcb *entity.TCB) {
	mBLOCKED.Lock()
	hilosBLOCKED.Add(tcb)
	mBLOCKED.Unlock()
}

func QuitarHiloDeBLOCKED(tcb *entity.TCB) {
	mBLOCKED.Lock()
	for i := 0; i < hilosBLOCKED.Size(); i++ {
		hi, _ := hilosBLOCKED.Get(i)
		if hi.TID == tcb.TID && hi.PID == tcb.PID {
			hilosBLOCKED.Remove(i)
		}
	}
	mBLOCKED.Unlock()
}

func AgregarHiloEnEXIT(tcb *entity.TCB) {
	mEXIT.Lock()
	hilosEXIT.Add(tcb)
	mEXIT.Unlock()
}

func AsignarHiloEXEC(tcb *entity.TCB) {
	hiloEXEC = tcb
}

func ObtenerHiloDeLista(hilos list.ArrayList[*entity.TCB], index int) *entity.TCB {
	h, err := hilos.Get(index)
	if err != nil {
		h = nil
	}
	return h
}

func PopHiloREADY() *entity.TCB {
	mREADY.Lock()
	h, err := hilosREADY.Get(0)
	if err != nil {
		mREADY.Unlock()
		return nil
	}
	hilosREADY.Remove(0)
	mREADY.Unlock()
	return h
}

func ObtenerHiloREADY(index int) *entity.TCB {
	mREADY.Lock()
	h := ObtenerHiloDeLista(hilosREADY, index)
	mREADY.Unlock()
	return h
}

func ObtenerHiloBLOCKED(index int) *entity.TCB {
	mBLOCKED.Lock()
	h := ObtenerHiloDeLista(hilosBLOCKED, index)
	mBLOCKED.Unlock()
	return h
}

// PorPrioridad esto prox es la logica para ordenar una lista de hilos por prioridad
type PorPrioridad []*entity.TCB

func (a PorPrioridad) Len() int {
	return len(a)
}

func (a PorPrioridad) Less(i, j int) bool {
	return a[i].Priority < a[j].Priority
}

func (a PorPrioridad) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func OrdenarHilosREADYPorPrioridad() {

	// Ordenar la lista en orden ascendente
	hilosREADY.Sort(func(a, b *entity.TCB) bool {
		return a.Priority <= b.Priority // Orden ascendente de Prioridades
	})

}

//func ObtenerTodosLosHilosSegunAlgoritmo() list.ArrayList[*entity.TCB] {
//	algoritmo := config.GetInstance().SchedulerAlgorithm
//
//	var hilos list.ArrayList[*entity.TCB]
//
//	if algoritmo == "FIFO" || algoritmo == "PRIORIDADES" {
//		hilos = hilosREADY
//		hilos.Add(hiloEXEC)
//		for i := 0; i < hilosBLOCKED.Size(); i++ {
//			var hi = ObtenerHiloBLOCKED(i)
//			hilos.Add(hi)
//		}
//
//	} else if algoritmo == "CMN" {
//
//		hilos = ObtenerHilosREADYDeTodasLasColas()
//		for i := 0; i < hilosBLOCKED.Size(); i++ {
//			var hi = ObtenerHiloBLOCKED(i)
//			hilos.Add(hi)
//		}
//
//	}
//	return hilos
//}

func ObtenerHiloPorTIDdeProcesoX(tid uint32, pid uint32) (*entity.TCB, error) {

	for i := 0; i < listaGlobalDeHilos.Size(); i++ {
		var hi = ObtenerHiloDeLista(listaGlobalDeHilos, i)
		if hi.TID == tid && hi.PID == pid {
			return hi, nil
		}
	}
	return nil, nil
}

func DesbloquearHilosSegunAlgoritmo(hilos list.ArrayList[*entity.TCB]) {

	for i := 0; i < hilos.Size(); i++ {
		var hilo = ObtenerHiloDeLista(hilos, i)
		AgregarHiloEnREADYSegunAlgoritmo(hilo)
		QuitarHiloDeBLOCKED(hilo)
	}
}

func TerminarTodosLosHilos(proceso *entity.PCB) {
	var tids = proceso.TIDs
	for i := 0; i < tids.Size(); i++ {
		var hilo, _ = ObtenerHiloPorTIDdeProcesoX(uint32(i), proceso.PID)
		AgregarHiloEnEXIT(hilo)
		QuitarHiloDeREADY(hilo)
	}
}

func ObtenerhiloEXEC() *entity.TCB {
	return hiloEXEC
}

//FUNCIONES COLAS MULTINIVEL

func ObtenerHiloColaPrioridadX(prioridad uint32, index int) *entity.TCB {
	var cola = mapColas[prioridad]
	var hilo = ObtenerHiloDeLista(cola, index)

	return hilo
}

func PopHiloColaPrioridadX(prioridad uint32) *entity.TCB {
	var cola = mapColas[prioridad]
	var hilo = ObtenerHiloDeLista(cola, 0)
	if hilo != nil {
		cola.Remove(0)
		mapColas[prioridad] = cola
	}
	return hilo
}

func AgregarAColaDePrioridadX(prioridad uint32, hilo *entity.TCB) {
	var cola = mapColas[prioridad]
	//for i := 0; i < cola.Size(); i++ {
	//	var hi = ObtenerHiloDeLista(cola, i)
	//	if EstaHiloEnBLOCKED(hi.PID, hi.TID) || EstaHiloEnREADY(hi) {
	//		return
	//	}
	//}
	cola.Add(hilo)
	mapColas[prioridad] = cola
}

func EstaHiloEnEXIT(pid uint32, tid uint32) bool {
	for i := 0; i < hilosEXIT.Size(); i++ {
		var hi = ObtenerHiloDeLista(hilosEXIT, i)
		if hi.TID == tid && hi.PID == pid {
			return true
		}
	}
	return false
}

func EstaHiloEnBLOCKED(pid uint32, tid uint32) bool {
	for i := 0; i < hilosBLOCKED.Size(); i++ {
		var hi = ObtenerHiloDeLista(hilosBLOCKED, i)
		if hi.TID == tid && hi.PID == pid {
			return true
		}
	}
	return false
}

func EstaHiloEnREADY(tcb *entity.TCB) bool {
	hilos := mapColas[tcb.Priority]
	for i := 0; i < hilos.Size(); i++ {
		var hi = ObtenerHiloDeLista(hilos, i)
		if hi.TID == tcb.TID && hi.PID == tcb.PID {
			return true
		}
	}
	return false
}

func EsHiloUnico(prioridad uint32) bool {
	var cola = mapColas[prioridad]
	if cola.Size() == 0 {
		return true
	} else {
		return false
	}
}

func ObtenerHilosREADYDeTodasLasColas() list.ArrayList[*entity.TCB] {
	var hilos list.ArrayList[*entity.TCB]
	for i := 0; i < prioridadesEnUso.Size(); i++ {
		var prioridad, _ = prioridadesEnUso.Get(i)
		cola := mapColas[prioridad]
		for i := 0; i < cola.Size(); i++ {
			var hi = ObtenerHiloColaPrioridadX(prioridad, i)
			hilos.Add(hi)
		}
	}
	return hilos
}

func QuitarDeColaDePrioridadX(prioridad uint32, hilo *entity.TCB) {
	var cola = mapColas[prioridad]
	for i := 0; i < cola.Size(); i++ {
		hi, _ := cola.Get(i)
		if hi.TID == hilo.TID && hi.PID == hilo.PID {
			cola.Remove(i)
		}
	}
	mapColas[prioridad] = cola
}

func OrdenarPrioridades() {
	prioridadesEnUso.Sort(func(a, b uint32) bool {
		return a < b // Orden ascendente de Prioridades
	})
}

func GetPrioridades() list.ArrayList[uint32] {
	return prioridadesEnUso
}

func GetPrioridad(index int) (uint32, string) {
	prioridad, err := prioridadesEnUso.Get(index)
	if err != nil {
		return 0, "error"
	} else {
		return prioridad, ""
	}
}

func RemovePrioridad(index int) {
	prioridadesEnUso.Remove(index)
}

func AddPrioridadEnOrden(prioridad uint32) {
	OrdenarPrioridades()
	for i := 0; i < prioridadesEnUso.Size(); i++ {
		pi, _ := prioridadesEnUso.Get(i)
		if pi == prioridad {
			return
		}
	}
	prioridadesEnUso.Add(prioridad)
	OrdenarPrioridades()
}

// funciones debug
func GetHilosREADY() list.ArrayList[*entity.TCB] {
	return hilosREADY
}

func GetHilosBLOCKED() list.ArrayList[*entity.TCB] {
	return hilosBLOCKED
}

func GetHilosEXIT() list.ArrayList[*entity.TCB] {
	return hilosEXIT
}

func GetHilosCola0() list.ArrayList[*entity.TCB] {
	return mapColas[uint32(0)]
}

func GetHilosCola1() list.ArrayList[*entity.TCB] {
	return mapColas[uint32(1)]
}

func GetHilosCola5() list.ArrayList[*entity.TCB] {
	return mapColas[uint32(5)]
}

func GetHilosCola6() list.ArrayList[*entity.TCB] {
	return mapColas[uint32(6)]
}

func FiltrarHilosDuplicadosREADY() {
	var hilosPresentes = ObtenerHilosREADYDeTodasLasColas()
	hilosPresentes.Sort(func(a, b *entity.TCB) bool {
		if a.PID != b.PID {
			return a.PID <= b.PID
		} else if a.PID == b.PID {
			return a.TID <= b.TID
		}
		return false
	})

	for i := 1; i < hilosPresentes.Size(); i++ {
		h0 := ObtenerHiloDeLista(hilosPresentes, i-1)
		h1 := ObtenerHiloDeLista(hilosPresentes, i)
		if h0.PID == h1.PID && h0.TID == h1.TID {
			QuitarDeColaDePrioridadX(h0.Priority, h0)
		}
	}

}
