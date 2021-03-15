package asteroids

import (
	"math"

	"github.com/michaelmcallister/asteroids/asteroid"
	"github.com/michaelmcallister/asteroids/player"
	"github.com/michaelmcallister/asteroids/vector"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const ScreenWidth, ScreenHeight = 640, 480

const (
	rotationDegrees       = 4
	startingAsteroidCount = 3
	childAsteroidsToSpawn = 2
	additionalLifeScore   = 10000
)

var scores = map[asteroid.Size]int{
	asteroid.Small:  100,
	asteroid.Medium: 50,
	asteroid.Large:  20,
}

type Game struct {
	player    *player.Player
	asteroids []*asteroid.Asteroid
	score     int
	level     int
}

func NewGame() *Game {
	g := &Game{}
	g.reset()
	return g
}

func (g *Game) spawnAsteroids(a *asteroid.Asteroid) {
	// Shrink the asteroid size and spawn children from it.
	a.Shrink()
	count := 0
	for _, asteroid := range g.asteroids {
		if count == childAsteroidsToSpawn {
			break
		}
		if !asteroid.Alive {
			g.asteroids = append(g.asteroids, a.SpawnChild())
			count++
		}
	}
}

func (g *Game) aliveAsteroids() int {
	c := 0
	for _, a := range g.asteroids {
		if a.Alive {
			c++
		}
	}
	return c
}

func (g *Game) reset() {
	g.player = player.New(ScreenWidth, ScreenHeight)
	g.level = 0
	g.score = 0
	g.asteroids = nil

	m := int(math.Pow(startingAsteroidCount, (childAsteroidsToSpawn + 1)))
	for i := 0; i < m; i++ {
		g.asteroids = append(g.asteroids, asteroid.New(ScreenWidth, ScreenHeight))
	}
	for i := 0; i < startingAsteroidCount; i++ {
		g.asteroids[i].Alive = true
	}
}

func (g *Game) handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Rotate(-rotationDegrees)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Rotate(rotationDegrees)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Thrust()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.player.Shoot()
	}
}

func (g *Game) Update() error {
	// If it's game over pause the game until Enter is pressed.
	if g.player.Lives == 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
		return nil
	}

	g.handleInput()

	// Killed all asteroids.
	if g.aliveAsteroids() == 0 {
		g.level++
		for i := 0; i < startingAsteroidCount+g.level; i++ {
			g.asteroids[i] = asteroid.New(ScreenWidth, ScreenHeight)
			g.asteroids[i].Alive = true
		}
	}

	// Check for collisions.
	for _, asteroid := range g.asteroids {
		if !asteroid.Alive {
			continue
		}
		// Asteroid collided with player.
		if asteroid.Collision(g.player.Location, g.player.HitRadius) {
			g.player.Kill()
		}
	}
	translation := vector.V2D{X: -ScreenWidth / 2, Y: -ScreenHeight / 2}
	for _, bullet := range g.player.Bullets {
		if !bullet.Alive {
			continue
		}
		//convert bullet screen space location to world space to compare
		//with asteroids world space to detect a collision
		world := bullet.Location.Add(translation)
		for _, a := range g.asteroids {
			// No collision, moving on.
			if !a.Collision(world, 1) {
				continue
			}
			// We've collided, destroy the bullet.
			bullet.Alive = false

			// Calculate the score.
			g.score += scores[a.Size]

			// Add new life every additionalLifeScore.
			if g.score%additionalLifeScore == 0 {
				g.player.Lives++
			}

			// If the asteroid is the smallest it can get, no more children.
			if a.Size == asteroid.Small {
				a.Alive = false
				continue
			} else {
				g.spawnAsteroids(a)
			}
		}
	}

	g.player.Update()
	for _, a := range g.asteroids {
		a.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.player.Lives == 0 {
		drawGameOver(screen)
	}
	g.player.Draw(screen)
	for _, a := range g.asteroids {
		a.Draw(screen)
	}
	drawLives(screen, g.player.Lives)
	drawScore(screen, g.score)
}

func (*Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}
