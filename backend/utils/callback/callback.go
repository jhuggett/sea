package callback

import "fmt"

func functionsAreTheSame(a, b interface{}) bool {
	ap := fmt.Sprintf("%p", a)
	bp := fmt.Sprintf("%p", b)

	return ap == bp
}
