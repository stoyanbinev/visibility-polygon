package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/googollee/go-socket.io"
)

type Response struct {
	ViewerX     float64       `json:"viewer_x"`
	ViewerY     float64       `json:"viewer_y"`
	LightX      float64       `json:"light_x"`
	LightY      float64       `json:"light_y"`
	AreaPercent float64       `json:"area_percent"`
	Segments    [][][]float64 `json:"segments"`
	Polygon     [][]float64   `json:"polygon"`
}

var sceneData = &Scene{}

func main() {

	var position []float64
	var inputSegments [][][]float64

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/scene", "input", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)

		return ComputeTextfile(msg)
	})

	server.OnEvent("/scene", "light", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		json.Unmarshal([]byte(msg), &position)

		return changeLight(position[0], position[1])
	})

	server.OnEvent("/scene", "segments", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		json.Unmarshal([]byte(msg), &inputSegments)

		return changeSegments(inputSegments)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	log.Println("Creating SocketIO server")

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Calculates area with Shoelace formula
func CalculateArea(polygons [][]float64) float64 {

	area := 0.0

	j := len(polygons) - 1
	for i := 0; i < len(polygons); i++ {
		area += (polygons[j][0] + polygons[i][0]) * (polygons[j][1] - polygons[i][1])
		j = i
	}

	// Round to 1 decimal place
	return math.Round(math.Abs(area/2.0)*100) / 100
}

// Convert input float 3D array of polygons to segments object
func ConvertPolygonsToSegments(polygons [][][]float64) []*Segment {
	var segments []*Segment
	log.Println("Converting ", len(polygons), " segments")
	for i := 0; i < len(polygons); i++ {
		for j := 0; j < len(polygons[i]); j++ {
			k := j + 1
			if k == len(polygons[i]) {
				k = 0
			}

			segments = append(segments, NewSegment(polygons[i][j][0], polygons[i][j][1], polygons[i][k][0], polygons[i][k][1]))

		}
	}

	return segments
}

// Compute scene
func compute(polygons [][][]float64, lightX, lightY float64) [][]float64 {

	log.Println("Computing scene")

	sceneData.Position = NewPoint(lightX, lightY)
	sceneData.Buffer = make(map[int]int)

	sceneData.Segments = ConvertPolygonsToSegments(polygons)

	result := sceneData.Render(sceneData.Segments, sceneData.Position)

	var response [][]float64

	for i := 0; i < len(result); i++ {
		response = append(response, []float64{result[i].x, result[i].y})
	}

	return response
}

// Renders scene with new light position
func changeLight(lightX, lightY float64) string {

	log.Println("Rendering with light postion: ", lightX, " ", lightY)
	sceneData.Position = NewPoint(lightX, lightY)

	result := sceneData.Render(sceneData.Segments, sceneData.Position)

	var resultedPolygon [][]float64

	for i := 0; i < len(result); i++ {
		resultedPolygon = append(resultedPolygon, []float64{result[i].x, result[i].y})
	}

	areaPercent := 100 / ((sceneData.Width * sceneData.Height) / CalculateArea(resultedPolygon))

	response := &Response{}
	response.AreaPercent = areaPercent
	response.Polygon = resultedPolygon

	responseJson, _ := json.Marshal(response)

	return string(responseJson)

}

// Renders scene with edited segments
func changeSegments(inputSegments [][][]float64) string {

	log.Println("Rendering with edited segments")

	sceneData.Segments = ConvertPolygonsToSegments(inputSegments)

	result := sceneData.Render(sceneData.Segments, sceneData.Position)

	var resultedPolygon [][]float64

	for i := 0; i < len(result); i++ {
		resultedPolygon = append(resultedPolygon, []float64{result[i].x, result[i].y})
	}

	areaPercent := 100 / ((sceneData.Width * sceneData.Height) / CalculateArea(resultedPolygon))

	response := &Response{}
	response.AreaPercent = areaPercent
	response.Polygon = resultedPolygon

	responseJson, _ := json.Marshal(response)

	return string(responseJson)

}

// Compute scene from text file
func ComputeTextfile(input string) string {

	log.Println("Computing from text file")

	input = strings.Replace(input, "\r", "", -1)

	lines := strings.Split(input, "\n")
	dimensions := strings.Split(lines[0], " ")

	x, err := strconv.ParseFloat(dimensions[0], 64)
	if err != nil {
		log.Println("Error", err)
	}

	y, err := strconv.ParseFloat(dimensions[1], 64)
	if err != nil {
		log.Println("Error", err)
	}

	light := strings.Split(lines[1], " ")

	lightX, err := strconv.ParseFloat(light[0], 64)
	if err != nil {
		log.Println("Error", err)
	}

	lightY, err := strconv.ParseFloat(light[1], 64)
	if err != nil {
		log.Println("Error", err)
	}

	var polygons [][][]float64

	polygons = append(polygons, [][]float64{{0, 0}, {x, 0}, {x, y}, {0, y}})

	sceneData.Width = y
	sceneData.Height = x

	log.Println("Converting file line to segments")

	for i := 3; i < len(lines); i++ {

		points := strings.Split(lines[i], " ")
		counter := 0
		segment := [][]float64{}
		point := []float64{}
		for j := 1; j < len(points); j++ {
			coordinate, err := strconv.ParseFloat(points[j], 64)
			if err != nil {
				fmt.Printf("Error")
			}

			if counter < 2 {
				point = append(point, coordinate)
				counter++
			} else {
				segment = append(segment, point)
				point = nil
				point = append(point, coordinate)
				counter = 1
			}
		}

		segment = append(segment, point)
		point = nil
		polygons = append(polygons, segment)
		segment = nil
	}
	result := compute(polygons, lightX, lightY)

	log.Println("Calculating lit surface")
	areaPercent := 100 / ((sceneData.Width * sceneData.Height) / CalculateArea(result))

	response := &Response{x, y, lightX, lightY, areaPercent, polygons, result}

	q, _ := json.Marshal(response)

	return string(q)
}
