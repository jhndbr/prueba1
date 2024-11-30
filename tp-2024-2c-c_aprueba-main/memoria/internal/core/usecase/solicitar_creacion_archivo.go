package usecase

/*
import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/utils/infra/rest"
	"log/slog"
	"net/http"
)

func solicitarCreacionDeArchivo(nombre string, size uint64) {
	// Cargar la configuración del servidor de Filesystem
	conf := config.GetInstance()
	// Crear la estructura del archivo que se va a solicitar
	archivo := map[string]interface{}{
		"nombre": nombre,
		"size":   size,
	}

	// Convertir la estructura a JSON
	body, err := json.Marshal(archivo)
	if err != nil {
		slog.Error("Error al codificar el archivo en JSON: ", err)
		return
	}

	// Construir la URL del servidor de Filesystem
	url := fmt.Sprintf("http://%s:%d/fs/crearArchivo", conf.IpFilesystem, conf.PortFilesystem)
	// Enviar la solicitud POST a Filesystem
	status, err := rest.Send(url, body)
	if err != nil {
		slog.Error("No se pudo solicitar la creación del archivo en Filesystem: ", err)
		return
	}

	if status == fmt.Sprint(http.StatusOK) {
		slog.Info("Archivo creado exitosamente en Filesystem")
	} else {
		slog.Error("Error en la creación del archivo. Estado:", status)
	}

}
*/
