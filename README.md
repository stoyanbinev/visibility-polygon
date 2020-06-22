
## Description
This projects implements the Algorithm of Asano (https://arxiv.org/pdf/1403.3905.pdf) to solve the Visibility Polygon problem (https://en.wikipedia.org/wiki/Visibility_polygon). The algorithm has a complexity of O(Nlog(N)) and is implemented in Go, with display of interactive environment done through a local web server and socket.io.

## To Run

To run:
`go run server.go algorithm.go`

To test:
`go test`

## Interface

After you have ran the application you can go on `localhost:8000` to view the interface. There is an `Open` link on the bottom left from which you can load a scene file. There is a sample scene(scene.txt) provided in the project. The source of light can be moved by clicking and dragging. Also, the shapes can be altered by clicking and dragging on the edges. On the bottom left the percentage of lit area in the scene is dynamically changed.
![alt text](https://i.imgur.com/RwvvCO9.gif "gif")

![alt text](https://i.imgur.com/o2ue1K6.png "Image 1")
![alt text](https://i.imgur.com/mncXv3S.png "Image 2")

