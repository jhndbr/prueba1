package usecase

import (
	"errors"
	"fmt"
	"log/slog"
	"math"

	"github.com/sisoputnfrba/tp-golang/memoria/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/memoria/internal/infra/config"
)

// esta previsto la fragmentacion interna externa???
func InicializarParticion(tamanios []int) {
	base := 0

	for i, tamPart := range tamanios {
		if base+tamPart > config.GetInstance().MemorySize {
			slog.Debug("Error: la partición %d excede el límite de memoria.\n", i)
			break
		}

		particion := &entity.Particion{
			Base:   base,
			Limite: tamPart,
			Libre:  true,
			PID:    0,
		}
		entity.Particiones = append(entity.Particiones, particion)

		slog.Debug("Partición creada ->", "Particion", i, "base:", particion.Base, "Tamaño:", tamPart)

		// Actualizar la base para la siguiente partición
		base += tamPart
	}
	slog.Debug("Inicializando entity.Particiones")
	slog.Debug(fmt.Sprintf("entity.Particiones: %v", entity.Particiones))
}

func CalcularBase(index int, tamanios []int) int {
	var base int
	for i := 0; i < index; i++ {
		base += tamanios[i]
	}
	return base
}
func InicializarMemoriaDinamica(memorySize int) {
	entity.Particiones = []*entity.Particion{
		&entity.Particion{
			Base:   0,
			Limite: memorySize - 1,
			Libre:  true,
			PID:    0,
		},
	}
}

func AsignarParticion(pid uint32, tam int, scheme string) (int, error) {
	slog.Debug(fmt.Sprintf("Asignando partición %s PID: %d Tamaño: %d", scheme, pid, tam))
	if scheme == "FIJAS" {
		slog.Debug("Asignando partición fija")
		return AsignarParticionFija(pid, tam)
	} else {
		slog.Debug("Asignando partición dinámica")
		return AsignarParticionDinamica(pid, tam)
	}
}

func AsignarParticionFija(pid uint32, tam int) (int, error) {
	for _, p := range entity.Particiones {
		slog.Debug("base menos limite", int(math.Abs(float64(p.Base-p.Limite))))
		if p.Libre && (int(math.Abs(float64(p.Base-p.Limite))) >= tam) || p.Libre && p.Limite == p.Base {
			p.Libre = false
			p.PID = pid
			// VER QUE HACER CON LO DE MEMORIA DE USUARIO
			//for j := p.Base; j < p.Base+tam; j++ {
			//	entity.MemoriaUsuario[j] = 1
			//}

			return p.Base, nil

		}
	}
	return -1, errors.New("No hay particiones fijas disponibles para este tamaño")
}

func LiberarParticion(pid uint32) {

	for i, p := range entity.Particiones {
		if !p.Libre && p.PID == pid {

			p.Libre = true
			p.PID = 0

			// verificamos si la particion anterior es libre
			if i > 0 && entity.Particiones[i-1].Libre {

				entity.Particiones[i-1].Limite = p.Limite
				entity.Particiones = append(entity.Particiones[:i], entity.Particiones[i+1:]...)
				i--
			}

			// verificams si la particion siguiente es libre
			if i < len(entity.Particiones)-1 && entity.Particiones[i+1].Libre {

				entity.Particiones[i].Limite = entity.Particiones[i+1].Limite
				entity.Particiones = append(entity.Particiones[:i+1], entity.Particiones[i+2:]...)
			}
			/*
				for j := p.Base; j < p.Limite; j++ {
					entity.MemoriaUsuario[j] = 0
				}

			*/

			break
		}
	}
}

func AsignarParticionDinamica(pid uint32, tam int) (int, error) {
	var totalLibre int = 0
	for _, p := range entity.Particiones {
		if p.Libre {
			totalLibre += p.Limite - p.Base

		}
	}
	if totalLibre < tam {
		return -1, errors.New("No hay suficiente espacio libre, se requiere compactación")
	}

	switch config.GetInstance().SearchAlgorithm {
	case "FIRST":
		return AsignarParticionFirstFit(pid, tam)
	case "BEST":
		return AsignarParticionBestFit(pid, tam)
	case "WORST":
		return AsignarParticionWorstFit(pid, tam)
	default:
		return -1, errors.New("Algoritmo de búsqueda no soportado")
	}
}

func ProcesoAsociado(pid uint32, base int) bool {
	for _, p := range entity.Particiones {

		if p.PID == pid {
			return true
		}
	}
	return false
}

// Implementación de AsignarParticionFirstFit
func AsignarParticionFirstFit(pid uint32, tam int) (int, error) {
	for i, p := range entity.Particiones {
		if p.Libre && (p.Limite-p.Base) >= tam {
			base := p.Base

			if (p.Limite - p.Base) > tam {
				nuevaParticion := &entity.Particion{

					Base:   base + tam,
					Limite: p.Limite,
					Libre:  true,
				}

				entity.Particiones[i].Limite = base + tam
				entity.Particiones[i].Libre = false
				entity.Particiones[i].PID = pid

				entity.Particiones = append(entity.Particiones[:i+1], append([]*entity.Particion{nuevaParticion}, entity.Particiones[i+1:]...)...)
			} else {
				entity.Particiones[i].Libre = false
				entity.Particiones[i].PID = pid
			}
			/*
				for j := base; j < base+tam; j++ {
					entity.MemoriaUsuario[j] = 1 // O un valor que indique que está ocupado
				}
			*/
			return base, nil
		}
	}
	return -1, errors.New("No hay Particiones disponibles")
}

func AsignarParticionBestFit(pid uint32, tam int) (int, error) {
	var mejorParticion *entity.Particion
	mejorIdx := -1

	for i, p := range entity.Particiones {
		if p.Libre && (p.Limite-p.Base) >= tam {
			if mejorParticion == nil || (p.Limite-p.Base) < (mejorParticion.Limite-mejorParticion.Base) {
				mejorParticion = p

				mejorIdx = i
			}
		}
	}

	if mejorParticion != nil {
		base := mejorParticion.Base

		if (mejorParticion.Limite - mejorParticion.Base) > tam {
			entity.Particiones = append(entity.Particiones[:mejorIdx+1], entity.Particiones[mejorIdx:]...)
			entity.Particiones[mejorIdx+1] = &entity.Particion{Base: base + tam, Limite: mejorParticion.Limite, Libre: true}
			entity.Particiones[mejorIdx].Limite = base + tam
		}
		entity.Particiones[mejorIdx].Libre = false
		entity.Particiones[mejorIdx].PID = pid
		/*
			for j := base; j < base+tam; j++ {
				entity.MemoriaUsuario[j] = 1 // O un valor que indique que está ocupado
			}

		*/
		return entity.Particiones[mejorIdx].Base, nil
	}

	return -1, errors.New("No hay entity.Particiones disponibles")
}

func AsignarParticionWorstFit(pid uint32, tam int) (int, error) {
	var peorParticion *entity.Particion
	peorIdx := -1

	for i, p := range entity.Particiones {
		if p.Libre && (p.Limite-p.Base) >= tam {
			if peorParticion == nil || (p.Limite-p.Base) > (peorParticion.Limite-peorParticion.Base) {
				peorParticion = p

				peorIdx = i
			}
		}
	}

	if peorParticion != nil {
		base := peorParticion.Base

		if (peorParticion.Limite - peorParticion.Base) > tam {
			entity.Particiones = append(entity.Particiones[:peorIdx+1], entity.Particiones[peorIdx:]...)
			entity.Particiones[peorIdx+1] = &entity.Particion{Base: base + tam, Limite: peorParticion.Limite, Libre: true}
			entity.Particiones[peorIdx].Limite = base + tam
		}
		entity.Particiones[peorIdx].Libre = false
		entity.Particiones[peorIdx].PID = pid
		/*
			for j := base; j < base+tam; j++ {
				entity.MemoriaUsuario[j] = 1 // O un valor que indique que está ocupado
			}
		*/
		return entity.Particiones[peorIdx].Base, nil
	}

	return -1, errors.New("No hay entity.Particiones disponibles")
}

//	func Compactar() {
//		var nuevaBase int = 0
//		var particionesOcupadas []*entity.Particion
//		var espacioLibreTotal int = 0
//
//		for _, particion := range entity.Particiones {
//			if particion.Libre {
//				espacioLibreTotal += particion.Limite - particion.Base
//			} else {
//				particionesOcupadas = append(particionesOcupadas, particion)
//			}
//		}
//
//		for _, particion := range particionesOcupadas {
//			tamano := particion.Limite - particion.Base
//
//			particion.Base = nuevaBase
//			particion.Limite = nuevaBase + tamano
//			nuevaBase += tamano
//		}
//
//		particionLibre := &entity.Particion{
//			Base:   nuevaBase,
//			Limite: nuevaBase + espacioLibreTotal,
//			Libre:  true,
//			PID:    0,
//		}
//
//		entity.Particiones = append(particionesOcupadas, particionLibre)
//
//		slog.Info("Memoria compactada y particiones libres consolidadas.")
//		for i, p := range entity.Particiones {
//			fmt.Printf("Partición %d: Base=%d, Limite=%d, Libre=%v, PID=%d\n", i, p.Base, p.Limite, p.Libre, p.PID)
//		}
//	}
func Compactar() {
	var nuevaBase int = 0
	var particionesOcupadas []*entity.Particion
	var espacioLibreTotal int = 0

	memoriaTemporal := make([]byte, len(entity.MemoriaUsuario))

	for _, particion := range entity.Particiones {
		if particion.Libre {
			espacioLibreTotal += particion.Limite - particion.Base
		} else {
			particionesOcupadas = append(particionesOcupadas, particion)
		}
	}

	for _, particion := range particionesOcupadas {
		tamano := particion.Limite - particion.Base

		copy(memoriaTemporal[nuevaBase:nuevaBase+tamano], entity.MemoriaUsuario[particion.Base:particion.Limite])

		particion.Base = nuevaBase
		particion.Limite = nuevaBase + tamano

		//para actualizar no se si es necesario
		for i, contexto := range entity.MemoriaSistema {
			if contexto.PID == particion.PID {
				entity.MemoriaSistema[i].Base = uint32(particion.Base)
				entity.MemoriaSistema[i].Limite = uint32(particion.Limite)
			}
		}

		nuevaBase += tamano
	}
	copy(entity.MemoriaUsuario, memoriaTemporal)

	particionLibre := &entity.Particion{
		Base:   nuevaBase,
		Limite: nuevaBase + espacioLibreTotal,
		Libre:  true,
		PID:    0,
	}

	entity.Particiones = append(particionesOcupadas, particionLibre)

	slog.Debug("Memoria compactada y particiones libres consolidadas.")

}
