package rest

import (
	"bytes"
	"errors"
	"net/http"
)

// Send function used in order to sent messages http.
func Send(url string, method string, body []byte) (*http.Response, error) {
	// Crear una solicitud HTTP
	var req *http.Request
	var err error

	// Crear la solicitud HTTP
	switch method {

	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
		break
	case http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		break
	case http.MethodPut:
		req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
		break
	case http.MethodDelete:
		req, err = http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))
		break
	default:
		return nil, errors.New("method invalid")
	}

	// Si hubo un error al crear la solicitud, lo manejamos dentro del callback
	if err != nil {
		return nil, err
	}

	// Establecer el tipo de contenido en caso de que haya un cuerpo
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud con un cliente HTTP
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
