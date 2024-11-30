package dto

type Syscall struct {
	Code              string `json:"code_id"`
	FileName          string `json:"file_name"`
	ProcessSize       uint32 `json:"process_size"`
	PrioridadHiloMain uint32 `json:"prioridad_hilo_main"`
	PID               uint32 `json:"pid"`
	TID               uint32 `json:"tid"`
	Priority          uint32 `json:"priority"`
	MutexID           string `json:"mutex_id"`
	Tiempo            uint32 `json:"tiempo"`
}
