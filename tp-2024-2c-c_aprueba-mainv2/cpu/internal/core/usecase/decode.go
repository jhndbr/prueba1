package usecase

import (
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"log/slog"
)

var DireccionFisica uint32
var huboSegmentation bool

func Decode(ctx *dto.Context) {
	switch ctx.IR.Code {
	case "READ_MEM":
		registroDireccion := ctx.IR.Args[1]
		offset := dto.ObtenerValorRegistro(ctx, registroDireccion)
		MMU(offset, ctx.Base, ctx.Limit)
		break
	case "WRITE_MEM":

		registroDireccion := ctx.IR.Args[0]
		offset := dto.ObtenerValorRegistro(ctx, registroDireccion)
		MMU(offset, ctx.Base, ctx.Limit)

		break
	case "SET":
		break
	case "SUM":
		break
	case "JNZ":
		break
	case "LOG":
		break
	case "MUTEX_CREATE":
		break
	case "MUTEX_LOCK":
		break
	case "MUTEX_UNLOCK":
		break
	case "DUMP_MEMORY":
		break
	case "IO":
		break
	case "PROCESS_CREATE":
		break
	case "THREAD_CREATE":
		break
	case "THREAD_CANCEL":
		break
	case "THREAD_JOIN":
		break
	case "THREAD_EXIT":
		break
	case "PROCESS_EXIT":
		break
	case "SUB":
		break
	default:
		slog.Error("Instruction Invalid.", slog.String("Code:", ctx.IR.Code))
		return
	}
	slog.Info("Decoded -", slog.String("Code:", ctx.IR.Code))

	Execute(ctx)

}
