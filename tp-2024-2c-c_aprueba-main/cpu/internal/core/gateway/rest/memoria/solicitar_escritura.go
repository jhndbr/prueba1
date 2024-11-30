package memoria

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
)

// type Request struct {
// 	Direction uint32 `json:"direction"`
// 	Value     uint32 `json:"value"`
// }

func EnviarSolicitudEscribirEnMemoria(direccionFisica uint32, valor uint32) error {
	slog.Debug("Enviando solicitud de escritura en memoria")
	conf := config.GetInstance()
	url := fmt.Sprintf("http://%s:%d/memory/escribirMemoria", conf.IPMemory, conf.PortMemory)
	valorEnBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(valorEnBytes, valor)
	// [0005]
	request := struct {
        Direction int  `json:"direction"`
        Value     []byte `json:"value"`
    }{
        Direction: int(direccionFisica),
        Value:     valorEnBytes,
    }

	body, _ := json.Marshal(request)
	_, err := rest.Send(url, "POST", body)
	if err != nil {
		slog.Debug("Error al enviar solicitud de escritura en memoria")
		return err
	}
	slog.Debug("Solicitud de escritura en memoria enviada")
	return nil
}
