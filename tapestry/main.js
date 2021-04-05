function setup() {

  createCanvas(800, 800)
  background(255)
  // noSmooth()


  let previousXs = []
  for (let x = 0; x < 800; x++) {
    previousXs.push(10)
  }
  strokeWeight(0.5)
  stroke(45)
  for (let x = 0; x < 155; x++) {
    for (let y = 1; y < 800; y++) {
      let newX = previousXs[y] + 3 + random()
      line(previousXs[y - 1], y - 1, newX, y)
      previousXs[y] = newX

      if (y > 750) {
        if (random() < 0.1) {
          break
        }
      }
    }
    // if (x === 50) {
    //   for (let y = 500; y < 800; y++) {
    //     previousXs[y] += y / 60 * 10 * random()
    //   }
    // }

  }
}


