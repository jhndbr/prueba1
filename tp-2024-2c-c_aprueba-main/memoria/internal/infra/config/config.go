package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"sync"
)

// Config Estructura que representa la configuración
type Config struct {
	Port            int    `json:"port"`
	MemorySize      int    `json:"memory_size"`
	InstructionPath string `json:"instruction_path"`
	ResponseDelay   int    `json:"response_delay"`
	IPKernel        string `json:"ip_kernel"`
	PortKernel      int    `json:"port_kernel"`
	IPCPU           string `json:"ip_cpu"`
	PortCPU         int    `json:"port_cpu"`
	IPFilesystem    string `json:"ip_filesystem"`
	PortFilesystem  int    `json:"port_filesystem"`
	Scheme          string `json:"scheme"`
	SearchAlgorithm string `json:"search_algorithm"`
	Partitions      []int  `json:"partitions"`
	LogLevel        string `json:"log_level"`
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
			slog.Error("Error al cargar la configuración: %v", err)
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
			configPath = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/memoria/config.json" // "/home/utnso/Documents/tp-2024-2c-c_aprueba/memoria/config.json"
		}
		config, err = createConfig(configPath)
		if err != nil {
			log.Fatalf("Error al cargar la configuración: %v", err)
		}
	})
	return config
}
