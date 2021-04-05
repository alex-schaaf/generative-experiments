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
	"time"
)

const width = 640
const height = 640
const sensorAngle = 45
const sensorDistance = 9

var iter = 100

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

// move the agent by 1 (if step would be within bounds)
// else perturb direction randomly and try to move again
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

// perturbRandomly shuffles the direction of an agent by adding a
// random number between [-angle, angle]
func (a *Agent) perturbRandomly(angle float64) {
	a.direction = float64(int(a.direction+randSymmetricRange(angle)) % 360)
}

// perturb takes care of modifying the direction of the agent in line
// with Jones (2010) paper to follow pheromone trails on the underlying grid
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

// readSensors reads out the front left (FL), front (F) and front right (FR)
// sensors of an agent, will read 0 if sensor position is out of bounds
func (a *Agent) readSensors(grid [width * height]float64) {
	a.fl = a.getSensorReading(float64(int(a.direction-a.sensorAngle)%360), grid)
	a.f = a.getSensorReading(a.direction, grid)
	a.fr = a.getSensorReading(float64(int(a.direction+a.sensorAngle)%360), grid)
}

// getSensorReading at given angle in agent.sensorDistance from the grid
// return 0 if sensor position would be out of bounds
func (a *Agent) getSensorReading(angle float64, grid [width * height]float64) float64 {
	x, y := move(a.x, a.y, angle, a.sensorDistance)
	if isOutOfBounds(x, y, width, height) {
		return 0
	}
	return grid[int(math.Floor(y))*width+int(math.Floor(x))]
}

// deposit a given amount of pheromone trail onto the grid at the agent's position
func (a *Agent) deposit(grid *[width * height]float64, amount float64) {
	grid[int(math.Floor(a.y))*width+int(math.Floor(a.x))] += amount
}

func main() {
	var grid [width * height]float64
	colorPalette := GetColorPalette(255)
	var images []*image.Paletted

	t0 := time.Now()

	var agents []*Agent
	for i := 0; i < 100; i++ {
		agent := &Agent{
			x:              rand.Float64() * width,
			y:              rand.Float64() * height,
			direction:      rand.Float64() * 360,
			sensorAngle:    sensorAngle,
			sensorDistance: sensorDistance,
			fl:             0,
			f:              0,
			fr:             0}
		agents = append(agents, agent)
	}

	// run simulation
	var delays []int
	for iter > 0 {
		for _, agent := range agents {
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
		delays = append(delays, 8)
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

	t1 := time.Now()

	td := t1.Second() - t0.Second()
	fmt.Printf("This took %d Seconds.", td)
}

// degToRad converts degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Move x,y along given angle with given distance.
func move(x float64, y float64, angle float64, distance float64) (float64, float64) {
	return x + math.Cos(degToRad(angle))*distance, y + math.Sin(degToRad(angle))*distance
}

// Check if given x,y coordinates are within canvas bounds
func isOutOfBounds(x float64, y float64, width float64, height float64) bool {
	return x <= 0 || x >= width-1 || y <= 0 || y >= height-1
}

// blur every pixel of the grid with an average of all 8 neighbors and itself
func blur(grid *[width * height]float64) [width * height]float64 {
	var newGrid [width * height]float64
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			newGrid[idx(x, y)] = (grid[idx(x, y)] + grid[idx(x-1, y)] + grid[idx(x+1, y)] + grid[idx(x, y-1)] + grid[idx(x, y+1)] + grid[idx(x-1, y-1)] + grid[idx(x-1, y+1)] + grid[idx(x+1, y-1)] + grid[idx(x+1, y+1)]) / 9
		}
	}
	return newGrid
}

// decay the grid pheromone values by a factor of 0.9
func decay(grid *[width * height]float64) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grid[idx(x, y)] *= 0.9
		}
	}
}

// idx converts x,y coordinates into the corresponding flattened array index
func idx(x int, y int) int {
	return y*width + x
}

// randSymmetricRange returns a random float in [-boundary, boundary]
func randSymmetricRange(boundary float64) float64 {
	return (rand.Float64() - 0.5) * 2 * boundary
}

// CreateImage creates a new image with a black background and white dots for each agent
func CreateImage(agents []*Agent, colorPalette []color.Color) *image.Paletted {
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, colorPalette)
	for _, agent := range agents {
		img.SetColorIndex(int(math.Floor(agent.y)), int(math.Floor(agent.x)), uint8(254))
	}
	return img
}

// GetColorPalette creates a grayscale color palette with nColors steps
func GetColorPalette(nColos int) []color.Color {
	palette := []color.Color{color.Gray{}}
	nColos--
	for i, delta := 1, 255/nColos; i < nColos; i++ {
		palette = append(palette, color.Gray{Y: uint8(delta * i)})
	}
	palette = append(palette, color.Gray{Y: 255})
	return palette
}
