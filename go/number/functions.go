package number

import (
	"golang.org/x/exp/constraints"
)

//Max return the maximum value given variadic argument
func Max[T constraints.Ordered](values ...T) (r T) {
	return MinMax(false, values...)
}

//Min return the maximum value given variadic argument
func Min[T constraints.Ordered](values ...T) (r T) {
	return MinMax(true, values...)
}

//MinMax return the min or max given variadic argument
func MinMax[T constraints.Ordered](min bool, values ...T) (r T) {
	if values == nil {
		return
	}
	r = values[0]
	for _, v := range values {
		if min {
			if v < r {
				r = v
			}
			continue
		}
		if v > r {
			r = v
		}
	}
	return
}
