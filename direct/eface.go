package direct

import (
	"reflect"
	"unsafe"
)

type rtype struct {
	_    uintptr // size of the type
	_    uintptr // number of bytes in the type that are pointers
	hash uint32  // hash of the type
}

type eface struct {
	rtype *rtype
	data  unsafe.Pointer
}

func typeHash(value any) uint32 {
	return (*eface)(unsafe.Pointer(&value)).rtype.hash
}

func same(a, b any) bool {
	return (*eface)(unsafe.Pointer(&a)).data == (*eface)(unsafe.Pointer(&b)).data
}

func isPointer(value any) bool {
	rv := reflect.ValueOf(value)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}

	return true
}
