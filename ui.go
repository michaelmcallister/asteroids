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
	//go:embed resources/PressStart2P-Regular.ttf
	pixeboy_ttf []byte
)

type gameText string

const (
	gameOverText    gameText = "Game Over"
	gameOverSubText gameText = "Press <enter> to restart"
)

const (
	gameStartTitle      gameText = "Asteroids"
	gameStartSubText    gameText = "Press <enter> to play"
	gameStartBottomText gameText = "Michael McAllister 2021"
)

var fontSizes = map[gameText]float64{
	gameOverText:        64.0,
	gameOverSubText:     32.0,
	gameStartTitle:      64.0,
	gameStartSubText:    16.0,
	gameStartBottomText: 8.0,
}

var gameFont *truetype.Font

func init() {
	gameFont, _ = truetype.Parse(pixeboy_ttf)
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

func drawStart(screen *ebiten.Image) {
	sHalf := ScreenWidth / 2
	sThird := ScreenHeight / 3

	f1 := truetype.NewFace(gameFont, &truetype.Options{Size: fontSizes[gameStartTitle]})
	f2 := truetype.NewFace(gameFont, &truetype.Options{Size: fontSizes[gameStartSubText]})
	f3 := truetype.NewFace(gameFont, &truetype.Options{Size: fontSizes[gameStartBottomText]})
	b1 := text.BoundString(f1, string(gameStartTitle))
	b2 := text.BoundString(f2, string(gameStartSubText))
	b3 := text.BoundString(f3, string(gameStartBottomText))

	topX := sHalf - (b1.Dx() / 2)
	topY := sThird
	subX := ScreenWidth/2 - b2.Dx()/2
	subY := sThird + b1.Dy() + ScreenHeight/4
	bottomX := ScreenWidth/2 - b3.Dx()/2

	text.Draw(screen, string(gameStartTitle), f1, topX, topY, color.White)
	text.Draw(screen, string(gameStartSubText), f2, subX, subY, color.White)
	text.Draw(screen, string(gameStartBottomText), f3, bottomX, ScreenHeight, color.White)
}

func drawGameOver(screen *ebiten.Image) {
	sHalf := ScreenWidth / 2
	sThird := ScreenHeight / 3

	f1 := truetype.NewFace(gameFont, &truetype.Options{Size: fontSizes[gameOverText]})
	f2 := truetype.NewFace(gameFont, &truetype.Options{Size: fontSizes[gameStartSubText]})

	b1 := text.BoundString(f1, string(gameOverText))
	b2 := text.BoundString(f2, string(gameOverSubText))
	topX := sHalf - (b1.Dx() / 2)
	topY := sThird
	subX := ScreenWidth/2 - b2.Dx()/2
	subY := sThird + b1.Dy() + sThird

	text.Draw(screen, string(gameOverText), f1, topX, topY, color.White)
	text.Draw(screen, string(gameOverSubText), f2, subX, subY, color.White)
}

func drawScore(screen *ebiten.Image, score int) {
	t := fmt.Sprintf("%06d", score)
	f1 := truetype.NewFace(gameFont, &truetype.Options{Size: 16.0})
	b1 := text.BoundString(f1, t)
	xPos := b1.Dx() - b1.Dx()
	yPos := b1.Dy() + 10

	text.Draw(screen, t, f1, xPos, yPos, color.White)
}
