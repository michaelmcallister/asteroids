package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/michaelmcallister/asteroids"
)

func main() {
	ebiten.SetWindowSize(asteroids.ScreenWidth, asteroids.ScreenHeight)
	ebiten.SetWindowTitle("Asteroids")
	p := asteroids.NewGame()
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}
