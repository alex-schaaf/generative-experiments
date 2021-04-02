let particles = []

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
  for (let i = 0; i < 5; i++) {
    particles.push(new Particle(random(size), random(size)))
  }
  noStroke()
  frameRate(24)
}

/**
 * Particle class
 */
class Particle {
  constructor(x, y) {
    this.x = x
    this.y = y
    this.direction = random(360)

    this.sensorAngle = 45
    this.sensorDistance = 9
    this.FL = 0
    this.F = 0
    this.FR = 0

    draw()
  }

  move() {
    let { x, y } = move(this.x, this.y, this.direction, 5)
    if (!isOutOfBounds(x, y)) {
      this.x = x
      this.y = y
    } else {
      this.perturb(180)
      this.move()
    }

  }

  getNextXY(degrees, step) {
    return {
      x: this.x + Math.cos(degToRad(degrees)) * step,
      y: this.y + Math.sin(degToRad(degrees)) * step
    }
  }

  readSensors(grid) {
    this.FL = this.getSensorReading((this.direction - this.sensorAngle) % 360, grid)
    this.F = this.getSensorReading(this.direction, grid)
    this.FR = this.getSensorReading((this.direction + this.sensorAngle) % 360, grid)
  }

  getSensorReading(angle, grid) {
    let { x, y } = move(this.x, this.y, angle, this.sensorDistance)
    if (isOutOfBounds(x, y)) {
      return 0
    }
    return grid.getValue(x, y)
  }

  perturb(rng) {
    this.direction = (this.direction + random(-rng, rng)) % 360
  }

  draw() {
    fill("red")
    circle(this.x, this.y, 5)
    // fill("white")
    // circle(
    //   this.x + Math.cos(degToRad(this.direction)) * 9,
    //   this.y + Math.sin(degToRad(this.direction)) * 9,
    //   5, 5
    // )
    // circle(
    //   this.x + Math.cos(degToRad(this.direction - 45)) * 9,
    //   this.y + Math.sin(degToRad((this.direction - 45) % 360)) * 9,
    //   5, 5
    // )
    // circle(
    //   this.x + Math.cos(degToRad(this.direction + 45)) * 9,
    //   this.y + Math.sin(degToRad((this.direction + 45) % 360)) * 9,
    //   5, 5
    // )
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
    // stroke(25, 25, 25, 255)
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
  particles.forEach(particle => {
    particle.readSensors(grid)
    particle.perturb(25)
    particle.move()
    deposit(particle, grid)
    particle.draw()
  })
}

function deposit(p, grid) {
  if (grid.outOfBounds(p.x, p.y)) {
    return
  }
  let { x, y } = grid.xyToIdx(p.x, p.y)
  grid.grid[x][y] += 15
}