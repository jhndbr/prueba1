package repository

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/utils/list"
)

var mutexProcessNew sync.Mutex
var mutexProcessReady sync.Mutex
var procesosNEW = list.ArrayList[*entity.PCB]{}
var procesosEXIT = list.ArrayList[*entity.PCB]{}
var procesosREADY = list.ArrayList[*entity.PCB]{}
var procesoEXEC *entity.PCB

func AgregarProcesoEnNEW(pcb *entity.PCB) {
	mutexProcessNew.Lock()
	procesosNEW.Add(pcb)
	mutexProcessNew.Unlock()
}

func ObtenerProcesoEnNEW(index int) *entity.PCB {
	mutexProcessNew.Lock()
	item := ObtenerProceso(procesosNEW, index)
	mutexProcessNew.Unlock()
	return item
}

func ObtenerProceso(procesos list.ArrayList[*entity.PCB], index int) *entity.PCB {
	item, err := procesosNEW.Get(index)
	if err != nil {
		item = nil
	}
	return item
}

func PopProcesoEnNEW() *entity.PCB {
	mutexProcessNew.Lock()
	h, err := procesosNEW.Get(0)
	if err != nil {
		mutexProcessNew.Unlock()
		return nil
	}
	procesosNEW.Remove(0)
	mutexProcessNew.Unlock()
	return h
}

func ObtenerProcesosNEW() list.ArrayList[*entity.PCB] {
	return procesosNEW
}

func AsignarProcesoEXEC(pcb *entity.PCB) {
	procesoEXEC = pcb
}

func ObtenerProcesoEXEC() *entity.PCB {
	return procesoEXEC
}

func AgregarProcesoEnREADY(pcb *entity.PCB) {
	mutexProcessReady.Lock()
	procesosREADY.Add(pcb)
	mutexProcessReady.Unlock()
}

func QuitarProcesoDeREADY(pcb *entity.PCB) {
	mutexProcessReady.Lock()
	for i := 0; i < procesosREADY.Size(); i++ {
		pi, _ := procesosREADY.Get(i)
		if pi.PID == pcb.PID {
			procesosREADY.Remove(i)
		}
	}
	mutexProcessReady.Unlock()
}

func AgregarProcesoEnEXIT(pcb *entity.PCB) {
	procesosEXIT.Add(pcb)
}

// Definimos una función de comparación para filtrar los pid
func isEqualThan(a, b uint32) bool {
	return b == a
}

func ObtenerProcesoREADYPorPID(pid uint32) *entity.PCB {

	for i := 0; i < procesosREADY.Size(); i++ {
		item := ObtenerProceso(procesosREADY, i)
		if isEqualThan(item.PID, pid) {
			return item
		}
	}

	return nil
}

//func AsignarProcesoEXECPorPID(pid uint32) error {
//
//	for i := 0; i < procesosREADY.Size(); i++ {
//		item := ObtenerProceso(procesosREADY, i)
//		if isEqualThan(item.PID, pid) {
//			AsignarProcesoEXEC(item)
//			return nil
//		}
//	}
//
//	return errors.New("PID not found")
//}
