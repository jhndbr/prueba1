package test

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/usecase"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
	"testing"
)

// TestInicializarParticion prueba que la función crea correctamente las particiones
func TestInicializarParticion(t *testing.T) {
	configInstance := config.GetInstance()

	if configInstance == nil {
		t.Fatal("La configuración no está inicializada.")
	}

	configInstance.MemorySize = 1024
	configInstance.Partitions = []int{512, 16, 32, 16, 256, 64, 128}

	if configInstance.Partitions == nil {
		t.Fatal("Las particiones no están inicializadas.")
	}

	entity.Particiones = []*entity.Particion{}

	tamanios := []int{512, 16, 32, 16, 256, 64, 128}

	usecase.InicializarParticion(tamanios)

	if len(entity.Particiones) != len(tamanios) {
		t.Errorf("Número de particiones incorrecto: esperado %d, obtenido %d", len(tamanios), len(entity.Particiones))
	}

	// Verificamos la base y límite de cada partición creada
	for i, particion := range entity.Particiones {
		expectedBase := usecase.CalcularBase(i, configInstance.Partitions) // Base esperada
		if particion.Base != expectedBase {
			t.Errorf("Base incorrecta en la partición %d: esperado %d, obtenido %d", i, expectedBase, particion.Base)
		}
		if particion.Limite != tamanios[i] {
			t.Errorf("Tamaño incorrecto en la partición %d: esperado %d, obtenido %d", i, tamanios[i], particion.Limite)
		}
		if !particion.Libre {
			t.Errorf("La partición %d debería estar libre, pero no lo está", i)
		}
	}

	pid := uint32(1)
	tam := 32

	_, err := usecase.AsignarParticionFija(pid, tam)
	if err != nil {
		t.Fatalf("Error al asignar partición fija: %v", err)
	}

	particionAsignada := false
	for _, particion := range entity.Particiones {
		if particion.PID == pid {
			particionAsignada = true
			//Verificamos que la partición esté ocupada correctamente
			for i := particion.Base; i < particion.Base+tam; i++ {
				if entity.MemoriaUsuario[i] != 1 {
					t.Errorf("La memoria en el índice %d no está ocupada", i)
				}
			}

			break
		}
	}

	if !particionAsignada {
		t.Errorf("No se encontró la partición con el PID %d asignado", pid)
	}
}
func TestAsignarParticionFirstFit(t *testing.T) {

	memorySize := 512
	entity.Particiones = usecase.InicializarMemoriaDinamica(memorySize)
	//entity.MemoriaUsuario = make([]int, memorySize)

	pid := uint32(1)
	tam := 30
	_, err := usecase.AsignarParticionFirstFit(pid, tam)
	if err != nil {
		t.Fatalf("Error al asignar partición: %v", err)
	}

	_, _ = usecase.AsignarParticionFirstFit(2, 50)
	_, _ = usecase.AsignarParticionFirstFit(3, 50)
	_, _ = usecase.AsignarParticionFirstFit(4, 70)
	_, _ = usecase.AsignarParticionFirstFit(5, 20)
	_, _ = usecase.AsignarParticionFirstFit(6, 10)
	_, _ = usecase.AsignarParticionFirstFit(7, 100)
	_, _ = usecase.AsignarParticionFirstFit(8, 200)
	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------\n")

	usecase.LiberarParticion(2)

	usecase.LiberarParticion(4)

	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------\n")

	_, _ = usecase.AsignarParticionWorstFit(6, 70)
	_, _ = usecase.AsignarParticionWorstFit(1000, 100)
	_, _ = usecase.AsignarParticionWorstFit(2000, 10)

	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------\n")

	usecase.LiberarParticion(5)
	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------\n")
	usecase.LiberarParticion(6)
	usecase.LiberarParticion(2000)
	usecase.LiberarParticion(6)
	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------\n")

	usecase.Compactar()
	for i, p := range entity.Particiones {
		fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------\n")

}
