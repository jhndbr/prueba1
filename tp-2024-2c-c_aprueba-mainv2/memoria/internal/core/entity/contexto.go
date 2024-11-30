package entity

type ContextoEjecucion struct {
	PID    uint32 `json:"pid"` // ID del proceso
	TID    uint32 `json:"tid"` // ID del hilo
	AX     uint32 `json:"ax"`
	BX     uint32 `json:"bx"`
	CX     uint32 `json:"cx"`
	DX     uint32 `json:"dx"`
	EX     uint32 `json:"ex"`
	FX     uint32 `json:"fx"`
	GX     uint32 `json:"gx"`
	HX     uint32 `json:"hx"`
	PC     int    `json:"pc"`
	Base   uint32 `json:"base"`
	Limite uint32 `json:"limite"`
}
