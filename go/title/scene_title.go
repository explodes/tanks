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

	layers := tempura.NewLayers(numLayers)

	titleScene := &titleScene{
		g:      game,
		layers: layers,
	}

	layers[layerMenu].Add(&tempura.Object{
		Tag:      "tanks",
		Drawable: tempura.NewImageDrawable(tanksImage),
		Pos:      tempura.V(core.ScreenWidth*0.5-float64(tanksImage.Bounds().Dx())*0.25, core.ScreenHeight*0.5-float64(tanksImage.Bounds().Dy())*0.25),
		Size:     tempura.V(float64(tanksImage.Bounds().Dx())*0.5, float64(tanksImage.Bounds().Dy())*0.5),
	})

	return titleScene, nil
}

func (s *titleScene) Update(dt float64) error {
	if obj := s.firstTouch(); obj != nil {
		return &core.ChangeGameError{Game: obj.Tag}
	}
	return nil
}

func (s *titleScene) Draw(image *ebiten.Image) {
	s.layers.Draw(image)
}

func objectBoundsContainsPoint(obj *tempura.Object, x, y float64) bool {
	return obj.Pos.X <= x &&
		obj.Pos.X+obj.Size.X >= x &&
		obj.Pos.Y <= y &&
		obj.Pos.Y+obj.Size.Y >= y
}

func (s *titleScene) firstTouch() *tempura.Object {
	for _, touch := range ebiten.Touches() {
		x, y := touch.Position()
		if obj := s.firstTouchAt(float64(x), float64(y)); obj != nil {
			return obj
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return s.firstTouchAt(float64(x), float64(y))
	}
	return nil
}

func (s *titleScene) firstTouchAt(x, y float64) *tempura.Object {
	xf, yf := float64(x), float64(y)
	core.DebugLog("click at %f, %f", x, y)
	iter := s.layers.IteratorTop()
	for obj, ok := iter(); ok; obj, ok = iter() {
		if objectBoundsContainsPoint(obj, xf, yf) {
			core.DebugLog("clicked object %v", obj)
			return obj
		}
	}
	return nil
}
