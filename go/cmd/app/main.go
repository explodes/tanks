package main

import (
	"log"

	"github.com/explodes/tanks/go/tanks"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	game, err := tanks.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.Run(game.Update, tanks.ScreenWidth, tanks.ScreenHeight, 1, tanks.Title); err != nil {
		log.Fatal(err)
	}
}
