GeeFlight: another way to organize your golang code
=================================================

Welcome to GeeFlight, functionalized your business logic without if-err everywhere. The idea is inspired by [caolan/async](http://caolan.github.io/async/).

*Menu*

- [Installation](#installation)
- [Example](#example)
- [Usage](#usage)

Installation
------------

    $ go get github.com/szqmtl/geeflight


Example
-------

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

Output:

    i: 1, j: 1
    i: 2, j: 1
    i: 3, j: 2
    i: 5, j: 6, error: <nil>
    i: 1, j: 1
    i: 2, j: 1
    i: 3, j: 2
    i: 5, j: 6
    i: 0, j: 0, error: sum is more than 10
    i: 1, j: 1
    i: 2, j: 1
    i: 3, j: 2
    i: 5, j: 6
    i: 0, j: 0, error: sum is more than 10

Usage
-----

- Method Waterfall takes two arguments: function list and optional result handle function
- If giving a non function parameter, a panic raises
- Assumptions/approach of function list/result handle function
    - The preceding function output is the following one input
    - The last function output is the result handler function's input
    - The last return item of a function in the list is error
    - The error object is not as part of the function input, but it is in result handler argument list
    - If the error object is not nil, the function output is the result handle function input
    - If the quantity of a output is more than the quantity of the input argument, the excessive part is abandoned, but for the result handle function, the error object is always kept as the last argument.
    - If the quantity of a output is less than the quantity of the input argument, the system default value is assigned to the missing part, but for the result handle function, the error object is always kept as the last argument.
    - The return of result handle function is ignored.