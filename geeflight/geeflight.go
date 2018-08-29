package geeflight

import (
	"errors"
	"fmt"
	"reflect"
)

func resultHandlerExists(handler ...interface{}) (error, bool, reflect.Value) {
	hLen := len(handler)
	if hLen >= 1 {
		h := handler[0]
		if reflect.TypeOf(h).Kind() != reflect.Func {
			return errors.New("the final handler is not a function"), false, reflect.Value{}
		} else {
			return nil, true, reflect.ValueOf(h)
		}
	}
	return nil, false, reflect.Value{}
}

func initParams(funcType reflect.Type, reflectedParams []reflect.Value) {
	for idx := 0; idx < len(reflectedParams); idx++ {
		expectedType := funcType.In(idx)
		reflectedParams[idx] = reflect.New(expectedType).Elem()
	}
}

func makeParams(funcType reflect.Type, input []reflect.Value, outCount int) []reflect.Value {
	params := make([]reflect.Value, outCount)
	initParams(funcType, params)
	retLen := len(input)
	if retLen <= outCount {
		copy(params, input)
	} else {
		copy(params, input[0:outCount])
	}
	return params
}

func Waterfall(fa []interface{}, resultHandler ...interface{})  {
	err, ok, callBack := resultHandlerExists(resultHandler...)
	if err != nil {
		panic( err.Error())
	}

	var cv reflect.Type
	if ok {
		cv = reflect.TypeOf(resultHandler[0])
	}
	faLen := len(fa)
	var fi []reflect.Value
	for idx, p := range fa {
		if reflect.TypeOf(p).Kind() != reflect.Func {
			panic(fmt.Sprintf("One of the element(%d) is not a function", idx))
		}
		v := reflect.TypeOf(p)
		params := makeParams(v, fi, v.NumIn())

		val := reflect.ValueOf(p)
		r := val.Call(params)
		retLen := len(r)
		if idx < faLen - 1 {
			if retLen > 0 {

				ret1 := r[retLen-1].Interface()
				if ret1 != nil {
					// fmt.Printf("Error calling %d: %v\n", idx, ret1)
					if ok {
						params := makeParams(cv, r, cv.NumIn())
						params[cv.NumIn()-1] = r[retLen-1]
						callBack.Call(params)
						return
					}
				}
				fi = r[0:retLen-1]
			}else{
				panic(fmt.Sprintf("the empty return is from %d function", idx))
			}
		} else {
			if ok {
				params := makeParams(cv, r[0:retLen-1], cv.NumIn())
				params[cv.NumIn()-1] = r[retLen-1]
				callBack.Call(params)
				return
			}
		}
	}
	return
}

func convertInterfaceArrayToValueArray(ia []interface{}) []reflect.Value {
	var rva []reflect.Value = make([]reflect.Value, len(ia))
	for i := range ia {
		rva[i] = reflect.ValueOf(rva[i])
	}
	return rva
}

const ErrorMsgFuncPointerExpected string = "the first argument should be a function pointer"
const ErrorMsgFuncReturnExpected string = "there should be at one return value from the function"
const ErrorMsgWrongArgNumber string = "the number of Guard argument expects to be executed function argument + function return - error"

func Guard(args ...interface{}) {
	ret, withErr := exeGuard(args...)
	if withErr {
		panic(ret)
	}
}

func CGuard(args ...interface{}) {
	l := len(args)
	if l < 2 {
		return
	}

	oo := args[0]
	ret, withErr := exeGuard(args[1:]...)
	fmt.Printf("ret: %v, withErr: %v\n", ret, withErr)
	if withErr {
		panic([...]interface{}{oo, ret})
	}
}

func exeGuard(args ...interface{}) (interface{}, bool) {
	l := len(args)
	if l < 1 {
		return nil, false
	}


	fp := args[0]
	at := reflect.TypeOf(fp)
	if at.Kind() != reflect.Func {
		return ErrorMsgFuncPointerExpected, true
	}

	if at.NumOut() < 1 {
		return ErrorMsgFuncReturnExpected, true
	}

	if l != (at.NumIn() + at.NumOut()) {
		return ErrorMsgWrongArgNumber, true
	}

	var ret []reflect.Value
	if at.NumIn() > 0 {
		ret = reflect.ValueOf(fp).Call(makeParams(at, convertInterfaceArrayToValueArray(args[1:at.NumIn()]), at.NumIn()))
	}else{
		ret = reflect.ValueOf(fp).Call(nil)
	}

	fr := ret[at.NumOut()-1]
	if fr.Interface() != nil {
		return fr.Interface(), true
	}

	for i, out := range args[at.NumIn()+ 1: at.NumOut()]{
		reflect.ValueOf(out).Elem().Set(ret[i])
	}
	return nil, false
}

func CatchGuard(handler func(error)) {
	r := recover()

	if r == nil {
		return
	}

	if IsError(r) {
		handler(r.(error))
		return
	}
	panic(r)
}

func CatchCGuard(handler func(interface{}, error)) {
	r := recover()

	if r == nil {
		return
	}

	if reflect.TypeOf(r).Kind() != reflect.Array && len(r.([]interface{})) != 2{
		panic(r)
	}

	ra := r.([2]interface{})
	if IsError(ra[1]) {
		handler(ra[0], ra[1].(error))
		return
	}
	panic(r)
}

func CatchIntCGuard(handler func(int, error)) {
	r := recover()

	if r == nil {
		return
	}

	if reflect.TypeOf(r).Kind() != reflect.Array && len(r.([]interface{})) != 2{
		panic(r)
	}

	ra := r.([2]interface{})
	if IsInt(ra[0]) && IsError(ra[1]) {
		handler(ra[0].(int), ra[1].(error))
		return
	}
	panic(r)
}

func IsSameType(v interface{}, myType interface{}) bool {
	return reflect.TypeOf(v) == reflect.TypeOf(myType)
}

func IsError(e interface{}) bool {
	_, ok := (e).(error)
	return ok
}

func IsInt(e interface{}) bool {
	_, ok := (e).(int)
	return ok
}