package usecase

import (
	"log/slog"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/kernel"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/memoria"
)

func Execute(ctx *dto.Context) {

	ctx.PC++
	slog.Info("instruccion Ejecutada", "TID", ctx.TID, "EJECUTANDO - Instruccion:", ctx.IR.Code, "parametros:", ctx.IR.Args)
	switch ctx.IR.Code {

	case "SET": // set "ax" "1"
		registro := ctx.IR.Args[0]
		// Convierte el valor STRING a UINT64
		value, err := strconv.ParseUint(ctx.IR.Args[1], 10, 32)
		if err != nil {
			return
		}
		dto.ModificarContexto(ctx, registro, uint32(value))

		break

	case "READ_MEM":
		if !huboSegmentation {
			slog.Info("SOLICITANDO VALOR DE MEMORIA", "direccion", DireccionFisica)
			valor, err := memoria.EnviarSolicitudLeerValorMemoria(DireccionFisica)
			if err != nil {
				return
			}
			slog.Info("VALOR RECIBIDO DE LA SOLICITUD READ_MEM", "valor", valor)
			registroDatos := ctx.IR.Args[0]
			dto.ModificarContexto(ctx, registroDatos, *valor)
		} else {
			slog.Error("Error SEGMENTATION FAULT")
			err := memoria.ActualizarContextoDeEjecucion(ctx)
			if err != nil {
				return
			}
			kernel.EnviarMotivoInterrupcion(ctx.PID, ctx.TID, "SEGMENTATION FAULT")

			return // no sigue ejecutando
		}

		break

	case "WRITE_MEM":
		if !huboSegmentation {
			registroDestino := ctx.IR.Args[1]
			valor := dto.ObtenerValorRegistro(ctx, registroDestino)
			slog.Info("VALOR A ESCRIBIR EN MEMORIA", "valor", valor)
			slog.Info("VALOR ENVIADO A MEMORIA", "direccion", DireccionFisica)
			err := memoria.EnviarSolicitudEscribirEnMemoria(DireccionFisica, valor)
			if err != nil {
				return
			}
		} else {
			slog.Error("Error SEGMENTATION FAULT")
			err := memoria.ActualizarContextoDeEjecucion(ctx)
			if err != nil {
				slog.Error("Error", "Err:", err.Error())
				return
			}
			kernel.EnviarMotivoInterrupcion(ctx.PID, ctx.TID, "SEGMENTATION FAULT")

			return // no sigue ejecutando
		}

		break

	case "SUM":
		registroDestino := ctx.IR.Args[0]
		registroOrigen := ctx.IR.Args[1]
		valorDestino := dto.ObtenerValorRegistro(ctx, registroDestino)
		valorOrigen := dto.ObtenerValorRegistro(ctx, registroOrigen)
		total := valorDestino + valorOrigen
		dto.ModificarContexto(ctx, registroDestino, total)
		break

	case "SUB":
		registroDestino := ctx.IR.Args[0]
		registroOrigen := ctx.IR.Args[1]
		valorDestino := dto.ObtenerValorRegistro(ctx, registroDestino)
		valorOrigen := dto.ObtenerValorRegistro(ctx, registroOrigen)
		total := valorDestino - valorOrigen
		dto.ModificarContexto(ctx, registroDestino, total)
		break

	case "JNZ":
		registro := ctx.IR.Args[0]
		newValuePC, err := strconv.ParseUint(ctx.IR.Args[1], 10, 32)
		if err != nil {
			return
		}
		valorRegistro := dto.ObtenerValorRegistro(ctx, registro)
		if valorRegistro != 0 {
			dto.ModificarContexto(ctx, "PC", uint32(newValuePC))
		}
		break

	case "LOG":
		registro := ctx.IR.Args[0]
		valorRegistro := dto.ObtenerValorRegistro(ctx, registro)
		slog.Info("LOG:", "registro:", registro, "valor:", valorRegistro)
		break

	case "DUMP_MEMORY":
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudDUMP(ctx.PID, ctx.TID)
		return

	case "IO":
		tiempo, err := strconv.ParseUint(ctx.IR.Args[0], 10, 32)
		if err != nil {
			return
		}
		err = memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudIO(uint32(tiempo))
		return
		//break

	case "PROCESS_CREATE":
		fileInstruction := ctx.IR.Args[0]
		size, err := strconv.ParseUint(ctx.IR.Args[1], 10, 32)
		if err != nil {
			return
		}
		priority, err := strconv.ParseUint(ctx.IR.Args[2], 10, 32)
		if err != nil {
			return
		}

		err = memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudCrearProceso(fileInstruction, uint32(size), uint32(priority))
		return

	case "PROCESS_EXIT":
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudFinalizarProceso(ctx.PID)
		return

	case "THREAD_CREATE":
		fileInstruction := ctx.IR.Args[0]
		priority, err := strconv.ParseUint(ctx.IR.Args[1], 10, 32)
		if err != nil {
			return
		}
		err = memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudCrearHilo(fileInstruction, uint32(priority))
		return

	case "THREAD_JOIN":
		tid, err := strconv.ParseUint(ctx.IR.Args[0], 10, 32)
		if err != nil {
			return
		}
		err = memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudThreadJoin(uint32(tid))
		return

	case "THREAD_CANCEL":
		tid, err := strconv.ParseUint(ctx.IR.Args[0], 10, 32)
		if err != nil {
			return
		}
		err = memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudThreadCancel(uint32(tid))
		return

	case "THREAD_EXIT":
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudFinalizarHilo(ctx.PID, ctx.TID)
		return

	case "MUTEX_CREATE":
		mutexId := ctx.IR.Args[0]
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudMutexCreate(mutexId)
		return

	case "MUTEX_LOCK":
		mutexId := ctx.IR.Args[0]
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudMutexLock(mutexId)
		return

	case "MUTEX_UNLOCK":
		mutexId := ctx.IR.Args[0]
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarSolicitudMutexUnlock(mutexId)
		return

	default:

		slog.Error("Instruction Invalid", slog.String("code=", ctx.IR.Code))
	}

	//slog.Info("CONTEXTO DESPUES DEL EXECUTE:", "registros=", fmt.Sprintf("%+v", ctx.Register))

	if dto.Interrumpir {
		slog.Info("Desalojando hilo")
		err := memoria.ActualizarContextoDeEjecucion(ctx)
		if err != nil {
			return
		}
		kernel.EnviarMotivoInterrupcion(ctx.PID, ctx.TID, "DESALOJO")
	} else {
		Fetch(ctx)
	}

}
