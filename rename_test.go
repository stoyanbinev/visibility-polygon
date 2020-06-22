package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistance_Basic(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(3, 4)
	assert.Equal(t, Distance(a, b), 5.0, "Basic test for point distance")
}

func TestDistance_SamePoints(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(0, 0)
	assert.Equal(t, Distance(a, b), 0.0, "Test for points with same coordinates")
}

func TestDistance_NegativeCoords(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(-3, -4)
	assert.Equal(t, Distance(a, b), 5.0, "Test for point with negative coordinates")
}

func TestEquals_SamePoints(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(0, 0)
	assert.Equal(t, Equals(a, b), true, "Test for points with same coordinates")
}

func TestEquals_DifferentPoints(t *testing.T) {
	a := NewPoint(1, 0)
	b := NewPoint(0, 0)
	assert.Equal(t, Equals(a, b), false, "Test for points with same coordinates")
}

func TestEquals_PolarOpposites(t *testing.T) {
	a := NewPoint(-1, 0)
	b := NewPoint(1, 0)
	assert.Equal(t, Equals(a, b), false, "Test for points with polar opposite coordinates")
}

func TestIntersectLines_NoIntersection(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(3, 3)
	c := NewPoint(0, 4)
	d := NewPoint(0, 5)
	assert.Equal(t, IntersectLines(a, b, c, d), &Point{}, "Test for vectors with no intersection")
}

func TestIntersectLines_BasicIntersection(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(3, 3)
	c := NewPoint(0, 4)
	d := NewPoint(4, 0)
	assert.Equal(t, IntersectLines(a, b, c, d), NewPoint(2, 2), "Test for a basic intersection of vectors")
}

func TestAngle_Basic(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(3, 3)

	assert.Equal(t, Angle(a, b), 45.0, "Test for a basic angle bewtween points")
}

func TestAngle_SamePoint(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(0, 0)

	assert.Equal(t, Angle(a, b), 0.0, "Test for angle from same point")
}

func TestAngle_OppositePoint(t *testing.T) {
	a := NewPoint(1, 1)
	b := NewPoint(-1, -1)

	assert.Equal(t, Angle(a, b), -135.0, "Test for angle from opposite side")
}

func TestAngle_MirrorPoint(t *testing.T) {
	a := NewPoint(1, 1)
	b := NewPoint(-1, 1)

	assert.Equal(t, Angle(a, b), 180.0, "Test for angle from mirroring side")
}

func TestRender_Basic(t *testing.T) {

	a := [][][]float64{{{0, 0}, {4, 0}, {4, 4}, {0, 4}}}
	point := NewPoint(1, 1)
	s := &Scene{}
	assert.Equal(t, s.Render(ConvertPolygonsToSegments(a), point), []*Point{NewPoint(0, 4), NewPoint(0, 0), NewPoint(4, 0), NewPoint(4, 4)}, "Test for angle from mirroring side")
}
