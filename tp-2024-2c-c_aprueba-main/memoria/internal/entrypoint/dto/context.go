package dto

import (
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
)

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

type Instruction struct {
	Code string   `json:"code"` // Code del proceso
	Args []string `json:"args"` // Args
}

func NewContext(ctx *entity.ContextoEjecucion) *Context {
	return &Context{
		PID: ctx.PID,
		TID: ctx.TID,
		Register: &Register{
			AX: ctx.AX,
			BX: ctx.BX,
			CX: ctx.CX,
			DX: ctx.DX,
			EX: ctx.EX,
			FX: ctx.FX,
			GX: ctx.GX,
			HX: ctx.HX,
		},
		Base:  ctx.Base,
		Limit: ctx.Limite,
		PC:    uint32(ctx.PC),
		IR:    nil,
	}
}
