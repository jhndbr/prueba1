package usecase

import (
	"encoding/binary"
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/core/entity"
	"log"
	"log/slog"
	"os"
	"time"
)

func Crear_archivo_bloques(pathBloques string, tamanio int) {

	archivo, err := os.OpenFile(pathBloques, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()

	err = archivo.Truncate(int64(tamanio)) // setea el tamaño del archivo
	if err != nil {
		log.Fatal(err)
	}
}

func escribirEnBloqueIndice(fd_archivoBloques *os.File, posicionBloqueIndice int, bloquesAsignados []int, nombre_archivo string) {
	// [0] [1,2,3]
	offset := int64(posicionBloqueIndice * entity.TamanioBloque)

	_, errOffset := fd_archivoBloques.Seek(offset, 0)

	if errOffset != nil {
		slog.Error("Error al buscar el offset", "Error", errOffset)
		return
	}

	bytesAEscribir := make([]byte, 4*len(bloquesAsignados))

	for i, bloque := range bloquesAsignados {
		// empieza a escribir desde el indice i*4 es decir si i = 0 -> empieza a escribir desde la posicion bytesAEscribir[0] hasta bytesAEscribir[3] = 4 y asi sucesivamente
		binary.LittleEndian.PutUint32(bytesAEscribir[i*4:], uint32(bloque))

	}
	// al final te queda un array de bytes de 4 bytes * el tamaño de los bloquesAsignados
	// por ejemplo se tiene 2 bloques asignados [1,2] => bytesAEscribir = [0,0,0,1,0,0,0,2]
	_, err := fd_archivoBloques.Write(bytesAEscribir)
	if err != nil {
		slog.Error("Error al escribir en el archivo", "Error", err)
		return
	}

	slog.Info("##", "Acceso Bloque - Archivo:", nombre_archivo, "Tipo Bloque:", "INDICE", "Bloque File System:", posicionBloqueIndice)

	// hay que esperar el tiempo delayBlock en milisegundos ante cada acceso a bloqe.dat
	time.Sleep(time.Duration(entity.DelayBlock) * time.Millisecond)

}

func escribirEnBloquesDatos(fd *os.File, contenido []byte, bloquesDatosAsignados []int, nombre_archivo string) {

	// divide el contenido en subarrays del tamaño  block_size
	//4 bytes
	//[12,23,64,25,84,35]
	//[[12,23,64,25], [84,35]]
	contenidoSubArrays := DividirContenido(contenido)
	i := 0
	// [4,5]
	// [ [0,1,2,3,4,5,6,7] , [8,9,10,11,12,13,14,15] , [16,17,18,19,20,21] ]
	for _, bloque := range bloquesDatosAsignados {

		offset := int64(bloque * entity.TamanioBloque)
		_, err := fd.Seek(offset, 0)
		if err != nil {
			slog.Error("Error al buscar el offset", "Error", err)
			return
		}

		_, err = fd.Write(contenidoSubArrays[i])
		if err != nil {
			slog.Error("Error al escribir los DATOS en el archivo", "Error", err)
			return
		}
		i++
		slog.Info("##", "Acceso Bloque - Archivo:", nombre_archivo, "Tipo Bloque:", "DATO", "Bloque File System:", bloque)
		time.Sleep(time.Duration(entity.DelayBlock) * time.Millisecond)
	}

}
func DividirContenido(contenido []byte) [][]byte {
	/*
		tamañoBloque = 8
		i = 0
		contenido = [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21] bytes
		append(arrayDeSubArrays, contenido[0:8]) // devuelve desde la posicion 0 hasta la 7 ya que el 8 no lo incluye
		arrayDeSubArrays = [[0,1,2,3,4,5,6,7]]
		i = 8
		append(arrayDeSubArrays, contenido[8:16])
		arrayDeSubArrays = [[0,1,2,3,4,5,6,7], [8,9,10,11,12,13,14,15]]
		i = 16
		16+8 < 21? NO -> else -> rrayDeSubArrays = append(arrayDeSubArrays, contenido[i:])
		append(arrayDeSubArrays, contenido[16:] // toma a partir de la posicion 16 hasta el final
		arrayDeSubArrays = [[0,1,2,3,4,5,6,7], [8,9,10,11,12,13,14,15], [16,17,18,19,20,21]]
		arrayDeSubArrays tamaño = [ [8] [8] [6] ]
	*/
	arrayDeSubArrays := make([][]byte, 0)
	for i := 0; i < len(contenido); i += entity.TamanioBloque {
		if i+entity.TamanioBloque < len(contenido) {
			arrayDeSubArrays = append(arrayDeSubArrays, contenido[i:i+entity.TamanioBloque])
		} else {
			arrayDeSubArrays = append(arrayDeSubArrays, contenido[i:])
		}
	}
	return arrayDeSubArrays
}
func VerContenidorBloqueIndice(posicionBloqueIndice uint32) {
	fd_archivoBloques, errBloque := os.OpenFile("MOUNT_DIR/bloques.dat", os.O_RDWR, 0666) // abre el archivo de bloques
	if errBloque != nil {
		slog.Error("Error al abrir el archivo de bloques", "Error", errBloque)
	}
	defer fd_archivoBloques.Close()

	offset := int64(posicionBloqueIndice) * int64(entity.TamanioBloque)

	_, errOffset := fd_archivoBloques.Seek(offset, 0)

	if errOffset != nil {
		slog.Error("Error al buscar el offset", "Error", errOffset)
		return
	}

	bytesLeidos := make([]byte, entity.TamanioBloque)
	slog.Info("contenido de bytesLeidos con make([]byte, entity.TamanioBloque)", "bytesLeidos", bytesLeidos)
	// lee los bytes del archivo y los guarda en bytesLeidos, lee cada 4 bytes
	_, err := fd_archivoBloques.Read(bytesLeidos)
	if err != nil {
		slog.Error("Error al leer el archivo", "Error", err)
		return
	}
	slog.Info("contenido de bytesLeidos luego del fd.Read() ", "bytesLeidos", bytesLeidos)

	for i := 0; i < len(bytesLeidos); i += 4 {
		// convierte los 4 bytes a un entero
		numeroBloque := binary.LittleEndian.Uint32(bytesLeidos[i : i+4])
		if numeroBloque != 0 {
			slog.Info("PUNTEROS DEL BLOQUE INDICE:", "Bloque n°:", numeroBloque)
		}
	}
}
