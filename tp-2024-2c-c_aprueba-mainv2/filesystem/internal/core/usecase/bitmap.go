package usecase

import (
	"encoding/hex"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/filesystem/internal/core/entity"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Crear_archivo_bitmap(pathBitmap string, tamanio int) {

	archivo, err := os.OpenFile(pathBitmap, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()

	err = archivo.Truncate(int64(tamanio))
	if err != nil {
		log.Fatal(err)
	}

	// aca quiero inicializar bitmap
	inicializar_bitmap(pathBitmap, tamanio)

}
func inicializar_bitmap(pathBitmap string, tamanio int) {
	archivo, err := os.OpenFile(pathBitmap, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()

	// Inicializa el archivo con ceros
	entity.Bitmap = make([]byte, tamanio)
	_, err = archivo.Write(entity.Bitmap) // Escribe el bitmap en el archivo bitmap.dat
	if err != nil {
		log.Fatal(err)
	}
}
func Cargar_bitmap(path string, tamanio int) {
	archivo, err := os.Open(path) // Abre el archivo bitmap.dat
	if err != nil {
		log.Fatal(err)
	}

	defer archivo.Close() // al finalizar Cargar_bitmap() cierra el archivo

	entity.Bitmap = make([]byte, tamanio)
	_, err = archivo.Read(entity.Bitmap)
	if err != nil {
		log.Fatal(err)
	}
}
func HayEspacioDisponible(bloquesNecesarios int) bool {
	contador := 0
	for i := 0; i < len(entity.Bitmap); i++ {
		if entity.Bitmap[i] == 0 { // si hay bloqe libre aumenta el contador
			contador++
		}
		if contador == bloquesNecesarios {
			return true
		}
	}
	return false
}
func MarcarBloqueOcupado(posicion int) {
	entity.Bitmap[posicion] = 1
	ActualizarBitmapEnArchivo() // para que se actualice el bitmap en el archivo bitmap.dat
}
func DesmarcarBloqueOcupado(posicion int) {
	entity.Bitmap[posicion] = 0
	ActualizarBitmapEnArchivo() // para que se actualice el bitmap en el archivo bitmap.dat
}
func ActualizarBitmapEnArchivo() {
	archivoBitmap, err := os.OpenFile("MOUNT_DIR/bitmap.dat", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer archivoBitmap.Close()

	_, err = archivoBitmap.WriteAt(entity.Bitmap, 0) // Escribe el bitmap en el archivo bitmap.dat desde la posición 0
	if err != nil {
		log.Fatal(err)
	}
}
func VerContenidoArchivoBitmapBinario() {
	path := "MOUNT_DIR/bitmap.dat"

	// Leer el contenido del archivo
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error al leer el archivo %s: %v", path, err)
	}

	// Convertir el contenido a formato binario y agrupar en bloques de 8 bits
	var binContent strings.Builder
	for i, byteVal := range content {
		binContent.WriteString(fmt.Sprintf("%08b", byteVal))
		if (i+1)%1 == 0 && i != len(content)-1 {
			binContent.WriteString(" ") // Espacio entre cada byte
		}
	}

	// Imprimir el contenido en formato binario
	fmt.Printf("Contenido de bitmap.dat en formato binario:\n%s\n", binContent.String())

}
func VerContenidoArchivoBitmapHexa() {
	// Ruta al archivo bitmap.dat
	path := "MOUNT_DIR/bitmap.dat"

	// Leer el contenido del archivo
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error al leer el archivo %s: %v", path, err)
	}

	// Convertir el contenido a formato hexadecimal y separarlo en pares de dígitos
	hexContent := hex.EncodeToString(content)
	var formattedHex strings.Builder
	for i := 0; i < len(hexContent); i += 2 {
		formattedHex.WriteString(hexContent[i:i+2] + " ")
	}

	// Imprimir el contenido en formato hexadecimal
	fmt.Printf("Contenido de bitmap.dat en formato hexadecimal:\n%s\n", strings.TrimSpace(formattedHex.String()))

}
func VerContenidoBitmap() {
	for i, bit := range entity.Bitmap {
		fmt.Printf("Posición %d: %d\n", i, bit)
	}
}
func buscarBloqueLibre() int {
	for i := 0; i < len(entity.Bitmap); i++ {
		if entity.Bitmap[i] == 0 {
			return i
		}
	}
	return -1
}
func cantidadBloquesLibres() int {
	contador := 0
	for i := 0; i < len(entity.Bitmap); i++ {
		if entity.Bitmap[i] == 0 {
			contador++
		}
	}
	return contador
}
