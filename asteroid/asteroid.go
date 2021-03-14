package asteroid

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/michaelmcallister/asteroids/vector"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const asteroidVertices = 10

type Size int

const (
	Small Size = iota
	Medium
	Large
)

type Asteroid struct {
	Alive                           bool
	HitRadius                       float64
	Size                            Size
	rotation                        float64
	Location, velocity, translation vector.V2D
	objVert, worldVert              [asteroidVertices]vector.V2D
}

func New(screenWidth, screenHeight int) *Asteroid {
	rand.Seed(time.Now().UnixNano())
	asteroid := new(Asteroid)
	asteroid.translation = vector.V2D{X: float64(screenWidth / 2), Y: float64(screenHeight / 2)}

	// Start asteroid in random location, avoiding the centre.
	lx := rand.Intn((screenWidth / 2) + 35)
	ly := rand.Intn((screenHeight / 2) + 35)

	//give asteroid random velocity
	vx := rand.Float64() + 0.1/2
	vy := rand.Float64() + 0.1/2

	degrees := rand.Float64() + 0.1

	//50% chance the sign of the variable will be switched to negative
	if rand.Intn(1) == 1 {
		vx = -vx
		lx = -lx
		degrees = -degrees
	}

	//50% chance the sign of the variable will be switched to negative
	if rand.Intn(1) == 1 {
		vy = -vy
		ly = -ly
	}

	asteroid.Size = Large
	asteroid.HitRadius = 35
	asteroid.rotation = degrees
	asteroid.Location = vector.V2D{X: float64(lx), Y: float64(ly)}
	asteroid.velocity = vector.V2D{X: vx, Y: vy}
	asteroid.objVert[0] = vector.V2D{X: 0, Y: .4}
	asteroid.objVert[1] = vector.V2D{X: .2, Y: .3}
	asteroid.objVert[2] = vector.V2D{X: .2, Y: .1}
	asteroid.objVert[3] = vector.V2D{X: .4, Y: 0}
	asteroid.objVert[4] = vector.V2D{X: .3, Y: -.2}
	asteroid.objVert[5] = vector.V2D{X: .1, Y: -.2}
	asteroid.objVert[6] = vector.V2D{X: 0, Y: -0.3}
	asteroid.objVert[7] = vector.V2D{X: -.2, Y: -0.2}
	asteroid.objVert[8] = vector.V2D{X: -.4, Y: -0}
	asteroid.objVert[9] = vector.V2D{X: -.3, Y: .3}

	asteroid.Scale(88)
	for i := range asteroid.objVert {
		// converts verts from obj space to world space and translate world '
		// space to screen space.
		asteroid.worldVert[i] = asteroid.worldVert[i].Add(asteroid.objVert[i]).Add(asteroid.translation)
	}
	asteroid.Update()
	return asteroid
}

func (a *Asteroid) Update() {
	a.Location = a.Location.Add(a.velocity)
	//update each vert of the asteroid to reflect the changes made to the asteroids location vector
	//and rotation amount, then translate the new vert location to screen space
	for i := range a.objVert {
		a.worldVert[i] = a.objVert[i].Add(a.Location).Add(a.translation)
		a.objVert[i] = a.objVert[i].Rotate(a.rotation)
	}
	a.bounds()
}

func (a *Asteroid) bounds() {
	if a.Location.X < -a.translation.X {
		a.Location.X = a.translation.X
	}
	if a.Location.X > a.translation.X {
		a.Location.X = -a.translation.X
	}
	if a.Location.Y < -a.translation.Y {
		a.Location.Y = a.translation.Y
	}
	if a.Location.Y > a.translation.Y {
		a.Location.Y = -a.translation.Y
	}
}

func (a *Asteroid) Scale(n float64) {
	if n < 0 {
		for j := range a.objVert {
			a.objVert[j] = a.objVert[j].Divide(-n)
		}
	}
	if n > 0 {
		for j := range a.objVert {
			a.objVert[j] = a.objVert[j].Multiply(n)
		}
	}
}

func (a *Asteroid) Shrink() {
	switch a.Size {
	case Small:
		a.Alive = false
		return
	case Medium:
		a.Size = Small
	case Large:
		a.Size = Medium
	}
	a.HitRadius /= 2
	a.Scale(-2.0)
}

func (a *Asteroid) Collision(v vector.V2D, radius float64) bool {
	if !a.Alive || a == nil {
		return false
	}
	p1 := math.Pow(a.Location.X-v.X, 2)
	p2 := math.Pow(a.Location.Y-v.Y, 2)
	distance := math.Sqrt(p1 + p2)
	sum := a.HitRadius + radius
	return distance < sum
}

func (a *Asteroid) SpawnChild() *Asteroid {
	b := New(int(a.translation.X*2), int(a.translation.Y*2))
	b.HitRadius = a.HitRadius
	b.Location = a.Location
	b.objVert = a.objVert
	b.Size = a.Size
	b.Alive = true
	return b
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	if !a.Alive || a == nil {
		return
	}
	ebitenutil.DrawLine(screen, a.worldVert[0].X, a.worldVert[0].Y, a.worldVert[1].X, a.worldVert[1].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[1].X, a.worldVert[1].Y, a.worldVert[2].X, a.worldVert[2].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[2].X, a.worldVert[2].Y, a.worldVert[3].X, a.worldVert[3].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[3].X, a.worldVert[3].Y, a.worldVert[4].X, a.worldVert[4].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[4].X, a.worldVert[4].Y, a.worldVert[5].X, a.worldVert[5].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[5].X, a.worldVert[5].Y, a.worldVert[6].X, a.worldVert[6].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[6].X, a.worldVert[6].Y, a.worldVert[7].X, a.worldVert[7].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[7].X, a.worldVert[7].Y, a.worldVert[8].X, a.worldVert[8].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[8].X, a.worldVert[8].Y, a.worldVert[9].X, a.worldVert[9].Y, color.White)
	ebitenutil.DrawLine(screen, a.worldVert[9].X, a.worldVert[9].Y, a.worldVert[0].X, a.worldVert[0].Y, color.White)
}
