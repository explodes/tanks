package main

import (
	"log"

	_ "github.com/explodes/tanks/go/cmd/games_registry"
	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tanks/go/overworld"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	game, err := overworld.NewOverworld()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(game.Update, core.ScreenWidth, core.ScreenHeight, 1, core.Title); err != nil {
		log.Fatal(err)
	}
}
