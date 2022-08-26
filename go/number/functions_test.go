package number

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxFloat64(t *testing.T) {
	testCases := []struct {
		value    []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 10, 5, 3, 1000, 1001}, 1001},
		{[]float64{5885.588997, 0.5566, 5885.588998}, 5885.588998},
		{[]float64{0, 10, 10, 36}, 36},
		{[]float64{100, 10, 10, 36}, 100},
		{[]float64{100.5, 10, 10, 36}, 100.5},
		{[]float64{-1, -1, 1.0001, 1.00001, 0}, 1.0001},
	}
	for _, v := range testCases {
		got := Max(v.value...)
		assert.Equal(t, v.expected, got, "", v.value, v.expected, got)
	}
}

func TestMaxInt(t *testing.T) {
	testCases := []struct {
		value    []int
		expected int
	}{
		{[]int{1, 2, 3, 10, 5, 3, 1000, 1001}, 1001},
		{[]int{5885, 0, 5885}, 5885},
		{[]int{10, 11, 9}, 11},
	}
	for _, v := range testCases {
		got := Max(v.value...)
		assert.Equal(t, v.expected, got, "", v.value, v.expected, got)
	}
}

func TestMinFloat64(t *testing.T) {
	testCases := []struct {
		value    []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 10, 5, 3, 1000, 1001}, 1},
		{[]float64{5885.588997, 0.5566, 5885.588998}, 0.5566},
		{[]float64{0, 10, 10, 36}, 0},
		{[]float64{100.5, 10, 10, 36}, 10},
		{[]float64{-1, 0, 0, 1}, -1},
		{[]float64{-1, -1, -1.0001, 0}, -1.0001},
		{[]float64{-1, -1, -1.0001, -1.00001, 0}, -1.0001},
	}
	for _, v := range testCases {
		got := Min(v.value...)
		assert.Equal(t, v.expected, got, "", v.value, v.expected, got)
	}
}

func TestMinInt(t *testing.T) {
	testCases := []struct {
		value    []int
		expected int
	}{
		{[]int{1, 2, 3, 10, 5, 3, 1000, 1001}, 1},
		{[]int{5885, 0, 5885}, 0},
		{[]int{10, 11, 9}, 9},
		{[]int{-5, -6, 5, 0, 0, 8}, -6},
	}
	for _, v := range testCases {
		got := Min(v.value...)
		assert.Equal(t, v.expected, got, "", v.value, v.expected, got)
	}
}
