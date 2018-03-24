package main

import (
	"testing"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWaterfall(t *testing.T) {
	// assert := assertions.New(t)
	ErrorDefault := "DEFAULT"
	ErrorMsg := "IN ERROR"
	var finalError = errors.New(ErrorDefault)

	var initFuncR1 = func() (int, error) { return 1, nil }
	var initFuncR2 = func() (int, int, error) { return 1, 2, nil }
	// var initFuncR3 = func() (int, int, int, error) { return 1, 2, 3, nil }

	// var funcP1 = func(i int) (error) { return errors.New(ErrorMsg) }
	// var funcP1R1 = func(i int) (int, error) { return i, errors.New(ErrorMsg) }
	// var funcP1R2 = func(i int) (int, int, error) { return i, i * 2, errors.New(ErrorMsg) }
	var funcP1R3 = func(i int) (int, int, int, error) { return i, i * 2, i * 3, errors.New(ErrorMsg) }

	// var funcP2 = func(i, j int) (error) { return nil }
	// var funcP2R1 = func(i, j int) (int, error) { return i * 2, nil }
	// var funcP2R2 = func(i, j int) (int, int, error) { return i, j * 2, nil }
	var funcP2R3 = func(i, j int) (int, int, int, error) { return i, j * 2, j * 3, nil }

	var funcP3 = func(i, j, k int) (error) { return nil }
	var funcP3R1 = func(i, j, k int) (int, error) { return i + j + k, nil }
	// var funcP3R2 = func(i, j, k int) (int, int, error) { return i, j + k, nil }
	// var funcP3R3 = func(i, j, k int) (int, int, int, error) { return i, j * 2, k * 3, nil }

	var callBack = func(err error) {
		finalError = err
	}
	var callBackP1 = func(i int, err error) {
		finalError = err
	}
	// var callBackP2 = func(i, j int, err error) {
	// 	finalError = err
	// }
	// var callBackP3 = func(i, j, k int, err error) {
	// 	finalError = err
	// }

	Convey("Checking function array, and result handler", t, func() {

		Convey("Checking function array pass 1", func() {
			fs := []interface{}{initFuncR2, funcP2R3, funcP3}
			Waterfall(fs, callBack)

			So(finalError.Error(), ShouldEqual, ErrorDefault)

		})

		Convey("Checking function array pass 2", func() {
			fs := []interface{}{initFuncR2, funcP2R3, funcP3R1}
			Waterfall(fs, callBackP1)

			So(finalError.Error(), ShouldEqual, ErrorDefault)

		})

		Convey("Checking function return error", func() {
			fs := []interface{}{initFuncR1, funcP1R3, funcP3R1}
			Waterfall(fs, callBackP1)

			So(finalError.Error(), ShouldEqual, ErrorMsg)

		})

		Convey("Checking value in function array error", func() {
			fs := []interface{}{initFuncR1, "abc", funcP1R3, funcP3R1}

			var p = func () {
				Waterfall(fs, callBackP1)
			}

			So(p, ShouldPanic)
		})

		Convey("Checking value in result handle function error", func() {
			fs := []interface{}{initFuncR1, funcP1R3, funcP3R1}

			var p = func () {
				Waterfall(fs, "abcc")
			}

			So(p, ShouldPanic)
		})

	})
}
