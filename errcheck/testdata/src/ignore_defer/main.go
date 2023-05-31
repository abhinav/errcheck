// ensure that the package keyword is not equal to file beginning
// to test correct position calculations.
package ignore_defer

import (
	"fmt"
	"log"
	"os"
)

func a() error {
	fmt.Println("this function returns an error") // ok, excluded
	return nil
}

func b() (int, error) {
	fmt.Println("this function returns an int and an error") // ok, excluded
	return 0, nil
}

func c() int {
	fmt.Println("this function returns an int") // ok, excluded
	return 7
}

func rec() {
	defer func() {
		recover()     // want "unchecked error"
		_ = recover() // ok, assigned to blank
	}()
	defer recover() // ok, ignore defer

	os.Open("filename.ext") // want "unchecked error"

	f, err := os.Open("filename.ext")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // ok, ignore defer
}

type MyError string

func (e MyError) Error() string {
	return string(e)
}

func customError() error {
	return MyError("an error occurred")
}

func customConcreteError() MyError {
	return MyError("an error occurred")
}

func customConcreteErrorTuple() (int, MyError) {
	return 0, MyError("an error occurred")
}

type MyPointerError string

func (e *MyPointerError) Error() string {
	return string(*e)
}

func customPointerError() *MyPointerError {
	e := MyPointerError("an error occurred")
	return &e
}

func customPointerErrorTuple() (int, *MyPointerError) {
	e := MyPointerError("an error occurred")
	return 0, &e
}

type ErrorMakerInterface interface {
	MakeNilError() error
}
type ErrorMakerInterfaceWrapper interface {
	ErrorMakerInterface
}

func main() {
	// Single error return
	defer a() // ok, ignore defer
	defer func() {
		a() // want "unchecked error"
	}()

	// Return another value and an error
	defer b() // ok, ignore defer
	defer func() {
		b() // want "unchecked error"
	}()

	// Return a custom error type
	defer customError() // ok, ignore defer
	defer func() {
		customError() // want "unchecked error"
	}()

	// Return a custom concrete error type
	defer customConcreteError() // ok, ignore defer
	defer func() {
		customConcreteError() // want "unchecked error"
	}()

	defer customConcreteErrorTuple() // ok, ignore defer
	defer func() {
		customConcreteErrorTuple() // want "unchecked error"
	}()

	// Return a custom pointer error type
	defer customPointerError() // ok, ignore defer
	defer func() {
		customPointerError() // want "unchecked error"
	}()

	defer customPointerErrorTuple() // ok, ignore defer
	defer func() {
		customPointerErrorTuple() // want "unchecked error"
	}()

	// Method with a single error return
	x := t{}
	defer x.a() // ok, ignore defer
	defer func() {
		x.a() // want "unchecked error"
	}()

	// Method call on a struct member
	y := u{x}
	defer y.t.a() // ok, ignore defer
	defer func() {
		y.t.a() // want "unchecked error"
	}()

	m1 := map[string]func() error{"a": a}
	defer m1["a"]() // ok, ignore defer
	defer func() {
		m1["a"]() // want "unchecked error"
	}()
}

type t struct{}

func (x t) a() error {
	return nil
}

type u struct {
	t t
}
