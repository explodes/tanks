package title

import (
	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tempura"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var _ core.Scene = (*titleScene)(nil)

const (
	layerBackground = iota
	layerMenu
	layerForeground
	numLayers
)

type titleScene struct {
	g *Game

	layers tempura.Layers
}

func NewTitleScene(game *Game) (core.Scene, error) {
	loader := game.context.Loader()

	tanksImage, err := loader.EbitenImage("images/ic_tanks.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	shipwreckImage, err := loader.EbitenImage("images/ship_blue.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	layers := tempura.NewLayers(numLayers)

	titleScene := &titleScene{
		g:      game,
		layers: layers,
	}

	menu := []struct {
		name  string
		image *ebiten.Image
	}{
		{"tanks", tanksImage},
		{"shipwreck", shipwreckImage},
	}

	w := core.ScreenWidth / float64(len(menu)) * 0.75
	h := w

	dx := w * 1.10

	x := (core.ScreenWidth - (float64(len(menu)) * dx)) * 0.5
	y := core.ScreenHeight*0.5 - h*0.5

	for _, m := range menu {
		layers[layerMenu].Add(&tempura.Object{
			Tag:      m.name,
			Drawable: tempura.NewImageDrawable(m.image),
			Pos:      tempura.V(x, y),
			Size:     tempura.V(w, h),
		})
		x += dx
	}

	return titleScene, nil
}

func (s *titleScene) Update(dt float64) error {
	if obj := s.menuTouch(); obj != nil {
		return &core.ChangeGameError{Game: obj.Tag}
	}
	return nil
}

func (s *titleScene) Draw(image *ebiten.Image) {
	s.layers.Draw(nil, image)
}

func objectBoundsContainsPoint(obj *tempura.Object, x, y float64) bool {
	return obj.Pos.X <= x &&
		obj.Pos.X+obj.Size.X >= x &&
		obj.Pos.Y <= y &&
		obj.Pos.Y+obj.Size.Y >= y
}

func (s *titleScene) menuTouch() *tempura.Object {
	for _, touch := range ebiten.Touches() {
		x, y := touch.Position()
		if obj := s.menuTouchAt(float64(x), float64(y)); obj != nil {
			return obj
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return s.menuTouchAt(float64(x), float64(y))
	}
	return nil
}

func (s *titleScene) menuTouchAt(x, y float64) *tempura.Object {
	xf, yf := float64(x), float64(y)
	iter := s.layers[layerMenu].Iterator()
	for obj, ok := iter(); ok; obj, ok = iter() {
		if objectBoundsContainsPoint(obj, xf, yf) {
			return obj
		}
	}
	return nil
}
