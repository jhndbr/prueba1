package contract

type Register struct {
	AX uint32 `json:"ax"` // AX general purpose register
	BX uint32 `json:"bx"` // BX general purpose register
	CX uint32 `json:"cx"` // CX general purpose register
	DX uint32 `json:"dx"` // DX general purpose register
	EX uint32 `json:"ex"` // EX general purpose register
	FX uint32 `json:"fx"` // FX general purpose register
	GX uint32 `json:"gx"` // GX general purpose register
	HX uint32 `json:"hx"` // HX general purpose register
	PC uint32 `json:"pc"` // PC Program Counter register
	IR string `json:"ir"` // IR instruction register
}

type Context struct {
	PID      uint32   `json:"pid"`      // PID Process ID
	TID      uint32   `json:"tid"`      // TID Thread ID
	Register Register `json:"register"` // Register
	Base     uint32   `json:"base"`     // Base memory
	Limit    uint32   `json:"limit"`    // Limit memory
}
