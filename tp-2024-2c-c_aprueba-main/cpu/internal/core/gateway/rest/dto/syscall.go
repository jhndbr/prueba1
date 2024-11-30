package dto

type ProcessCreate struct {
	Code              string `json:"code_id"`
	FileName          string `json:"file_name"`
	ProcessSize       uint32 `json:"process_size"`
	PrioridadHiloMain uint32 `json:"prioridad_hilo_main"`
}
type ProcessExit struct {
	Code string `json:"code_id"`
	PID  uint32 `json:"pid"`
}
type ThreadCreate struct {
	Code     string `json:"code_id"`
	FileName string `json:"file_name"`
	Priority uint32 `json:"priority"`
}
type ThreadJoin struct {
	Code string `json:"code_id"`
	TID  uint32 `json:"tid"`
}
type ThreadCancel struct {
	Code string `json:"code_id"`
	TID  uint32 `json:"tid"`
}
type ThreadExit struct {
	Code string `json:"code_id"`
	TID  uint32 `json:"tid"`
	PID  uint32 `json:"pid"`
}
type Mutex struct {
	Code    string `json:"code_id"`
	MutexID string `json:"mutex_id"`
}
type IO struct {
	Code   string `json:"code_id"`
	Tiempo uint32 `json:"tiempo"`
}
type Dump struct {
	Code string `json:"code_id"`
	TID  uint32 `json:"tid"`
	PID  uint32 `json:"pid"`
}
