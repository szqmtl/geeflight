package geeflight

import (
	"testing"
	"errors"
)

func TestWaterfall(t *testing.T){
	ERR_DEFAULT := "DEFAULT"
	ERR_MSG := "IN ERROR"
	var finalError = errors.New(ERR_DEFAULT)
	var func1 = func()(int, error) { return 1, nil }
	var func2 = func(i int)(int, int, error) { return i, i*2, errors.New(ERR_MSG) }
	var func3 = func(i, j int)(int, error) { return j*2, nil }
	var callBack = func(i int, err error) {
		finalError = err
	}

	fs := []interface{}{func1, func2, func3}

	Waterfall(fs, callBack)

	if finalError.Error() != ERR_MSG {
		t.Error("The error is not corrected(%v,%v)", finalError, ERR_MSG)
	}
}

