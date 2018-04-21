package tanks

import (
	"github.com/explodes/tanks/go/core"
	"github.com/hajimehoshi/ebiten"
)

func Begin() bool {
	return len(ebiten.Touches()) > 0
}

func BlueRotate() bool {
	for _, touch := range ebiten.Touches() {
		x, _ := touch.Position()
		if x < core.ScreenWidth/2 {
			return true
		}
	}
	return false
}

func RedRotate() bool {
	for _, touch := range ebiten.Touches() {
		x, _ := touch.Position()
		if x > core.ScreenWidth/2 {
			return true
		}
	}
	return false
}
