/**
 * @Author: jie.an
 * @Description:
 * @File:  func_opt.go
 * @Version: 1.0.0
 * @Date: 2020/7/23 16:49
 */
package example

import "fmt"

func funcOpt() {
	h := NewHouse(
		WithConcrete(),
		WithoutFireplace(),
		WithFloors(3),
	)
	fmt.Println(h)
}

////Functional Options
// https://www.sohamkamani.com/golang/options-pattern/
type House struct {
	Material     string
	HasFireplace bool
	Floors       int
}

// `NewHouse` is a constructor function for `*House`
// NewHouse now takes a slice of option as the rest arguments
func NewHouse(opts ...HouseOption) *House {
	const (
		defaultFloors       = 2
		defaultHasFireplace = true
		defaultMaterial     = "wood"
	)

	h := &House{
		Material:     defaultMaterial,
		HasFireplace: defaultHasFireplace,
		Floors:       defaultFloors,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *House as the argument
		opt(h)
	}

	// return the modified house instance
	return h
}

type HouseOption func(*House)

func WithConcrete() HouseOption {
	return func(h *House) {
		h.Material = "concrete"
	}
}

func WithoutFireplace() HouseOption {
	return func(h *House) {
		h.HasFireplace = false
	}
}
func WithFloors(floors int) HouseOption {
	return func(h *House) {
		h.Floors = floors
	}
}
