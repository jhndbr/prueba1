package memoria

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/core/gateway/rest/dto"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"io/ioutil"
	"log/slog"
	"net/http"
)

func EnviarSolicitudInstruccion(context *dto.Context) (*dto.Instruction, error) {
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/instruction?pid=%d&tid=%d&pc=%d", conf.IPMemory, conf.PortMemory, context.PID, context.TID, context.PC)

	var response *http.Response
	response, err := rest.Send(url, "GET", nil)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		slog.Info("Has not found any instruction")
		return nil, nil
	}

	byteValue, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	var instruction *dto.Instruction
	err = json.Unmarshal(byteValue, &instruction)

	if err != nil {
		slog.Error("Error once has try to deserealize instruction from Memory. Error=", err.Error())
		return nil, err
	}

	return instruction, nil
}
