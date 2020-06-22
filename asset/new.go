package main

import (
	"math"
	"sort"
)

type Scene struct {
	segments []*Segment
	position *Point
	buffer   map[int]int
	heap     []int
	polygon  []*Point
}

type Point struct {
	x float64
	y float64
}

type Segment struct {
	a *Point
	b *Point
}

type Kept struct {
	num       int
	variation bool
	angle     float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func NewSegment(x1, y1, x2, y2 float64) *Segment {
	return &Segment{&Point{x1, y1}, &Point{x2, y2}}
}

func (s *Scene) render(segments []*Segment, position *Point) []*Point {

	var bounded []*Segment
	s.segments = segments
	s.position = position

	minX := position.x
	minY := position.y
	maxX := position.x
	maxY := position.y
	for i := range segments {

		minX = math.Min(minX, segments[i].a.x)
		minY = math.Min(minY, segments[i].a.y)
		minX = math.Min(minX, segments[i].b.x)
		minY = math.Min(minY, segments[i].b.y)
		maxX = math.Max(maxX, segments[i].a.x)
		maxY = math.Max(maxY, segments[i].a.y)
		maxX = math.Max(maxX, segments[i].b.x)
		maxY = math.Max(maxY, segments[i].b.y)

		bounded = append(bounded, segments[i])
	}
	minX--
	minY--
	maxX++
	maxY++

	bounded = append(bounded, NewSegment(minX, minY, maxX, minY))
	bounded = append(bounded, NewSegment(maxX, minY, maxX, maxY))
	bounded = append(bounded, NewSegment(maxX, maxY, minX, maxY))
	bounded = append(bounded, NewSegment(minX, maxY, minX, minY))

	sorted := sortPoints(position, bounded)

	//var s.buffer map[int]int

	for i := 0; i < len(bounded); i++ {
		s.buffer[i] = -1
	}

	//var heap []int
	var vertex *Point
	start := NewPoint(position.x+1, position.y)

	for i := 0; i < len(bounded); i++ {
		a1 := angle(bounded[i].a, position)
		a2 := angle(bounded[i].b, position)
		active := false
		if a1 > -180 && a1 <= 0 && a2 <= 180 && a2 >= 0 && a2-a1 > 180 {
			active = true
		}
		if a2 > -180 && a2 <= 0 && a1 <= 180 && a1 >= 0 && a1-a2 > 180 {
			active = true
		}
		if active {
			s.insert(i, bounded, start)
		}
	}

	for i := 0; i < len(sorted); {

		extend := false
		shorten := false
		orig := i

		if sorted[i].variation {
			vertex = bounded[sorted[i].num].b
		} else {
			vertex = bounded[sorted[i].num].a
		}

		old_segment := s.heap[0]

		for sorted[i].angle < sorted[orig].angle+epsilon() {

			if s.buffer[sorted[i].num] != -1 {
				if sorted[i].num == old_segment {
					extend = true

					if sorted[i].variation {
						vertex = bounded[sorted[i].num].b
					} else {
						vertex = bounded[sorted[i].num].a
					}
				}
				s.remove(s.buffer[sorted[i].num], bounded, vertex)
			} else {
				s.insert(sorted[i].num, bounded, vertex)
				if s.heap[0] != old_segment {
					shorten = true
				}
			}
			i++
			if i == len(sorted) {
				break
			}
		}

		if extend {
			s.polygon = append(s.polygon, vertex)
			cur := intersectLines(bounded[s.heap[0]].a, bounded[s.heap[0]].b, position, vertex)
			if !equal(cur, vertex) {
				s.polygon = append(s.polygon, cur)
			}
		} else if shorten {
			s.polygon = append(s.polygon, intersectLines(bounded[old_segment].a, bounded[old_segment].b, position, vertex))
			s.polygon = append(s.polygon, intersectLines(bounded[s.heap[0]].a, bounded[s.heap[0]].b, position, vertex))
		}
	}
	return s.polygon

}

func sortPoints(position *Point, segments []*Segment) []*Kept {
	var points []*Kept
	for i := 0; i < len(segments); i++ {

		points = append(points, &Kept{i, false, angle(segments[i].a, position)})
		points = append(points, &Kept{i, true, angle(segments[i].b, position)})

	}

	sort.Slice(points, func(i, j int) bool {
		if points[i].angle == points[j].angle {
			return points[i].num < points[j].num
		}
		return points[i].angle < points[j].angle
	})
	return points
}

func angle(p1, p2 *Point) float64 {

	result := math.Atan2(p2.y-p1.y, p2.x-p1.x) * 180 / math.Pi
	return math.Floor(result*1000) / 1000
}

func (s *Scene) insert(index int, segments []*Segment, destination *Point) {

	intersect := intersectLines(segments[index].a, segments[index].b, s.position, destination)
	if intersect == (&Point{}) {
		return
	}
	cur := len(s.heap)
	s.heap = append(s.heap, index)
	s.buffer[index] = cur

	for cur > 0 {
		parent := parent(cur)
		if !lessThan(s.heap[cur], s.heap[parent], s.position, segments, destination) {
			break
		}
		s.buffer[s.heap[parent]] = cur
		s.buffer[s.heap[cur]] = parent
		temp := s.heap[cur]
		s.heap[cur] = s.heap[parent]
		s.heap[parent] = temp
		cur = parent
	}

}

func (s *Scene) remove(index int, segments []*Segment, destination *Point) {

	s.buffer[s.heap[index]] = -1
	if index == len(s.heap)-1 {
		s.heap = s.heap[:len(s.heap)-1]
		return
	}

	s.heap[index] = s.heap[len(s.heap)-1]
	s.heap = s.heap[:len(s.heap)-1]

	s.buffer[s.heap[index]] = index
	cur := index
	parent1 := parent(cur)
	if cur != 0 && lessThan(s.heap[cur], s.heap[parent1], s.position, segments, destination) {
		for cur > 0 {
			parent1 = parent(cur)
			if !lessThan(s.heap[cur], s.heap[parent1], s.position, segments, destination) {
				break
			}
			s.buffer[s.heap[parent1]] = cur
			s.buffer[s.heap[cur]] = parent1
			temp := s.heap[cur]
			s.heap[cur] = s.heap[parent1]
			s.heap[parent1] = temp
			cur = parent1
		}
	} else {
		for {
			left := child(cur)
			right := left + 1
			if left < len(s.heap) && lessThan(s.heap[left], s.heap[cur], s.position, segments, destination) && (right == len(s.heap) || lessThan(s.heap[left], s.heap[right], s.position, segments, destination)) {
				s.buffer[s.heap[left]] = cur
				s.buffer[s.heap[cur]] = left
				temp := s.heap[left]
				s.heap[left] = s.heap[cur]
				s.heap[cur] = temp
				cur = left
			} else if right < len(s.heap) && lessThan(s.heap[right], s.heap[cur], s.position, segments, destination) {
				s.buffer[s.heap[right]] = cur
				s.buffer[s.heap[cur]] = right
				temp := s.heap[right]
				s.heap[right] = s.heap[cur]
				s.heap[cur] = temp
				cur = right
			} else {
				break
			}
		}
	}

}

func intersectLines(a1, a2, b1, b2 *Point) *Point {

	dbx := b2.x - b1.x
	dby := b2.y - b1.y
	dax := a2.x - a1.x
	day := a2.y - a1.y

	u_b := dby*dax - dbx*day
	if u_b != 0 {
		ua := (dbx*(a1.y-b1.y) - dby*(a1.x-b1.x)) / u_b

		return NewPoint(a1.x-ua*-dax, a1.y-ua*-day)
	}
	return &Point{}
}

func parent(index int) int {
	return int(math.Floor(float64((index - 1) / 2)))
}

func lessThan(index1 int, index2 int, position *Point, segments []*Segment, destination *Point) bool {
	inter1 := intersectLines(segments[index1].a, segments[index1].b, position, destination)
	inter2 := intersectLines(segments[index2].a, segments[index2].b, position, destination)
	if !equal(inter1, inter2) {
		d1 := distance(inter1, position)
		d2 := distance(inter2, position)
		return d1 < d2
	}

	end1 := false
	if equal(inter1, segments[index1].a) {
		end1 = true
	}
	end2 := false
	if equal(inter2, segments[index2].a) {
		end2 = true
	}

	var a1, a2 float64
	if end1 {
		a1 = angle2(segments[index1].b, inter1, position)
	} else {
		a1 = angle2(segments[index1].a, inter1, position)
	}
	if end2 {
		a2 = angle2(segments[index2].b, inter2, position)
	} else {
		a2 = angle2(segments[index2].a, inter2, position)
	}

	if a1 < 180 {
		if a2 > 180 {
			return true
		}
		return a2 < a1
	}
	return a1 < a2
}

func equal(a, b *Point) bool {
	if math.Abs(a.x-b.x) < epsilon() && math.Abs(a.y-b.y) < epsilon() {
		return true
	}
	return false
}

func epsilon() float64 {
	return 0.0000001
}

func distance(a, b *Point) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return dx*dx + dy*dy
}

func child(index int) int {
	return 2*index + 1
}

func angle2(a, b, c *Point) float64 {
	a1 := angle(a, b)
	a2 := angle(b, c)
	a3 := a1 - a2
	if a3 < 0 {
		a3 += 360
	}
	if a3 > 360 {
		a3 -= 360
	}
	return a3
}

/*func main() {

	var segments []*Segment
	var scene Scene
	//segments = append(segments, NewSegment(600.6, 200.0, 646.0, 133.3))

	scene.position = NewPoint(250.0, 300.0)
	scene.buffer = make(map[int]int)

	r := [][]float64{{0, 0}, {800, 0}, {800, 500}, {0, 500}}
	a := [][]float64{{600.6, 200.0}, {646.0, 133.3}, {646.6, 261.0}}
	b := [][]float64{{131.0, 188.7}, {54.0, 136.8}, {86.6, 32.9}, {220.8, 32.9}, {238.9, 114.6}, {209.1, 163.6}}
	c := [][]float64{{412.1, 346.9}, {454.1, 251.1}, {537.6, 257.0}, {601.2, 350.3}, {528.8, 430.8}}

	var polygons [][][]float64

	polygons = append(polygons, r, a, b, c)

	for i := 0; i < len(polygons); i++ {
		for j := 0; j < len(polygons[i]); j++ {
			k := j + 1
			if k == len(polygons[i]) {
				k = 0
			}

			segments = append(segments, NewSegment(polygons[i][j][0], polygons[i][j][1], polygons[i][k][0], polygons[i][k][1]))

		}
	}
	scene.segments = segments

	result := scene.render(scene.segments, scene.position)
	for i := 0; i < len(result); i++ {
		log.Println(result[i].x, result[i].y)
	}

}*/
