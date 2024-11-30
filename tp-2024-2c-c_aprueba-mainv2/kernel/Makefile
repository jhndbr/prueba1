# Nombre del binario resultante
BINARY_NAME=app
BUILD_DIR=build

# Configuración de Go
GO_FILES=$(shell find . -type f -name '*.go')
GO_TEST_FILES=$(shell find . -type f -name '*_test.go')
GO_PACKAGES=$(shell go list ./...)

# Variables para pruebas y compilación cruzada
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Comandos de Make
.PHONY: all build run clean test fmt vet lint docker

# Compilar el binario
build: $(GO_FILES)
	@echo "Compilando el proyecto..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/app

# Ejecutar el proyecto
run: build
	@echo "Ejecutando el proyecto..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Ejecutar los tests
test:
	@echo "Ejecutando pruebas..."
	go test ./... -v

# Formatear el código
fmt:
	@echo "Formateando el código..."
	go fmt ./...

# Analizar el código con go vet
vet:
	@echo "Analizando el código..."
	go vet ./...

# Analizar el código con golint (necesitas tener golint instalado)
lint:
	@echo "Chequeando el código con golint..."
	golint ./...

# Limpiar archivos de compilación
clean:
	@echo "Limpiando archivos generados..."
	rm -rf $(BUILD_DIR)

# Compilación cruzada
build-cross:
	@echo "Compilando para OS=$(GOOS) y ARCH=$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME)_$(GOOS)_$(GOARCH) ./cmd/app

# Crear imagen de Docker
docker:
	@echo "Construyendo imagen de Docker..."
	docker build -t your_project_image .

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  build           - Compila el binario"
	@echo "  run             - Ejecuta el binario compilado"
	@echo "  test            - Ejecuta los tests"
	@echo "  fmt             - Formatea el código con go fmt"
	@echo "  vet             - Analiza el código con go vet"
	@echo "  lint            - Analiza el código con golint"
	@echo "  clean           - Elimina los archivos compilados"
	@echo "  build-cross     - Compila el binario para otros OS/Arquitecturas"
	@echo "  docker          - Construye una imagen Docker"
