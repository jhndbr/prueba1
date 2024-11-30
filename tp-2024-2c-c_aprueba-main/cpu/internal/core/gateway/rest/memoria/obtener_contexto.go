package memoria

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"io/ioutil"
	"log/slog"
)

// ObtenerContextoDeEjecucion recupera el contexto de un proceso
// Ejemplo: GET http://localhost:8003/contexto?pid=1&tid=1
func ObtenerContextoDeEjecucion(PID uint32, TID uint32) (*dto.Context, error) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/context?pid=%d&tid=%d", conf.IPMemory, conf.PortMemory, PID, TID)

	response, err := rest.Send(url, "GET", nil)
	if err != nil {
		slog.Error("Error al enviar contexto: ", "Error=", err.Error())
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		slog.Error("Error en deserealizar contexto: ", "Error=", err.Error())
		return nil, err
	}
	var ctx *dto.Context
	err = json.Unmarshal(byteValue, &ctx)
	if err != nil {
		slog.Error("Error en biding de contexto: ", "Error=", err.Error())
		return nil, err
	}

	return ctx, nil
}
