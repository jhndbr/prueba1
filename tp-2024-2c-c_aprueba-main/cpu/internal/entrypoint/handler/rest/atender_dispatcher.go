package rest

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/memoria"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/usecase"
	"io/ioutil"
	"log/slog"
	"net/http"
)

func instructionCicle(ctx *dto.Context) {

	// Busco en Memoria el Contexto de Ejecución
	ctx, err := memoria.ObtenerContextoDeEjecucion(ctx.PID, ctx.TID)

	if err != nil {
		return
	}
	slog.Info("contexto obtenido por memoria")

	// inicializo esta variable GLOBAL para comparar en el hilo CHECKINTERRUMP si el TID que manda kernel a interrumpir es el mismo que el que se esta ejecutando
	dto.CtxEnEjecucion = ctx

	// Busco en Memoria próxima instrucción e inicio el ciclo
	usecase.Fetch(ctx)

	/* LO COMENTO PQ LO PASE AL FETCH QUE AGREGASTE POR SI NO HAY MAS INSTRUCCIONES A EJECUTAR*/
	// Actualizo contexto en Memoria con lo ultímo ejecutado
	//err = memoria.ActualizarContextoDeEjecucion(ctx)
	//if err != nil {
	//	kernel.EnviarMotivoInterrupcion(ctx.TID, err.Error())
	//	return
	//}
	//
	//// Finalizo el proceso
	//kernel.EnviarMotivoInterrupcion(ctx.TID, "EXIT")
}

// AtenderDispatcher recibe las solicitudes de procesos a ejecutar
func AtenderDispatcher(w http.ResponseWriter, r *http.Request) {

	// Verifico el methods sea un post
	if r.Method != http.MethodPost {
		http.Error(w, "Method has not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request is Empty or has an erro", http.StatusInternalServerError)
		return
	}

	// Deserealizo el JSON a la estructura
	var context *dto.Context
	err = json.Unmarshal(body, &context)
	if err != nil {
		http.Error(w, "Request has not been deserialize", http.StatusBadRequest)
		return
	}

	slog.Info("RECIBO PROCESO A EJECUTAR", "PID=", context.PID, "TID=", context.TID)

	// Ejecuto en un hilo el ciclo de instrucción sin bloquear el servidor.
	go instructionCicle(context)

	// Respondo solicitud
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
