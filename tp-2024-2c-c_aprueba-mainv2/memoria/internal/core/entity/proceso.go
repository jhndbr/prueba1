package entity

import "time"

// ProcessState define los posibles estados de un proceso
type ProcessState int

// Definición de los estados del proceso utilizando iota
const (
	NEW        ProcessState = iota
	READY                   // Estado de listo para ejecutarse
	EXEC                 // Estado de ejecución
	BLOCKED                 // Estado de bloqueado esperando un recurso
	Suspended               // Estado de suspendido
	EXIT              // Estado de terminado
)

// String permite que se muestre un nombre legible para cada estado
func (state ProcessState) String() string {
	return [...]string{"NEW", "READY", "EXEC", "BLOCKED", "Suspended", "EXIT"}[state]
}

// PCB representa un bloque de control de procesos
type PCB struct {
	PID               uint32        // Identificador único del proceso
	TID               []uint32      // Lista de los identificadores de los hilos asociados al proceso
	State             ProcessState  // Estado actual del proceso
	Priority          uint32        // Prioridad del proceso
	PC                uint64        // Contador de programa (registro de PC)
	CPURegisters      [8]int        // Simulación de registros del CPU (generalmente más)
	MemoryBase        uint64        // Dirección de inicio del espacio de memoria
	MemoryLimit       uint64        // Dirección de fin del espacio de memoria
	OpenFiles         []string      // Lista de archivos abiertos
	CPUTimeUsed       time.Duration // Tiempo total de CPU usado por el proceso
	CreationTime      time.Time     // Momento de creación del proceso
	LastExecutionTime time.Time     // Momento de la última ejecución
}
