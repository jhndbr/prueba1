package dto

type Register struct {
	AX uint32 `json:"ax"` // AX general purpose register
	BX uint32 `json:"bx"` // BX general purpose register
	CX uint32 `json:"cx"` // CX general purpose register
	DX uint32 `json:"dx"` // DX general purpose register
	EX uint32 `json:"ex"` // EX general purpose register
	FX uint32 `json:"fx"` // FX general purpose register
	GX uint32 `json:"gx"` // GX general purpose register
	HX uint32 `json:"hx"` // HX general purpose register
}

type Context struct {
	PID      uint32       `json:"pid"`      // PID Process ID
	TID      uint32       `json:"tid"`      // TID Thread ID
	Register *Register    `json:"register"` // Register
	Base     uint32       `json:"base"`     // Base memory
	Limit    uint32       `json:"limit"`    // Limit memory
	PC       uint32       `json:"pc"`       // PC Program Counter register
	IR       *Instruction `json:"ir"`       // IR instruction register
}

var CtxEnEjecucion *Context

func ModificarContexto(contexto *Context, registro string, valor uint32) {
	switch registro {
	case "AX":
		contexto.Register.AX = valor
	case "BX":
		contexto.Register.BX = valor
	case "CX":
		contexto.Register.CX = valor
	case "DX":
		contexto.Register.DX = valor
	case "EX":
		contexto.Register.EX = valor
	case "FX":
		contexto.Register.FX = valor
	case "GX":
		contexto.Register.GX = valor
	case "HX":
		contexto.Register.HX = valor
	case "PC":
		contexto.PC = valor
	}
}
func ObtenerValorRegistro(contexto *Context, registro string) uint32 {
	switch registro {
	case "AX":
		return contexto.Register.AX
	case "BX":
		return contexto.Register.BX
	case "CX":
		return contexto.Register.CX
	case "DX":
		return contexto.Register.DX
	case "EX":
		return contexto.Register.EX
	case "FX":
		return contexto.Register.FX
	case "GX":
		return contexto.Register.GX
	case "HX":
		return contexto.Register.HX
	case "PC":
		return contexto.PC
	}
	return 0
}
