package entity

import "sync"

type Particion struct {
	Base   int
	Limite int
	Libre  bool
	PID    uint32
}

var Particiones []*Particion

var MemoriaUsuario []byte
var MemMutex sync.Mutex

func InicializarMemoriaUsuario(tam int) {
	MemoriaUsuario = make([]byte, tam)
}

// Memoria de sistema donde se almacenan los contextos, tambien se deberia almacenar el archivo de instrucciones?
var MemoriaSistema []ContextoEjecucion

// Inicializamos la memoria de sistema

func InicializarMemoriaSistema(tam int) {
	MemoriaSistema = make([]ContextoEjecucion, 0, tam) // Definir la capacidad
}
