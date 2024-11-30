package main

import (
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/core/usecase"
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/infra/config"
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/infra/http"
	"log/slog"
)

// var MNT_DIR = "../../MNT_DIR"
var MOUNT_DIR = "MOUNT_DIR"

func main() {
	// Configuro el Log de la Aplicación
	slog.SetLogLoggerLevel(slog.LevelInfo)
	conf := config.GetInstance()
	if conf == nil {
		slog.Error("Error loading config: config is nil")
		return
	}
	iniciar_filesystem(conf)

	//usecase.MarcarBloqueOcupado(3)
	//usecase.MarcarBloqueOcupado(4)
	//usecase.MarcarBloqueOcupado(7)
	//usecase.DesmarcarBloqueOcupado(3)
	//usecase.DesmarcarBloqueOcupado(4)

	usecase.VerContenidoBitmap()
	usecase.VerContenidoArchivoBitmapBinario()
	usecase.VerContenidoArchivoBitmapHexa()

	// funcion para ver si lee bien los enteros del bloque de indice especificado

	//usecase.VerContenidorBloqueIndice(uint32(0))
	//slog.Info("##", "config", conf)
	err := http.Server("127.0.0.1", conf.Port)
	if err != nil {
		slog.Error("Server error", "error", err)
		return
	}
}

func iniciar_filesystem(conf *config.Config) {

	/* verifico que exista la carpeta MNT_DIR y MOUNT_DIR */

	// DIRECTORIO PARA BITMAP.DAT Y BLOQUES.DAT
	err := usecase.Crear_directorio(MOUNT_DIR)
	if err != nil {
		slog.Error("Error creando MOUNT_DIR", "error", err)
	}

	// BITMAP
	if !usecase.ExisteArchivo(MOUNT_DIR + "/bitmap.dat") {
		slog.Info("creando bitmap.dat")
		usecase.Crear_archivo_bitmap(MOUNT_DIR+"/bitmap.dat", conf.Block_count/8)
	} else {
		slog.Info("bitmap.dat ya creado")
		slog.Info("cargando bitmap.dat")
		usecase.Cargar_bitmap(MOUNT_DIR+"/bitmap.dat", conf.Block_count/8)
	}

	// BLOQUES
	if !usecase.ExisteArchivo(MOUNT_DIR + "/bloques.dat") {
		slog.Info("creando bloques.dat")
		usecase.Crear_archivo_bloques(MOUNT_DIR+"/bloques.dat", conf.Block_count*conf.Block_size)
	} else {
		slog.Info("bloques.dat ya creado")
	}

	entity.TamanioBloque = conf.Block_size
	entity.DelayBlock = conf.Block_access_delay

}
