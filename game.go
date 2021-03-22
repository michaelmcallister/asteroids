package asteroids

import (
	"math"

	"github.com/michaelmcallister/asteroids/internal/asteroid"
	"github.com/michaelmcallister/asteroids/internal/player"
	"github.com/michaelmcallister/asteroids/internal/vector"

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

var translation = vector.V2D{X: -ScreenWidth / 2, Y: -ScreenHeight / 2}

type Game struct {
	player    *player.Player
	asteroids []*asteroid.Asteroid
	score     int
	level     int
	started   bool
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
	for i, asteroid := range g.asteroids {
		if count == childAsteroidsToSpawn {
			return
		}
		if !asteroid.Alive {
			g.asteroids[i] = a.SpawnChild()
			g.asteroids[i].Alive = true
			g.asteroids[i].Update()
			count++
		}
	}
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

	if !g.started {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.started = true
		}
		return nil
	}

	g.handleInput()

	// Check for collisions.
	bulletLocations := make(map[*player.Bullet](vector.V2D))
	for _, bullet := range g.player.Bullets {
		if !bullet.Alive {
			continue
		}
		//convert bullet screen space location to world space to compare
		//with asteroids world space to detect a collision
		bulletLocations[bullet] = bullet.Location.AddNew(translation)
	}
	aliveAsteroids := 0
	for _, astrd := range g.asteroids {
		if !astrd.Alive {
			continue
		}
		aliveAsteroids++
		// Check if player collided with asteroid.
		if astrd.Collision(g.player.Location, g.player.HitRadius) {
			g.player.Kill()
		}

		for blt, loc := range bulletLocations {
			if astrd.Collision(loc, 1) {
				// We've collided, destroy the bullet.
				blt.Alive = false

				// Calculate the score.
				g.score += scores[astrd.Size]

				// Add new life every additionalLifeScore.
				if g.score%additionalLifeScore == 0 {
					g.player.Lives++
				}

				// If the asteroid is the smallest it can get, no more children.
				if astrd.Size == asteroid.Small {
					astrd.Alive = false
					continue
				} else {
					g.spawnAsteroids(astrd)
				}
			}
		}
		astrd.Update()
	}

	if aliveAsteroids == 0 {
		g.level++
		for i := 0; i < startingAsteroidCount+g.level; i++ {
			g.asteroids[i] = asteroid.New(ScreenWidth, ScreenHeight)
			g.asteroids[i].Alive = true
		}
	}

	g.player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.started {
		drawStart(screen)
		return
	}
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
