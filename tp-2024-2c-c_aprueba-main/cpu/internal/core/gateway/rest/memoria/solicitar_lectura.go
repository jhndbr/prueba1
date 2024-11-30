package memoria

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"io/ioutil"
)

func EnviarSolicitudLeerValorMemoria(direccionFisica uint32) (*uint32, error) {

	conf := config.GetInstance()

	url := fmt.Sprintf("http://%s:%d/memory/leerMemoria?direccion=%d", conf.IPMemory, conf.PortMemory, direccionFisica)

	response, err := rest.Send(url, "GET", nil)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    byteValue, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }
	
    var result map[string][]byte
    err = json.Unmarshal(byteValue, &result)
    if err != nil {
        return nil, err
    }

    valor, ok := result["valor"]
    if !ok {
        return nil, fmt.Errorf("valor no encontrado en la respuesta")
    }
	if len(valor) != 4 {
        return nil, fmt.Errorf("longitud de valor incorrecta: se esperaban 4 bytes, se recibieron %d", len(valor))
    }
    valorUInt32 := binary.BigEndian.Uint32(valor)

    return &valorUInt32, nil

}
