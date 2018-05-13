package tanks

import (
	"github.com/hajimehoshi/ebiten"

	_ "image/jpeg"
	_ "image/png"

	"image"

	"image/color"
	"math/rand"

	"math"

	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tempura"
)

var _ core.Scene = (*gameScene)(nil)

const (
	shipPosPadding = 25
	shipSpeed      = 100
	shipW, shipH   = 50, 50

	fireRate = 0.5

	redRotNormal  = -90 * math.Pi / 180
	blueRotNormal = 0
)

const (
	layerBackground = iota
	layerShips
	layerDebris
	layerBullets
	numLayers
)

type gameScene struct {
	g    *Game
	time float64

	layers tempura.Layers
}

func NewGameScene(game *Game) (core.Scene, error) {
	loader := game.context.Loader()

	blueImage, err := loader.EbitenImage("images/ship_blue.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	redImage, err := loader.EbitenImage("images/ship_red.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	s := &gameScene{
		g:      game,
		layers: tempura.NewLayers(numLayers),
	}

	bg, err := newBackground()
	if err != nil {
		return nil, err
	}
	s.layers[layerBackground].Add(bg)

	s.layers[layerShips].Add(s.newShip(blueImage, "blue", 1, tempura.V(shipPosPadding, core.ScreenHeight*0.5), 0, 1, BlueSwitchDirection))
	s.layers[layerShips].Add(s.newShip(redImage, "red", -1, tempura.V(core.ScreenWidth-shipW-shipPosPadding, core.ScreenHeight*0.5), redRotNormal, -1, RedSwitchDirection))

	return s, nil
}

func (s *gameScene) Update(dt float64) error {
	s.time += dt
	s.layers.Update(dt)
	return nil
}

func (s *gameScene) Draw(image *ebiten.Image) {
	s.layers.Draw(nil, image)
}

func newBackground() (*tempura.Object, error) {
	const numStars = 1000

	bg := image.NewGray(image.Rect(0, 0, core.ScreenWidth, core.ScreenHeight))

	for i := 0; i < numStars; i++ {
		bg.SetGray(rand.Intn(core.ScreenWidth), rand.Intn(core.ScreenHeight), color.Gray{Y: uint8(128 + rand.Intn(255-128))})
	}

	img, err := ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	obj := &tempura.Object{
		Tag:      "background",
		Size:     tempura.V(core.ScreenWidth, core.ScreenHeight),
		Pos:      tempura.V(0, 0),
		Drawable: tempura.NewImageDrawable(img),
	}
	return obj, nil
}

func (s *gameScene) newShip(img *ebiten.Image, tag string, initialMoveDirection float64, pos tempura.Vec, rot float64, bulletDirection int, move func() bool) *tempura.Object {
	return &tempura.Object{
		Tag:       tag,
		Size:      tempura.V(shipW, shipH),
		Pos:       pos,
		Drawable:  tempura.NewImageDrawable(img),
		Velocity:  tempura.V(0, shipSpeed*initialMoveDirection),
		RotNormal: rot,
		PreSteps: tempura.MakeBehaviors(
			verticalReflect,
		),
		Steps: tempura.MakeBehaviors(
			tempura.Movement,
		),
	}
}

// verticalReflect reflects an object off top/bottom of screen if
// a collision is imminent
func verticalReflect(source *tempura.Object, dt float64) {
	dy := source.Pos.Y + source.Velocity.Y*dt
	switch {
	case dy < 0:
		fallthrough
	case dy+source.Size.Y > core.ScreenHeight:
		source.Velocity.Y = -source.Velocity.Y
	}
}
