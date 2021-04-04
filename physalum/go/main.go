package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"os"
)

const width = 240
const height = 240
const sensorAngle = 25
const sensorDistance = 9

type Agent struct {
	x              float64
	y              float64
	direction      float64
	sensorAngle    float64
	sensorDistance float64
	fl             float64
	f              float64
	fr             float64
}

func (a *Agent) move() {
	x, y := move(a.x, a.y, a.direction, 1)
	if isOutOfBounds(x, y, width, height) {
		a.perturbRandomly(180)
		a.move()
	} else {
		a.x = x
		a.y = y
	}
}

func (a *Agent) perturbRandomly(angle float64) {
	a.direction = float64(int(a.direction+randSymmetricRange(angle)) % 360)
}

func (a *Agent) perturb() {
	if a.f > a.fl && a.f > a.fr {
		return
	} else if a.f < a.fl && a.f < a.fr {
		if rand.Float32() < 0.5 {
			a.direction = float64(int(a.direction-a.sensorAngle) % 360)
		} else {
			a.direction = float64(int(a.direction+a.sensorAngle) % 360)
		}
	} else if a.fl < a.fr {
		a.direction = float64(int(a.direction+a.sensorAngle) % 360)
	} else if a.fr < a.fl {
		a.direction = float64(int(a.direction-a.sensorAngle) % 360)
	} else {
		return
	}
}

func (a *Agent) readSensors(grid [width * height]float64) {
	a.fl = a.getSensorReading(float64(int(a.direction-a.sensorAngle)%360), grid)
	a.f = a.getSensorReading(a.direction, grid)
	a.fr = a.getSensorReading(float64(int(a.direction+a.sensorAngle)%360), grid)
}

func (a *Agent) getSensorReading(angle float64, grid [width * height]float64) float64 {
	x, y := move(a.x, a.y, angle, a.sensorDistance)
	if isOutOfBounds(x, y, width, height) {
		return 0
	}
	return grid[int(math.Floor(y))*width+int(math.Floor(x))]
}

func (a *Agent) deposit(grid *[width * height]float64, amount float64) {
	grid[int(math.Floor(a.y))*width+int(math.Floor(a.x))] += amount
}

func main() {
	var grid [width * height]float64
	colorPalette := GetColorPalette(255)
	images := []*image.Paletted{}

	var agents = []Agent{}
	for i := 0; i < 1; i++ {
		agent := Agent{
			x:              rand.Float64() * width,
			y:              rand.Float64() * height,
			direction:      rand.Float64() * 360,
			sensorAngle:    sensorAngle,
			sensorDistance: sensorDistance}
		agents = append(agents, agent)
	}

	// run simulation
	var iter = 100
	delays := []int{}
	for iter > 0 {
		for _, agent := range agents {
			fmt.Println(agent)
			agent.readSensors(grid)
			agent.perturb()
			agent.move()
			agent.deposit(&grid, 5)
		}
		decay(&grid)
		grid = blur(&grid)

		img := CreateImage(agents, colorPalette)
		images = append(images, img)

		iter--
		delays = append(delays, 0)
	}

	anim := gif.GIF{Delay: delays, Image: images}

	file, err := os.OpenFile("go.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = gif.EncodeAll(file, &anim)
	if err != nil {
		log.Fatal(err)
	}
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Move x,y along given angle with given distance.
func move(x float64, y float64, angle float64, distance float64) (float64, float64) {
	return x + math.Cos(degToRad(angle))*distance, y + math.Sin(degToRad((angle)))*distance
}

// Check if given x,y coordinates are within canvas bounds
func isOutOfBounds(x float64, y float64, width float64, height float64) bool {
	return (x <= 0 || x >= width-1 || y <= 0 || y >= height-1)
}

func blur(grid *[width * height]float64) [width * height]float64 {
	var newGrid [width * height]float64
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			newGrid[idx(x, y)] = (grid[idx(x, y)] + grid[idx(x-1, y)] + grid[idx(x+1, y)] + grid[idx(x, y-1)] + grid[idx(x, y+1)] + grid[idx(x-1, y-1)] + grid[idx(x-1, y+1)] + grid[idx(x+1, y-1)] + grid[idx(x+1, y+1)]) / 9
		}
	}
	return newGrid
}

func decay(grid *[width * height]float64) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grid[idx(x, y)] *= 0.9
		}
	}
}

func idx(x int, y int) int {
	return y*width + x
}

func randSymmetricRange(boundary float64) float64 {
	return (rand.Float64() - 0.5) * 2 * boundary
}

func CreateImage(agents []Agent, colorPalette []color.Color) *image.Paletted {
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, colorPalette)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// c := 0
			img.SetColorIndex(x, y, 0) // black canvas
		}
	}
	for _, agent := range agents {
		img.SetColorIndex(int(math.Floor(agent.y)), int(math.Floor(agent.x)), uint8(254))
	}
	return img
}

func GetColorPalette(colors int) []color.Color {
	palette := []color.Color{color.Gray{0}}
	colors--
	for i, delta := 1, 255/colors; i < colors; i++ {
		palette = append(palette, color.Gray{uint8(delta * i)})
	}
	palette = append(palette, color.Gray{255})
	return palette
}
