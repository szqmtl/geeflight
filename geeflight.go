package main

import (
	"reflect"
	"errors"
	"fmt"
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
		if idx != faLen {
			if retLen > 0 {

				ret1 := r[retLen-1].Interface()
				if ret1 != nil {
					// fmt.Printf("Error calling %d: %v\n", idx, ret1)
					if ok {
						params := makeParams(cv, fi, cv.NumIn())
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
