let cell = 2  // cell size
let walkers = []


function setup() {
  let size = min(800, 1000)
  let canvas = createCanvas(size, size)
  canvas.parent("canvas-container")

  // let startX = width / 2


  for (let j = 0; j < 4; j++) {
    let potato = 250
    let startX = [potato, potato, width - potato, width - potato]
    let startY = [potato, height - potato, potato, height - potato]
    for (let i = 0; i < 1000; i++) {
      walkers.push(new Walker(startX[j], startY[j]))
    }
  }

  background(255)
  strokeWeight(2)
  draw()
}

function restart() {
  walkers = []  // delete existing walkers using garbage collection
  clear()  // clear canvas
  setup()  // rerun setup script
}

class Walker {
  constructor(x, y) {
    this.x = x
    this.px = x
    this.velocityX = random(-5, 5);

    this.y = y
    this.py = y
    this.velocityY = random(-5, 5);

    this.r = 1 || random(255)
    this.g = 1 || random(255)
    this.b = 1 || random(255)
    this.alpha = 5

    this.dampening = 0.2

    this.draw()
  }

  isOut() {
    return (this.x < 0 || this.x > width || this.y < 0 || this.y > height)
  }

  move() {
    this.x += this.dampening * this.velocityX;
    this.y += this.dampening * this.velocityY;
  }

  draw() {
    stroke(this.r, this.g, this.b, this.alpha);
    line(this.x, this.y, this.px, this.py)
    this.px = this.x
    this.py = this.y
  }

  updateVelocity() {
    this.velocityX += random(-0.9, 0.9)
    this.velocityY += random(-0.9, 0.9)
  }
  updateVelocityPerlin() {
    let locator = 0.009
    this.velocityX += map(noise(this.x * locator, this.y * locator), 0, 1, -1, 1) //+ random() *locator
    this.velocityY += map(noise(this.y * locator, this.x * locator), 0, 1, -1, 1) //+ random() *locator
  }

}

function draw() {
  walkers.forEach(walker => {
    if (!walker.isOut()) {
      walker.updateVelocityPerlin()
      walker.move()
      walker.draw()
    }
  })
}