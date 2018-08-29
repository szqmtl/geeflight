package geeflight

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestGuard(t *testing.T) {

	Convey("Checking only one nil return", t, func() {

		f := func() error { return nil }

		tf := func() { Guard(f) }

		assert.NotPanics(t, tf)
	})

	var ErrorMsg string = "non nil return"
	var RaisedError error = errors.New(ErrorMsg)

	Convey("Checking with only one non-nil return", t, func() {

		f := func() error { return RaisedError }

		tf := func() { Guard(f) }

		assert.PanicsWithValue(t, RaisedError, tf)
	})

	Convey("Checking with only one non-nil return and extra parameter", t, func() {

		f := func() error { return RaisedError }

		var i int
		tf := func() { Guard(f, &i) }

		assert.PanicsWithValue(t, ErrorMsgWrongArgNumber, tf)
	})

	Convey("Checking with two return", t, func() {

		var testInt int = 1
		f := func() (int, error) { return testInt, nil }

		var i int
		tf := func() { Guard(f, &i) }

		assert.NotPanics(t, tf)
		assert.Equal(t, testInt, i)
	})

	Convey("Checking with two return but missed one parameter", t, func() {

		var testInt int = 1
		f := func() (int, error) { return testInt, nil }

		tf := func() { Guard(f) }

		assert.PanicsWithValue(t, ErrorMsgWrongArgNumber, tf)
	})

	Convey("Checking with one return and error", t, func() {

		var testInt int = 1
		f := func() (int, error) { return testInt, RaisedError }

		var i int
		tf := func() { Guard(f, &i) }

		assert.PanicsWithValue(t, RaisedError, tf)
	})

	Convey("Checking with three return", t, func() {

		var testInt int = 1
		var testString string = "abc"
		f := func() (int, string, error) { return testInt, testString, nil }

		var i int
		var s string
		tf := func() { Guard(f, &i, &s) }

		assert.NotPanics(t, tf)
		assert.Equal(t, testInt, i)
		assert.Equal(t, testString, s)
	})

	Convey("Checking with two return and error", t, func() {

		var testInt int = 1
		var testString string = "abc"
		f := func() (int, string, error) { return testInt, testString, RaisedError }

		var i int
		var s string
		tf := func() { Guard(f, &i, &s) }

		assert.PanicsWithValue(t, RaisedError, tf)
	})
}

func TestCatchGuard(t *testing.T) {
	var ErrorMsg string = "non nil return"
	var RaisedError error = errors.New(ErrorMsg)

	Convey("Checking with valid panic", t, func() {
		defer CatchGuard(func(e error) {
			assert.Equal(t, RaisedError, e)
		})

		f := func() error { return RaisedError }

		Guard(f)
	})

	Convey("Checking with nil panic", t, func() {
		defer CatchGuard(func(e error) {
			fmt.Printf("error: %v\n", e)
		})

		f := func() error { return nil }

		Guard(f)
		assert.Equal(t, 1, 1)
	})


	Convey("Checking with string panic", t, func() {

		f := func() string { return ErrorMsg }


		tf := func() {
			defer CatchGuard(func(e error) {
				fmt.Printf("error: %v\n", e)
			})

			Guard(f)
		}

		assert.PanicsWithValue(t, ErrorMsg, tf)
	})
}

func TestCatchCGuard(t *testing.T) {
	var ErrorMsg string = "non nil return"
	var RaisedError error = errors.New(ErrorMsg)
	var testInt int = 1

	Convey("Checking with valid panic", t, func() {
		defer CatchCGuard(func(i interface{}, e error) {
			assert.Equal(t, testInt, i.(int))
			assert.Equal(t, RaisedError, e)
		})

		f := func()(error) { return RaisedError }

		CGuard(testInt, f)
	})

	Convey("Checking with nil panic", t, func() {
		defer CatchCGuard(func(i interface{}, e error) {
			fmt.Printf("error: %v\n", e)
		})

		f := func() error { return nil }

		CGuard(testInt, f)
		assert.Equal(t, 1, 1)
	})


	Convey("Checking with string panic", t, func() {

		f := func() string { return ErrorMsg }


		tf := func() {
			defer CatchCGuard(func(i interface{}, e error) {
				fmt.Printf("error: %v\n", e)
			})

			CGuard(testInt, f)
		}

		assert.PanicsWithValue(t, [...]interface{}{testInt, ErrorMsg}, tf)
	})


	Convey("Checking with string panic", t, func() {

		f := func() string { return ErrorMsg }


		tf := func() {
			defer CatchCGuard(func(i interface{}, e error) {
				fmt.Printf("error: %v\n", e)
			})

			CGuard(ErrorMsg, f)
		}

		assert.PanicsWithValue(t, [...]interface{}{ErrorMsg, ErrorMsg}, tf)
	})
}


func TestCatchIntCGuard(t *testing.T) {
	var ErrorMsg string = "non nil return"
	var RaisedError error = errors.New(ErrorMsg)
	var testInt int = 1

	Convey("Checking with valid panic", t, func() {
		defer CatchIntCGuard(func(i int, e error) {
			assert.Equal(t, testInt, i)
			assert.Equal(t, RaisedError, e)
		})

		f := func()(error) { return RaisedError }

		CGuard(testInt, f)
	})

	Convey("Checking with nil panic", t, func() {
		defer CatchIntCGuard(func(i int, e error) {
			fmt.Printf("error: %v\n", e)
		})

		f := func() error { return nil }

		CGuard(testInt, f)
		assert.Equal(t, 1, 1)
	})


	Convey("Checking with string panic", t, func() {

		f := func() string { return ErrorMsg }


		tf := func() {
			defer CatchIntCGuard(func(i int, e error) {
				fmt.Printf("error: %v\n", e)
			})

			CGuard(testInt, f)
		}

		assert.PanicsWithValue(t, [...]interface{}{testInt, ErrorMsg}, tf)
	})


	Convey("Checking with string panic", t, func() {

		f := func() string { return ErrorMsg }


		tf := func() {
			defer CatchIntCGuard(func(i int, e error) {
				fmt.Printf("error: %v\n", e)
			})

			CGuard(ErrorMsg, f)
		}

		assert.PanicsWithValue(t, [...]interface{}{ErrorMsg, ErrorMsg}, tf)
	})
}