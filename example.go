package main

import (
	"errors"
	"fmt"
	"github.com/szqmtl/geeflight/geeflight"
)

func sum(i, j int) (int, int, error) {
	fmt.Printf("i: %d, j: %d\n", i, j)
	if (i + j) > 10 {
		return 0, 0, errors.New("sum is more than 10")
	}
	return i + j, i*j, nil
}

func main() {
	var init = func()(int, int, error) { return sum(1, 1) }
	var handler = func(i, j int, err error){
		fmt.Printf("i: %d, j: %d, error: %v\n", i, j, err)
	}

	geeflight.Waterfall(
		[]interface{}{ init, sum, sum },
		handler )

	geeflight.Waterfall(
		[]interface{}{ init, sum, sum, sum },
		handler )

	geeflight.Waterfall(
		[]interface{}{ init, sum, sum, sum, sum },
		handler )
}
