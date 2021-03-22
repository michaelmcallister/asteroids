package player

import (
	"image/color"

	"github.com/michaelmcallister/asteroids/internal/vector"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	maxBullets       = 3
	playerVectors    = 3
	bulletVelocity   = 4.1
	thrustMultiplier = .6
)

type Player struct {
	HitRadius          float64
	Lives              int
	Location, velocity vector.V2D
	objVert, worldVert [playerVectors]vector.V2D
	Bullets            [maxBullets]*Bullet
	translation        vector.V2D
}

type Bullet struct {
	Location, velocity vector.V2D
	Alive              bool
}

var (
	bulletImg = ebiten.NewImage(1, 1)
)

func init() {
	bulletImg.Fill(color.White)
}

func New(screenWidth, screenHeight int) *Player {
	p := &Player{}
	p.Lives = 3
	p.HitRadius = 15
	p.objVert[0] = vector.V2D{X: 0, Y: 1.5}
	p.objVert[1] = vector.V2D{X: -1, Y: -1}
	p.objVert[2] = vector.V2D{X: 1, Y: -1}
	p.translation = vector.V2D{
		X: float64(screenWidth / 2),
		Y: float64(screenHeight / 2),
	}

	p.Scale(12)
	p.Rotate(180)
	for i := range p.objVert {
		p.worldVert[i] = p.objVert[i].AddNew(p.translation)
	}
	for i := range p.Bullets {
		p.Bullets[i] = new(Bullet)
	}
	return p
}

func (p *Player) Scale(n float64) {
	for j := range p.objVert {
		if n < 0 {
			p.objVert[j].Divide(-n)
		}
		if n > 0 {
			p.objVert[j].Multiply(n)
		}
	}
}

func (p *Player) Kill() {
	p.velocity = vector.V2D{X: 0, Y: 0}
	p.Location = vector.V2D{X: 0, Y: 0}
	if p.Lives > 0 {
		p.Lives--
	}
}

func (p *Player) Shoot() {
	for _, b := range p.Bullets {
		if !b.Alive {
			b.Alive = true
			b.Location = p.worldVert[0]
			b.velocity = p.GetDirection()
			b.velocity.Multiply(bulletVelocity)
			break
		}
	}
}

func (p *Player) Rotate(degrees float64) {
	for i := range p.objVert {
		p.objVert[i].Rotate(degrees)
	}
}

func (p *Player) Update() {
	p.velocity.Limit(2)
	p.Location.Add(p.velocity)
	for i := range p.worldVert {
		t := p.objVert[i].AddNew(p.Location)
		t.Add(p.translation)
		p.worldVert[i] = t
	}
	for _, b := range p.Bullets {
		b.Location.Add(b.velocity)
	}
	p.bounds()
}

func (p *Player) Thrust() {
	p.velocity.Add(p.GetDirection())
	p.velocity.Multiply(thrustMultiplier)
}

func (p *Player) GetDirection() vector.V2D {
	direction := p.objVert[0]
	direction.Normalize()
	return direction
}

func (p *Player) bounds() {
	if p.Location.X < -p.translation.X {
		p.Location.X = p.translation.X
	}

	if p.Location.X > p.translation.X {
		p.Location.X = -p.translation.X
	}

	if p.Location.Y < -p.translation.Y {
		p.Location.Y = p.translation.Y
	}

	if p.Location.Y > p.translation.Y {
		p.Location.Y = -p.translation.Y
	}

	for _, b := range p.Bullets {
		if b.Location.X < 0 || b.Location.X >= (p.translation.X*2) {
			b.Alive = false
		}

		if b.Location.Y < 0 || b.Location.Y >= (p.translation.Y*2) {
			b.Alive = false
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, p.worldVert[0].X, p.worldVert[0].Y, p.worldVert[1].X, p.worldVert[1].Y, color.White)
	ebitenutil.DrawLine(screen, p.worldVert[1].X, p.worldVert[1].Y, p.worldVert[2].X, p.worldVert[2].Y, color.White)
	ebitenutil.DrawLine(screen, p.worldVert[2].X, p.worldVert[2].Y, p.worldVert[0].X, p.worldVert[0].Y, color.White)

	for _, b := range p.Bullets {
		if b.Alive {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(b.Location.X, b.Location.Y)
			screen.DrawImage(bulletImg, op)
		}
	}
}
