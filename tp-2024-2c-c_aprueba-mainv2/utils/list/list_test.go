package list

//
//import (
//	"testing"
//)
//
//// Test para el método Add
//func TestAdd(t *testing.T) {
//	list := &ArrayList[int]{}
//
//	list.Add(10)
//	list.Add(20)
//
//	if list.Size() != 2 {
//		t.Errorf("Expected size 2, got %d", list.Size())
//	}
//
//	if list.Get(0) != 10 {
//		t.Errorf("Expected 10 at index 0, got %d", list.Get(0))
//	}
//
//	if list.Get(1) != 20 {
//		t.Errorf("Expected 20 at index 1, got %d", list.Get(1))
//	}
//}
//
//// Test para el método Remove
//func TestRemove(t *testing.T) {
//	list := &ArrayList[int]{}
//
//	list.Add(10)
//	list.Add(20)
//	list.Add(30)
//
//	list.Remove(1) // Eliminar el elemento en índice 1
//
//	if list.Size() != 2 {
//		t.Errorf("Expected size 2, got %d", list.Size())
//	}
//
//	if list.Get(1) != 30 {
//		t.Errorf("Expected 30 at index 1, got %d", list.Get(1))
//	}
//}
//
//// Test para el método Get
//func TestGet(t *testing.T) {
//	list := &ArrayList[int]{}
//
//	list.Add(10)
//	list.Add(20)
//
//	if list.Get(0) != 10 {
//		t.Errorf("Expected 10 at index 0, got %d", list.Get(0))
//	}
//
//	if list.Get(1) != 20 {
//		t.Errorf("Expected 20 at index 1, got %d", list.Get(1))
//	}
//}
//
//// Test para el método Size
//func TestSize(t *testing.T) {
//	list := &ArrayList[int]{}
//
//	if list.Size() != 0 {
//		t.Errorf("Expected size 0, got %d", list.Size())
//	}
//
//	list.Add(10)
//
//	if list.Size() != 1 {
//		t.Errorf("Expected size 1, got %d", list.Size())
//	}
//}
//
//// Test para el método Filter
//func TestFilter(t *testing.T) {
//	list := &ArrayList[int]{}
//
//	list.Add(10)
//	list.Add(20)
//	list.Add(30)
//	list.Add(40)
//
//	filtered := list.Filter(func(n int) bool {
//		return n > 20
//	})
//
//	if filtered.Size() != 2 {
//		t.Errorf("Expected size 2, got %d", filtered.Size())
//	}
//
//	if filtered.Get(0) != 30 {
//		t.Errorf("Expected 30 at index 0, got %d", filtered.Get(0))
//	}
//
//	if filtered.Get(1) != 40 {
//		t.Errorf("Expected 40 at index 1, got %d", filtered.Get(1))
//	}
//}
//
