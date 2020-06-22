package algorithm

import (
    "fmt"
    "math"
    "sort"
)

type Scene struct {
    output        []*Point
    segments      []*Segment
    open          []*Segment
    endpoints     []*EndPoint
    center        *Point
    width, height float64
}

type Point struct {
    x float64
    y float64
}

type EndPoint struct {
    x       float64
    y       float64
    begin   bool
    segment *Segment
    angle   float64
}

type Segment struct {
    p1, p2 *EndPoint
    d      float64
}

func NewPoint(x, y float64) *Point {
    return &Point{x, y}
}

func NewEndPoint(x, y float64) *EndPoint {
    return &EndPoint{x, y, false, nil, 0.0}
}

func (s *Scene) addSegment(x1, y1, x2, y2 float64) {
    var segment = &Segment{}
    var p1 = NewEndPoint(x1, y1)
    p1.segment = segment
    var p2 = NewEndPoint(x2, y2)
    p2.segment = segment

    segment = &Segment{}
    p1.x = x1
    p1.y = y1
    p2.x = x2
    p2.y = y2
    p1.segment = segment
    p2.segment = segment
    segment.p1 = p1
    segment.p2 = p2
    segment.d = 0.0
    s.segments = append(s.segments, segment)
    s.endpoints = append(s.endpoints, p1)
    s.endpoints = append(s.endpoints, p2)

}

func (s *Scene) setLightLocation(x, y float64) {
    s.center.x = x
    s.center.y = y

    const Pi = 3.14

    for i := range s.segments {

        dx := 0.5*(s.segments[i].p1.x+s.segments[i].p2.x) - x
        dy := 0.5*(s.segments[i].p1.y+s.segments[i].p2.y) - y

        s.segments[i].d = dx*dx + dy*dy

        s.segments[i].p1.angle = math.Atan2(s.segments[i].p1.y-y, s.segments[i].p1.x-x)
        s.segments[i].p2.angle = math.Atan2(s.segments[i].p2.y-y, s.segments[i].p2.x-x)

        dAngle := s.segments[i].p2.angle - s.segments[i].p1.angle
        if dAngle <= -Pi {
            dAngle += 2 * Pi
        }
        if dAngle > Pi {
            dAngle -= 2 * Pi
        }
        s.segments[i].p1.begin = (dAngle > 0.0)
        s.segments[i].p2.begin = !s.segments[i].p1.begin
    }
}

func leftOf(s *Segment, p *EndPoint) bool {
    // Cross product for 3D vectors
    cross := (s.p2.x-s.p1.x)*(p.y-s.p1.y) - (s.p2.y-s.p1.y)*(p.x-s.p1.x)
    return cross < 0
}

func leftOfPoint(s *Segment, p *Point) bool {
    // Cross product for 3D vectors
    cross := (s.p2.x-s.p1.x)*(p.y-s.p1.y) - (s.p2.y-s.p1.y)*(p.x-s.p1.x)
    return cross < 0

}

func interpolate(p *EndPoint, q *EndPoint, f float64) *EndPoint {
    return NewEndPoint(p.x*(1-f)+q.x*f, p.y*(1-f)+q.y*f)
}

// Check if segment a is in front of segment b
func segment_in_front_of(a *Segment, b *Segment, relativeTo *Point) bool {
    A1 := leftOf(a, interpolate(b.p1, b.p2, 0.0))
    A2 := leftOf(a, interpolate(b.p2, b.p1, 0.0))
    A3 := leftOfPoint(a, relativeTo)
    B1 := leftOf(b, interpolate(a.p1, a.p2, 0.0))
    B2 := leftOf(b, interpolate(a.p2, a.p1, 0.0))
    B3 := leftOfPoint(b, relativeTo)

    // If B is on one side of A and relativeTo on the other
    // this means that A is between B and the viewer
    if B1 == B2 && B2 != B3 {
        return true
    }
    if A1 == A2 && A2 == A3 {
        return true
    }
    if A1 == A2 && A2 != A3 {
        return false
    }
    if B1 == B2 && B2 == B3 {
        return false
    }

    return false

}

func (s *Scene) sweep() {
    maxAngle := 999.0

    // Sort endpoints
    sort.Slice(s.endpoints, func(i, j int) bool {

        // Firt by angle order
        if s.endpoints[i].angle > s.endpoints[j].angle {
            return true
        }
        if s.endpoints[i].angle < s.endpoints[j].angle {
            return false
        }
        // For ties put higher priority on begin nodes
        if !s.endpoints[i].begin && s.endpoints[j].begin {
            return true
        }
        if s.endpoints[i].begin && !s.endpoints[j].begin {
            return false
        }
        return true

    })

    s.open = nil
    beginAngle := 0.0

    // Traverse to collect active segments
    for _, p := range s.endpoints {

        var current_old *Segment
        if len(s.open) < 1 {
            current_old = nil
        } else {
            current_old = s.open[0]
        }

        if p.begin {

            var node *Segment
            i := 0
            if len(s.open) < 1 {
                node = nil
            } else {
                node = s.open[i]
            }

            for {
                if node == nil || segment_in_front_of(p.segment, node, s.center) {
                    break
                }
                i = i + 1
                node = s.open[i]
            }

            if node == nil {
                s.open = append(s.open, p.segment)
            } else {
                s.open = append([]*Segment{p.segment}, s.open...)
            }

            var current_new *Segment
            if len(s.open) < 1 {
                current_new = nil
            } else {
                current_new = s.open[0]
            }

            if current_old != current_new {

                s.addTriangle(beginAngle, p.angle, current_old)

            }

        }

    }

}

func lineIntersection(p1, p2, p3, p4 *Point) *Point {
    // http://paulbourke.net/geometry/lineline2d/
    s := ((p4.x-p3.x)*(p1.y-p3.y) - (p4.y-p3.y)*(p1.x-p3.x)) / ((p4.y-p3.y)*(p2.x-p1.x) - (p4.x-p3.x)*(p2.y-p1.y))

    return NewPoint(p1.x+s*(p2.x-p1.x), p1.y+s*(p2.y-p1.y))
}

// Create triangle of lit are
func (s *Scene) addTriangle(angle1 float64, angle2 float64, segment *Segment) {
    p1 := s.center
    p2 := NewPoint(s.center.x+math.Cos(angle1), s.center.y+math.Sin(angle1))
    p3 := NewPoint(0.0, 0.0)
    p4 := NewPoint(0.0, 0.0)

    // Stop the triangle at the intersecting segment
    p3.x = segment.p1.x
    p3.y = segment.p1.y
    p4.x = segment.p2.x
    p4.y = segment.p2.y

    pBegin := lineIntersection(p3, p4, p1, p2)

    p2.x = s.center.x + math.Cos(angle2)
    p2.y = s.center.y + math.Sin(angle2)
    pEnd := lineIntersection(p3, p4, p1, p2)

    s.output = append(s.output, pBegin, pEnd)

}
