package entity

import (
	"log"
	"strconv"
	"sync"
)

var SeguirMutex sync.Mutex

var Seguir int

var CONTINUAR_CICLO = 0
var SYSCALL = 1
var DETENER_CICLO = 2

type Instruction struct { // (SET AX BX), etc
	Code string   `json:"code"` // Code del proceso
	Args []string `json:"args"` // Args
}
type InstruccionesPorHilo struct {
	Instrucciones []Instruction
}

type Archivo struct {
	PID      uint32
	FilePath string
}

var ListaArchivos []Archivo

// Instruccion_a_ejecutar GLOBAL
var Instruccion_a_ejecutar *Instruction

var SEM_RECIBIRINSTRUCCION = make(chan string)

func ObtenerValorArgumento(valorTexto string) uint32 {

	valor, err := strconv.ParseUint(valorTexto, 10, 32) // retorna un UINT64
	if err != nil {
		log.Fatalf("Error al convertir el valor: %v", err)
	}
	// despues guardo el valor en el registro correspondiente pero en formato uint32
	return uint32(valor)
}
