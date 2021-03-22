package asteroids

import (
	_ "embed"
	"fmt"
	"image/color"

	"github.com/michaelmcallister/asteroids/internal/player"
	"github.com/michaelmcallister/asteroids/internal/vector"

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
	gameOverSubText gameText = "Press <Enter> to restart"
)

var (
	gameStartControls []gameText = []gameText{
		"Controls",
		"Shoot.....<Space>",
		"Thrust....<W> or <UP>",
		"Move.....<A>/<D> or <LEFT>/<RIGHT>",
	}
)

const (
	gameStartTitle      gameText = "Asteroids"
	gameStartSubText    gameText = "Press <Enter>"
	gameStartBottomText gameText = "Michael McAllister 2021"
)

var (
	font8x  font.Face
	font16x font.Face
	font64x font.Face
)

var gameFont *truetype.Font

func init() {
	gameFont, _ = truetype.Parse(pixeboy_ttf)
	font8x = truetype.NewFace(gameFont, &truetype.Options{Size: 8.0})
	font16x = truetype.NewFace(gameFont, &truetype.Options{Size: 16.0})
	font64x = truetype.NewFace(gameFont, &truetype.Options{Size: 64.0})
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
		}
		n.Location.Add(translation)
		n.Update()
		n.Draw(screen)
		offset += n.HitRadius
	}
}

func drawStart(screen *ebiten.Image) {
	sHalf := ScreenWidth / 2
	sThird := ScreenHeight / 3

	b1 := text.BoundString(font64x, string(gameStartTitle))
	b2 := text.BoundString(font16x, string(gameStartSubText))
	b4 := text.BoundString(font8x, string(gameStartBottomText))

	topX := sHalf - (b1.Dx() / 2)
	topY := sThird
	subX := ScreenWidth/2 - b2.Dx()/2
	subY := sThird + b1.Dy() + ScreenHeight/4
	bottomX := ScreenWidth/2 - b4.Dx()/2

	text.Draw(screen, string(gameStartTitle), font64x, topX, topY, color.White)
	text.Draw(screen, string(gameStartSubText), font16x, subX, topY+b1.Dy(), color.White)
	controlsY := sThird + b1.Dy() + ScreenHeight/4
	maxX := 0
	for _, t := range gameStartControls {
		b := text.BoundString(font16x, string(t))
		if b.Dx() > maxX {
			maxX = b.Dx()
		}
	}
	text.Draw(screen, string(gameStartControls[0]), font16x, maxX/2, subY, color.White)
	for i := 1; i < len(gameStartControls); i++ {
		t := gameStartControls[i]
		b := text.BoundString(font16x, string(t))
		controlsY += b.Dy()
		text.Draw(screen, string(t), font16x, ScreenWidth/8, controlsY, color.White)
	}
	text.Draw(screen, string(gameStartBottomText), font8x, bottomX, ScreenHeight, color.White)
}

func drawGameOver(screen *ebiten.Image) {
	sHalf := ScreenWidth / 2
	sThird := ScreenHeight / 3

	b1 := text.BoundString(font64x, string(gameOverText))
	b2 := text.BoundString(font16x, string(gameOverSubText))
	topX := sHalf - (b1.Dx() / 2)
	topY := sThird
	subX := ScreenWidth/2 - b2.Dx()/2
	subY := sThird + b1.Dy() + sThird

	text.Draw(screen, string(gameOverText), font64x, topX, topY, color.White)
	text.Draw(screen, string(gameOverSubText), font16x, subX, subY, color.White)
}

func drawScore(screen *ebiten.Image, score int) {
	t := fmt.Sprintf("%06d", score)
	b1 := text.BoundString(font16x, t)
	xPos := b1.Dx() - b1.Dx()
	yPos := b1.Dy() + 10

	text.Draw(screen, t, font16x, xPos, yPos, color.White)
}
