package list

import "fmt"

// List Definir la interfaz List
type List[T any] interface {
	Add(item T)                                          // Add Añadir un elemento
	Remove(index int)                                    // Remove Eliminar un elemento en el índice dado
	Get(index int) (T, error)                            // Get Obtener un elemento en el índice dado
	Size() int                                           // Size Retornar el tamaño de la lista
	Filter(value T, predicate func(a, b T) bool) List[T] // Filter Filtra elementos de la lista
	Sort(less func(a, b T) bool)                         // Sort Ordena una Lista de acuerdo al criterio
}

// ArrayList implements List
type ArrayList[T any] struct {
	items []T
}

// Add item to a List
func (list *ArrayList[T]) Add(item T) {
	list.items = append(list.items, item)
}

// Remove item from a List
func (list *ArrayList[T]) Remove(index int) {
	if index >= 0 && index < len(list.items) {
		list.items = append(list.items[:index], list.items[index+1:]...)
	}
}

// Get item from a List
func (list *ArrayList[T]) Get(index int) (T, error) {
	// Validar si el índice está dentro del rango
	if index < 0 || index >= len(list.items) {
		var zero T // Crear un valor cero del tipo genérico T
		return zero, fmt.Errorf("index out of range: %d", index)
	}
	return list.items[index], nil
}

// Size of the list
func (list *ArrayList[T]) Size() int {
	return len(list.items)
}

// Filter items from the list
func (list *ArrayList[T]) Filter(value T, predicate func(a, b T) bool) List[T] {
	filteredList := &ArrayList[T]{}

	for _, item := range list.items {
		if predicate(value, item) {
			filteredList.Add(item)
		}
	}
	return filteredList
}

// Sort ordena una lista
func (list *ArrayList[T]) Sort(less func(a, b T) bool) {
	size := list.Size()
	for i := 0; i < size-1; i++ {
		for j := 0; j < size-i-1; j++ {
			if !less(list.items[j], list.items[j+1]) {
				// Intercambiar los elementos si están en el orden incorrecto
				list.items[j], list.items[j+1] = list.items[j+1], list.items[j]
			}
		}
	}
}

func (list *ArrayList[T]) RemoveElement(element *T) {

}
