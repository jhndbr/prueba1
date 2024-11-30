package usecase

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/core/entity"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func Crear_directorio(pathDirectorio string) error {

	if _, err := os.Stat(pathDirectorio); os.IsNotExist(err) { // Si no existe el directorio lo crea
		err := os.Mkdir(pathDirectorio, 0777) // creamos directorio (0777) es permiso de lectra escritura y ejecucioon
		if err != nil {
			return fmt.Errorf("error al crear el directorio %s: %w", pathDirectorio, err)
		}
	}
	return nil
}

func Crear_archivo(pid uint32, tid uint32, timestamp string, tamanio int, contenido []byte) int{

	nombreArchivo := strconv.Itoa(int(pid)) + "_" + strconv.Itoa(int(tid)) + "_" + timestamp
	err := Crear_directorio("MOUNT_DIR/FILES")
	if err != nil {
		return -1
	}

	bloquesNecesarios := (tamanio + entity.TamanioBloque - 1) / entity.TamanioBloque
	fmt.Println("Bloques Necesarios:", bloquesNecesarios)

	if HayEspacioDisponible(bloquesNecesarios + 1) {

		fd_archivoBloques, errBloque := os.OpenFile("MOUNT_DIR/bloques.dat", os.O_RDWR, 0666) // abre el archivo de bloques
		if errBloque != nil {
			slog.Error("Error al abrir el archivo de bloques", "Error", errBloque)
		}
		defer fd_archivoBloques.Close()

		posicionBloqueIndice := buscarBloqueLibre() // busca el primer bloque libre

		if posicionBloqueIndice == -1 {
			slog.Error("No hay bloques libres")
			return -1
		} else {

			MarcarBloqueOcupado(posicionBloqueIndice)
			//VerContenidoArchivoBitmapHexa()
			slog.Info("##", "Bloque INDICE asignado:", posicionBloqueIndice, "Archivo:", nombreArchivo, "Bloques Libres:", cantidadBloquesLibres())
		}

		bloquesAsignados := make([]int, bloquesNecesarios)

		// recorro la cantidad de bloques necesarios para guardarlos en el array bloquesAsignados
		for i := 0; i < bloquesNecesarios; i++ {

			// busco el primer bloque libre
			posicionBloqueDatos := buscarBloqueLibre()

			MarcarBloqueOcupado(posicionBloqueDatos)

			bloquesAsignados[i] = posicionBloqueDatos

			slog.Info("##", "Bloque DATOS asignado:", posicionBloqueDatos, "Archivo:", nombreArchivo, "Bloques Libres:", cantidadBloquesLibres())

		}

		fmt.Println("bitmap despues de asignar bloques")
		VerContenidoArchivoBitmapHexa()

		crearArchivoMetadata(nombreArchivo, posicionBloqueIndice, tamanio)

		// escribo en el bloque indice los "punteros" (enteros) a los bloques que contienen los datos
		escribirEnBloqueIndice(fd_archivoBloques, posicionBloqueIndice, bloquesAsignados, nombreArchivo)

		// escribo el contenido en los bloques asignados
		escribirEnBloquesDatos(fd_archivoBloques, contenido, bloquesAsignados, nombreArchivo)

		slog.Info("##", "Fin de solicitud - Archivo:", nombreArchivo)
		return 1
	} else {
		// MANDAR A MEMORIA ERROR
		slog.Info("NO HAY ESPACIO DISPONIBLE")
		slog.Info("##", "Fin de solicitud - Archivo:", nombreArchivo)
		return 0
	}


}

func ExisteArchivo(path string) bool {
	_, err := os.Stat(path)    // verifica si existe el archivo, si no existe devuelve error
	return !os.IsNotExist(err) // si no existe devuelve true
}

func crearArchivoMetadata(nombreArchivo string, pos int, tamanio int) {

	nombreArchivo = strings.ReplaceAll(nombreArchivo, ":", "-") // Cambia ":" por "-"

	pathArchivoMetadata := "MOUNT_DIR/FILES/" + nombreArchivo + ".dmp"

	// Crea el archivo si no existe, si existe lo abre
	archivo, err := os.OpenFile(pathArchivoMetadata, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		slog.Error("Error al crear el archivo", "Error", err)
		return
	}
	defer archivo.Close()

	// escribe en texto el index_block y el tamaño del archivo
	_, err = archivo.WriteString("index_block:" + strconv.Itoa(pos) + "\n" + "size:" + strconv.Itoa(tamanio))

	//slog.Info("##", "Archivo Creado:", nombreArchivo, "Tamaño:", tamanio)
}
