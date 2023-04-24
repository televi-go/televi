package magic

import (
	"reflect"
	"unsafe"
)

func getRealPointer(v any) unsafe.Pointer {
	ptrToPtr := unsafe.Add(unsafe.Pointer(&v), 8)
	return *(*unsafe.Pointer)(ptrToPtr)
}

// InjectInPlace injects state callbacks directly into State fields
func InjectInPlace(v any, c func()) (hasInjected bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	rt := rv.Type()
	ptr := getRealPointer(v)
	for i := 0; i < rt.NumField(); i++ {
		field := rv.Field(i)
		if rt.Field(i).PkgPath != "" {
			continue
		}
		mountable, isMountable := field.Interface().(Mountable)

		if !isMountable {

			continue
		}
		hasInjected = true
		mountable.Mount(c, unsafe.Add(ptr, rt.Field(i).Offset))
	}
	return
}
