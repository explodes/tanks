// +build darwin freebsd linux windows js
// +build !android
// +build !ios

package tanks

import (
	"github.com/hajimehoshi/ebiten"
)

var _ Input = (*inputImpl)(nil)

type inputImpl struct {
}

func NewInput() Input {
	return &inputImpl{}
}

func (i *inputImpl) Begin() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

func (i *inputImpl) BlueRotate() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA)
}

func (i *inputImpl) RedRotate() bool {
	return ebiten.IsKeyPressed(ebiten.KeyL)
}
