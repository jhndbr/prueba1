package config

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"os"
	"sync"
)

// Config Estructura que representa la configuración
type Config struct {
	Port       int    `json:"port"`
	IPMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	IPKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	LogLevel   string `json:"log_level"`
}

var config *Config
var once sync.Once

// Función para leer y parsear el archivo de configuración
func createConfig(filename string) (*Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			slog.Error("Error al cargar la configuración: %v", err.Error())
		}
	}(configFile)

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetInstance() *Config {
	once.Do(func() {
		var err error
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/cpu/config.json"
		}
		config, err = createConfig(configPath)
		if err != nil {
			slog.Error("Error al crear la instancia de configuración", err.Error())
		}
	})
	return config
}
