package tanks

import "github.com/hajimehoshi/ebiten"

func (i *inputImpl) ToggleFullscreen() bool {
	return false
}

func (i *inputImpl) ToggleMute() bool {
	return false
}

func (i *inputImpl) Exit() bool {
	return false
}

func (i *inputImpl) Begin() bool {
	return len(ebiten.Touches()) > 0
}

func (i *inputImpl) BlueRotate() bool {
	for _, touch := range ebiten.Touches() {
		x, _ := touch.Position()
		if x < ScreenWidth/2 {
			return true
		}
	}
	return false
}

func (i *inputImpl) RedRotate() bool {
	for _, touch := range ebiten.Touches() {
		x, _ := touch.Position()
		if x > ScreenWidth/2 {
			return true
		}
	}
	return false
}
