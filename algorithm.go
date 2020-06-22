package main

import (
	"log"
	"math"
	"sort"
)

const epsilon float64 = 0.0000001

// Scene contains environemnt data for a scene
type Scene struct {
	Segments      []*Segment
	Position      *Point
	Buffer        map[int]int // k - segment index, v - heap position
	Heap          []int
	Polygon       []*Point
	Width, Height float64
}

// Point cointains coordinate values for a point in a coordinate system
type Point struct {
	x float64
	y float64
}

// Segment contains values for a line between two points in a coordinate system
type Segment struct {
	a *Point
	b *Point
}

// SegmentSortData contains information for a sorted array of points relative to veiwer
type SegmentSortData struct {
	Index   int
	isFirst bool // Is this the closes point from the segment to viewer
	Angle   float64
}

// NewPoint creates a Point object for the scene
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// NewSegment creates a Segment object for the scene
func NewSegment(x1, y1, x2, y2 float64) *Segment {
	return &Segment{&Point{x1, y1}, &Point{x2, y2}}
}

// Render calculates the visibility polygon of a scene
func (s *Scene) Render(segments []*Segment, position *Point) []*Point {

	var bound []*Segment
	s.Segments = segments
	s.Position = position
	s.Buffer = make(map[int]int)

	log.Println("Getting extremums with position")
	minX := position.x
	minY := position.y
	maxX := position.x
	maxY := position.y
	for i := range segments {

		// Find scene extremums
		minX = math.Min(minX, segments[i].a.x)
		minY = math.Min(minY, segments[i].a.y)
		minX = math.Min(minX, segments[i].b.x)
		minY = math.Min(minY, segments[i].b.y)
		maxX = math.Max(maxX, segments[i].a.x)
		maxY = math.Max(maxY, segments[i].a.y)
		maxX = math.Max(maxX, segments[i].b.x)
		maxY = math.Max(maxY, segments[i].b.y)

		bound = append(bound, segments[i])
	}

	log.Println("Creating segments from extremnus")
	// Create segments from extremums
	bound = append(bound, NewSegment(minX, minY, maxX, minY))
	bound = append(bound, NewSegment(maxX, minY, maxX, maxY))
	bound = append(bound, NewSegment(maxX, maxY, minX, maxY))
	bound = append(bound, NewSegment(minX, maxY, minX, minY))

	log.Println("Sorting ", len(bound), " segments")
	// Sort segments by angle to viewer position
	sorted := sortPoints(position, bound)

	for i := 0; i < len(bound); i++ {
		s.Buffer[i] = -1
	}

	start := NewPoint(position.x+1, position.y)

	log.Println("Checking for active segments")

	for i := 0; i < len(bound); i++ {
		log.Println("Points: ", bound[i].a.x, " ", bound[i].a.y)
		log.Println("Points: ", bound[i].b.x, " ", bound[i].b.y)
		a1 := Angle(bound[i].a, position)
		a2 := Angle(bound[i].b, position)

		active := false
		if a1 > -180 && a1 <= 0 && a2 <= 180 && a2 >= 0 && a2-a1 > 180 {
			active = true
		}
		if a2 > -180 && a2 <= 0 && a1 <= 180 && a1 >= 0 && a1-a2 > 180 {
			active = true
		}
		if active {
			log.Println("Inserting active segment with index: ", i)
			s.insert(i, bound, start)
		}
	}

	s.sweep(sorted, bound)

	s.Heap = nil
	s.Buffer = nil
	resp := s.Polygon
	s.Polygon = nil
	return resp

}

// Traverse list of segments and find shadow points
func (s *Scene) sweep(sorted []*SegmentSortData, bound []*Segment) {
	var vertex *Point
	log.Println("Sweeping ", len(sorted), " segments")
	for i := 0; i < len(sorted); {

		extend := false
		shorten := false
		orig := i

		// Pick according point from segment
		if sorted[i].isFirst {
			vertex = bound[sorted[i].Index].b
		} else {
			vertex = bound[sorted[i].Index].a
		}

		// Initialize variable for top of the heap segment
		oldSegment := s.Heap[0]

		//
		for sorted[i].Angle < sorted[orig].Angle+epsilon {
			// Segment already has point inserted
			if s.Buffer[sorted[i].Index] != -1 {

				// Current segment is original top of the heap
				if sorted[i].Index == oldSegment {
					extend = true

					// Pick according point from segment
					if sorted[i].isFirst {
						vertex = bound[sorted[i].Index].b
					} else {
						vertex = bound[sorted[i].Index].a
					}
				}
				log.Println("Removing segment with index: ", sorted[i].Index)
				s.remove(s.Buffer[sorted[i].Index], bound, vertex)
			} else {
				log.Println("Inserting segment with index: ", sorted[i].Index)
				s.insert(sorted[i].Index, bound, vertex)
				if s.Heap[0] != oldSegment {
					shorten = true
				}
			}
			i++
			if i == len(sorted) {
				break
			}
		}

		if extend {
			s.Polygon = append(s.Polygon, vertex)
			cur := IntersectLines(bound[s.Heap[0]].a, bound[s.Heap[0]].b, s.Position, vertex)
			if !Equals(cur, vertex) {
				s.Polygon = append(s.Polygon, cur)
			}
		} else if shorten {
			s.Polygon = append(s.Polygon, IntersectLines(bound[oldSegment].a, bound[oldSegment].b, s.Position, vertex))
			s.Polygon = append(s.Polygon, IntersectLines(bound[s.Heap[0]].a, bound[s.Heap[0]].b, s.Position, vertex))
		}
	}
}

// Sort points according to angle to viewer
func sortPoints(position *Point, segments []*Segment) []*SegmentSortData {
	var points []*SegmentSortData
	for i := 0; i < len(segments); i++ {
		points = append(points, &SegmentSortData{i, false, Angle(segments[i].a, position)})
		points = append(points, &SegmentSortData{i, true, Angle(segments[i].b, position)})

	}
	// Compare function, on equal angles compares indexes
	sort.Slice(points, func(i, j int) bool {
		if points[i].Angle == points[j].Angle {
			return points[i].Index < points[j].Index
		}
		return points[i].Angle < points[j].Angle
	})
	return points
}

// Angle calculates the angle between 2 points
func Angle(p1, p2 *Point) float64 {

	result := math.Atan2(p2.y-p1.y, p2.x-p1.x) * 180 / math.Pi
	return math.Floor(result*1000) / 1000

}

// insert inserts segment into heap of scene
func (s *Scene) insert(index int, segments []*Segment, destination *Point) {

	intersect := IntersectLines(segments[index].a, segments[index].b, s.Position, destination)
	if intersect == (&Point{}) {
		return
	}
	cur := len(s.Heap)
	s.Heap = append(s.Heap, index)
	s.Buffer[index] = cur

	// Heapify
	for cur > 0 {
		parent := parent(cur)
		if !lessThan(s.Heap[cur], s.Heap[parent], s.Position, segments, destination) {
			break
		}
		s.Buffer[s.Heap[parent]] = cur
		s.Buffer[s.Heap[cur]] = parent
		temp := s.Heap[cur]
		s.Heap[cur] = s.Heap[parent]
		s.Heap[parent] = temp
		cur = parent
	}

}

// remove removes segment from heap of scene
func (s *Scene) remove(index int, segments []*Segment, destination *Point) {

	s.Buffer[s.Heap[index]] = -1
	if index == len(s.Heap)-1 {
		// Remove bottom of heap
		s.Heap = s.Heap[:len(s.Heap)-1]
		return
	}

	s.Heap[index] = s.Heap[len(s.Heap)-1]
	// Remove bottom of heap
	s.Heap = s.Heap[:len(s.Heap)-1]

	s.Buffer[s.Heap[index]] = index
	cur := index
	parent1 := parent(cur)
	if cur != 0 && lessThan(s.Heap[cur], s.Heap[parent1], s.Position, segments, destination) {
		// Heapify
		for cur > 0 {
			parent1 = parent(cur)
			if !lessThan(s.Heap[cur], s.Heap[parent1], s.Position, segments, destination) {
				break
			}
			s.Buffer[s.Heap[parent1]] = cur
			s.Buffer[s.Heap[cur]] = parent1
			temp := s.Heap[cur]
			s.Heap[cur] = s.Heap[parent1]
			s.Heap[parent1] = temp
			cur = parent1
		}
	} else { // Index to be removed is top of the heap
		for {
			left := child(cur)
			right := left + 1
			if left < len(s.Heap) && lessThan(s.Heap[left], s.Heap[cur], s.Position, segments, destination) && (right == len(s.Heap) || lessThan(s.Heap[left], s.Heap[right], s.Position, segments, destination)) {
				// Heapify by exchanging left
				s.Buffer[s.Heap[left]] = cur
				s.Buffer[s.Heap[cur]] = left
				temp := s.Heap[left]
				s.Heap[left] = s.Heap[cur]
				s.Heap[cur] = temp
				cur = left
			} else if right < len(s.Heap) && lessThan(s.Heap[right], s.Heap[cur], s.Position, segments, destination) {
				// Heapify by exchanging right
				s.Buffer[s.Heap[right]] = cur
				s.Buffer[s.Heap[cur]] = right
				temp := s.Heap[right]
				s.Heap[right] = s.Heap[cur]
				s.Heap[cur] = temp
				cur = right
			} else {
				break
			}
		}
	}

}

// IntersectLines finds intersection of 2 vectors
func IntersectLines(a1, a2, b1, b2 *Point) *Point {

	dbx := b2.x - b1.x
	dby := b2.y - b1.y
	dax := a2.x - a1.x
	day := a2.y - a1.y

	uB := dby*dax - dbx*day
	if uB != 0 {
		ua := (dbx*(a1.y-b1.y) - dby*(a1.x-b1.x)) / uB

		return NewPoint(a1.x-ua*-dax, a1.y-ua*-day)
	}
	return &Point{}
}

// parent returns parent node in a binary tree
func parent(index int) int {
	return int(math.Floor(float64((index - 1) / 2)))
}

// lessThan compares two segments
func lessThan(index1 int, index2 int, position *Point, segments []*Segment, destination *Point) bool {
	inter1 := IntersectLines(segments[index1].a, segments[index1].b, position, destination)
	inter2 := IntersectLines(segments[index2].a, segments[index2].b, position, destination)
	// If points not identical compare distance
	if !Equals(inter1, inter2) {
		d1 := Distance(inter1, position)
		d2 := Distance(inter2, position)
		return d1 < d2
	}

	// if first intersection lies on first segment's first point
	end1 := false
	if Equals(inter1, segments[index1].a) {
		end1 = true
	}
	// if second intersection lies on second segment's first point
	end2 := false
	if Equals(inter2, segments[index2].a) {
		end2 = true
	}

	var a1, a2 float64
	if end1 {
		a1 = Angle2(segments[index1].b, inter1, position)
	} else {
		a1 = Angle2(segments[index1].a, inter1, position)
	}
	if end2 {
		a2 = Angle2(segments[index2].b, inter2, position)
	} else {
		a2 = Angle2(segments[index2].a, inter2, position)
	}

	if a1 < 180 {
		if a2 > 180 {
			return true
		}
		return a2 < a1
	}
	return a1 < a2
}

// Equals checks whether two points are close enough to be considered the same
func Equals(a, b *Point) bool {
	if math.Abs(a.x-b.x) < epsilon && math.Abs(a.y-b.y) < epsilon {
		return true
	}
	return false
}

// Distance calculates distance between two points
func Distance(a, b *Point) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return math.Sqrt(dx*dx + dy*dy)
}

func child(index int) int {
	return 2*index + 1
}

//Angle2 calculates the angle between 3 points
func Angle2(a, b, c *Point) float64 {
	a1 := Angle(a, b)
	a2 := Angle(b, c)
	a3 := a1 - a2
	if a3 < 0 {
		a3 += 360
	}
	if a3 > 360 {
		a3 -= 360
	}
	return a3
}
