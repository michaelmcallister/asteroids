package asteroids

import (
	_ "embed"
	"fmt"
	"image/color"

	"github.com/michaelmcallister/asteroids/player"
	"github.com/michaelmcallister/asteroids/vector"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var (
	fontGameOver font.Face
	fontSubText  font.Face
	fontScore    font.Face
	//go:embed resources/Pixeboy.ttf
	pixeboy_ttf []byte
)

const (
	gameOverText    = "Game Over"
	gameOverSubText = "Press <enter> to restart"
)

func init() {
	fo, _ := truetype.Parse(pixeboy_ttf)
	fontGameOver = truetype.NewFace(fo, &truetype.Options{Size: 64.0})
	fontSubText = truetype.NewFace(fo, &truetype.Options{Size: 32.0})
	fontScore = truetype.NewFace(fo, &truetype.Options{Size: 32.0})
}

func drawLives(screen *ebiten.Image, count int) {
	offset := 10.0
	translation := vector.V2D{X: -ScreenWidth / 2, Y: -ScreenHeight / 2}
	for i := 0; i < count; i++ {
		n := player.New(ScreenWidth, ScreenHeight)
		n.Scale(-2.0)
		//convert screen space vector into world space
		n.Location = vector.V2D{
			X: ScreenWidth - offset,
			Y: n.HitRadius,
		}.Add(translation)
		n.Update()
		n.Draw(screen)
		offset += n.HitRadius
	}
}

func drawGameOver(screen *ebiten.Image) {
	sHalf := ScreenWidth / 2
	sThird := ScreenHeight / 3

	b1 := text.BoundString(fontGameOver, gameOverText)
	b2 := text.BoundString(fontSubText, gameOverSubText)
	topX := sHalf - (b1.Dx() / 2)
	topY := sThird
	subX := ScreenWidth/2 - b2.Dx()/2
	subY := sThird + b1.Dy() + sThird

	text.Draw(screen, gameOverText, fontGameOver, topX, topY, color.White)
	text.Draw(screen, gameOverSubText, fontSubText, subX, subY, color.White)
}

func drawScore(screen *ebiten.Image, score int) {
	t := fmt.Sprintf("%06d", score)
	b1 := text.BoundString(fontScore, t)
	xPos := b1.Dx() - b1.Dx()
	yPos := b1.Dy()

	text.Draw(screen, t, fontScore, xPos, yPos, color.White)
}
