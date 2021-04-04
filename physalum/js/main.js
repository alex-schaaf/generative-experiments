let agents = []

// Convert degree to radians
function degToRad(deg) {
  return deg * (Math.PI / 180)
}

// Get new x,y moving along given angle with given distance.
function move(x, y, angle, distance) {
  return {
    x: x + Math.cos(degToRad(angle)) * distance,
    y: y + Math.sin(degToRad(angle)) * distance
  }
}

// Check if given x,y coordinates are within canvas bounds
function isOutOfBounds(x, y) {
  return (x <= 0 || x >= width - 1 || y <= 0 || y >= height - 1)
}

function setup() {
  let nCells = 80
  let size = 800
  // size -= size % nCells
  let cellSize = size / nCells
  let canvas = createCanvas(size, size)
  canvas.parent("canvas-container")

  grid = new Grid(nCells, nCells, cellSize)
  for (let i = 0; i < 10; i++) {
    agents.push(new Agent(size / 2, size / 2))
  }
  noStroke()
  frameRate(60)
}

/**
 * Agent class
 */
class Agent {
  constructor(x, y) {
    this.x = x
    this.y = y
    this.direction = random(360)

    this.sensorAngle = 25
    this.sensorDistance = 9
    this.FL = 0
    this.F = 0
    this.FR = 0

    draw()
  }


  // Move agent. If movement step would be out of bounds the agent direction
  // will be perturbed randomly and the function will be called recursively.
  move() {
    let { x, y } = move(this.x, this.y, this.direction, 5)
    if (!isOutOfBounds(x, y)) {
      this.x = x
      this.y = y
    } else {
      this.perturbRandomly(180)
      this.move()
    }

  }

  // Read front left, front and front right sensors.
  readSensors(grid) {
    this.FL = this.getSensorReading((this.direction - this.sensorAngle) % 360, grid)
    this.F = this.getSensorReading(this.direction, grid)
    this.FR = this.getSensorReading((this.direction + this.sensorAngle) % 360, grid)
  }

  // Get a reading from the grid at sensor distance from current agent position
  // along given angle
  getSensorReading(angle, grid) {
    let { x, y } = move(this.x, this.y, angle, this.sensorDistance)
    if (isOutOfBounds(x, y)) {
      return 0
    }
    return grid.getValue(x, y)
  }

  perturbRandomly(rng) {
    this.direction = (this.direction + random(-rng, rng)) % 360
  }

  perturb() {
    if (this.F > this.FL && this.F > this.FR) {
      return
    } else if (this.F < this.FL && this.F < this.FR) {
      if (random() < 0.5) {
        // rotate left
        this.direction = (this.direction - this.sensorAngle) % 360
      } else {
        // rotate right
        this.direction = (this.direction + this.sensorAngle) % 360
      }
    } else if (this.FL < this.FR) {
      // rotate right by sensor angle
      this.direction = (this.direction + this.sensorAngle) % 360
    } else if (this.FR < this.FL) {
      // rotate left by sensor angle
      this.direction = (this.direction - this.sensorAngle) % 360
    } else {
      // don't rotate
      return
    }
  }

  // Draw agent as a circle
  draw() {
    fill("red")
    circle(this.x, this.y, 5)
  }
}

/**
 * Grid class
 */
class Grid {
  constructor(rows, cols, cellSize) {
    this.rows = rows
    this.cols = cols
    this.cellSize = cellSize
    this.grid = []

    for (let r = 0; r < this.rows; r++) {
      let row = []
      for (let c = 0; c < this.cols; c++) {
        row.push(0)
      }
      this.grid.push(row)
    }
  }

  draw() {
    fill("black")
    for (let r = 0; r < this.rows; r++) {
      for (let c = 0; c < this.cols; c++) {
        let value = this.grid[r][c]
        fill(value)
        rect(r * this.cellSize, c * this.cellSize, this.cellSize, this.cellSize)
      }
    }
  }

  xyToIdx(x, y) {
    return { x: Math.floor(x / this.cellSize), y: Math.floor(y / this.cellSize) }
  }

  outOfBounds(x, y) {
    let { gridX, gridY } = this.xyToIdx(x, y)
    return (gridX < 0 || gridX >= this.cols || gridY < 0 || gridY >= this.rows)
  }

  getValue(xCoord, yCoord) {
    let { x, y } = this.xyToIdx(xCoord, yCoord)
    // console.log(`${xCoord} -> ${x}; ${yCoord} -> ${y}`)
    return this.grid[x][y]
  }
}


function draw() {
  // background(255);
  grid.draw()
  agents.forEach(agent => {
    agent.readSensors(grid)
    agent.perturb()
    agent.move()
    deposit(agent, grid)
    blur(grid)
    // decay(grid)
    // agent.draw()
  })
}

function deposit(agent, grid) {
  if (grid.outOfBounds(agent.x, agent.y)) {
    return
  }
  let { x, y } = grid.xyToIdx(agent.x, agent.y)
  grid.grid[x][y] += 5
}

// Convolution of an average-blur kernel
function blur(grid) {

}

// decay the grid
function decay(grid) {
  for (let r = 0; r < grid.rows; r++) {
    for (let c = 0; c < grid.cols; c++) {
      if (grid.grid[r, c] > 0) {
        grid.grid[r][c] -= 5
      }
    }
  }
}

