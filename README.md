GeeFlight: another way to organize your golang code
=================================================

Welcome to GeeFlight, functionalized your business logic without if-err anywhere. The idea is inspired by [caolan/async](http://caolan.github.io/async/).

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
    i: 5, j: 6, error: sum is more than 10
