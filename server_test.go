package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateArea_Basic(t *testing.T) {

	input := [][]float64{{0, 0}, {4, 0}, {4, 4}, {0, 4}}

	assert.Equal(t, CalculateArea(input), 16.0, "Basic test for n-figure surface")
}

func TestCalculateArea_Zero(t *testing.T) {

	input := [][]float64{{0, 0}, {0, 0}, {0, 0}, {0, 0}}

	assert.Equal(t, CalculateArea(input), 0.0, "Test for no surface")
}

func TestCalculateArea_Figure(t *testing.T) {

	input := [][]float64{{0, 0}, {2, 4}, {3, 2}, {3, 6}, {2, 9}}

	assert.Equal(t, CalculateArea(input), 9.5, "Test for no surface")
}
