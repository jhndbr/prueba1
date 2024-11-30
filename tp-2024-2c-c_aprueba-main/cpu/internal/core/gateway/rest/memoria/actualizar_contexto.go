package memoria

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"log/slog"
)

func ActualizarContextoDeEjecucion(ctx *dto.Context) error {
	slog.Info("ACTUALIZO CONTEXTO DE EJECUCION", "TID", ctx.TID)
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/actualizarContexto", conf.IPMemory, conf.PortMemory)
	body, _ := json.Marshal(dto.Context{
		PID:      ctx.PID,
		TID:      ctx.TID,
		Register: ctx.Register,
		PC:       ctx.PC,
		Base:     ctx.Base,
		Limit:    ctx.Limit,
	})

	_, err := rest.Send(url, "PUT", body)
	if err != nil {
		return err
	}
	return nil
}
