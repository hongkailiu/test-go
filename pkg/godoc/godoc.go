// Package godoc shows how to use godoc.
//
// See https://blog.golang.org/godoc-documenting-go-code
//
// by Hongkai Liu
package godoc

//Deprecated: NewHello is defined here to show
//the usage for `Deprecated:` in goDoc\
//
// BUG(who): The NewHello did not say hello nicely.
func NewHello() error {
	return nil
}

//NewNiceHello says hello nicely.
func NewNiceHello() error {
	return nil
}
