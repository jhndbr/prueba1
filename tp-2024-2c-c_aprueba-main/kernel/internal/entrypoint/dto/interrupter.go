package dto

type Interrupt struct {
	PID    uint32 `json:"pid"`
	TID    uint32 `json:"tid"`
	Motivo string `json:"motivo"`
}
