package pages

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Scene interface {
	View(ctx PageBuildContext)
}

type SceneWithInit interface {
	Init()
}

type View interface {
	Build(ctx PageBuildContext)
}

type State[T any] struct {
	state             *T
	innerChangeHandle chan<- struct{}
}

func (s State[T]) mount(handle chan<- struct{}, instantiateAddr unsafe.Pointer) {
	var ptr = s.state
	if s.state == nil {
		ptr = new(T)
	}
	*(*State[T])(instantiateAddr) = State[T]{
		state:             ptr,
		innerChangeHandle: handle,
	}
}

type stateMarker interface {
	stateImpl()
	mount(handle chan<- struct{}, instantiateAddr unsafe.Pointer)
}

func (s State[T]) stateImpl() {}

func (s State[T]) Get() T {
	return *s.state
}

func StateOf[T any](v T) State[T] {
	return State[T]{
		state: &v,
	}
}

func (s State[T]) Set(v T) {
	*s.state = v
	s.innerChangeHandle <- struct{}{}
}

func (s State[T]) SetFn(fn func(prev T) T) {
	s.Set(fn(s.Get()))
}

func makePointer(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Pointer {
		return value
	}
	val := value.Interface()
	return reflect.ValueOf(&val)
}

func printBytesOf[T any](v *T, tag string) {
	ptrToVal := unsafe.Pointer(v)
	size := unsafe.Sizeof(*v)
	bytes := make([]byte, size)
	for i := 0; i < int(size); i++ {
		currPtr := unsafe.Add(ptrToVal, uintptr(i))
		bytes[i] = *(*byte)(currPtr)
	}

	fmt.Printf("%s is %x | %x (%d bytes)\n", tag, bytes[0:8], bytes[8:], size)
}

func getPointerToInterfaceValue(v any) unsafe.Pointer {
	ptrToPtr := unsafe.Add(unsafe.Pointer(&v), 8)
	valPtr := *(*uintptr)(ptrToPtr)
	return unsafe.Pointer(valPtr)
}

func MountStates(v *Scene, handleChan chan<- struct{}) {

	rv := reflect.ValueOf(v).Elem()
	if rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	rt := rv.Type()
	ptrToVal := getPointerToInterfaceValue(*v)

	for i := 0; i < rt.NumField(); i++ {
		field := rv.Field(i)
		if rt.Field(i).PkgPath != "" {
			continue
		}
		mountable, isMountable := field.Interface().(stateMarker)

		if !isMountable {

			continue
		}
		mountable.mount(handleChan, unsafe.Add(ptrToVal, rt.Field(i).Offset))
	}
	sceneWithInit, isSceneWithInit := (*v).(SceneWithInit)
	if isSceneWithInit {
		go func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					fmt.Println("Recovered", recovered)
				}
			}()
			sceneWithInit.Init()
		}()
	}
}
