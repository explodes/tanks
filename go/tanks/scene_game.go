package tanks

import (
	"github.com/hajimehoshi/ebiten"

	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"math/rand"

	"fmt"

	"strconv"

	"github.com/explodes/tanks/go/tanksutil"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

var _ Scene = (*gameScene)(nil)

type Phase uint8

const (
	phaseCountdown Phase = iota
	phaseBattle
	phaseBlueVictory
	phaseRedVictory
)

const (
	tankRotatesPerSecond  = 0.5
	tankSpeed             = 215
	tankWidth, tankHeight = 170 * 3 / 10, 200 * 3 / 10

	tankCollisionScale = 0.75

	victoryMessageDuration = 3

	autoShotPerSecond = 0.5

	bulletSpeed = 580

	tagBackground = "background"
	tagBluePlayer = "bluePlayer"
	tagBlueBullet = "blueBullet"
	tagRedPlayer  = "redPlayer"
	tagRedBullet  = "redBullet"
)

const (
	layerBackground = iota
	layerTanks
	layerBullets
	numLayers
)

var (
	tankRotateOffset = tanksutil.DegToRad(-90)

	winningMessages = []string{
		"%s has become the champion",
		"%s is victorious",
		"%s was better",
	}

	countdownColors = []color.Color{
		colornames.Red,
		colornames.Blue,
		colornames.White,
	}
)

type gameScene struct {
	g    *Game
	time float64

	phase Phase

	messageFace font.Face
	message     *tanksutil.Text

	cannonSFX tanksutil.AudioPlayer

	bluePlayer *tanksutil.Object
	redPlayer  *tanksutil.Object

	victoryTime float64

	shot tanksutil.Drawable

	blueShotDelay float64
	redShotDelay  float64

	layers tanksutil.Layers
}

func NewGameScene(game *Game) (Scene, error) {
	cannonSFX, err := game.loader.SFX(game.audioContext, "wav", "sound/tank.wav")
	if err != nil {
		return nil, err
	}

	messageFace, err := game.loader.Face("fonts/DampfPlatzs.ttf", 42)
	if err != nil {
		return nil, err
	}

	shotImage, err := game.loader.EbitenImage("images/shot.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	shotDrawable := tanksutil.NewImageDrawable(shotImage)

	dirtImage, err := game.loader.EbitenImage("images/dirt.jpg", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	tanksImage, err := game.loader.EbitenImage("images/tanks.png", ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	blueTankDrawable := tanksutil.NewImageDrawableFrames(tanksImage, tanksutil.R(0, 0, 148, 333./2))
	redTankDrawable := tanksutil.NewImageDrawableFrames(tanksImage, tanksutil.R(0, 333./2, 148, 333))

	s := &gameScene{
		g:           game,
		phase:       phaseCountdown,
		cannonSFX:   cannonSFX,
		messageFace: messageFace,
		shot:        shotDrawable,
		layers:      tanksutil.NewLayers(numLayers),
	}

	rotBlue := tanksutil.DegToRad(135)
	rotRed := tanksutil.DegToRad(-45)
	if rand.Float64() < 0.5 {
		rotBlue, rotRed = rotRed, rotBlue
	}

	bluePlayer := &tanksutil.Object{
		Tag:       tagBluePlayer,
		Pos:       tanksutil.V(100, ScreenHeight/2-tankHeight/2),
		Size:      tanksutil.V(tankWidth, tankHeight),
		Drawable:  blueTankDrawable,
		Rot:       rotBlue,
		RotNormal: tankRotateOffset,

		Steps: tanksutil.MakeBehaviors(
			s.behaviorBlueRotateOnButton,
		),
		PostSteps: tanksutil.MakeBehaviors(
			s.reflectInBounds,
			s.behaviorBlueHitsRedBullet,
		),
	}
	s.bluePlayer = bluePlayer
	s.layers[layerTanks].Add(bluePlayer)

	redPlayer := &tanksutil.Object{
		Tag:       tagRedPlayer,
		Pos:       tanksutil.V(ScreenWidth-100-tankWidth, ScreenHeight/2-tankHeight/2),
		Size:      tanksutil.V(tankWidth, tankHeight),
		Drawable:  redTankDrawable,
		Rot:       rotRed,
		RotNormal: tankRotateOffset,

		Steps: tanksutil.MakeBehaviors(
			s.behaviorRedRotateOnButton,
		),
		PostSteps: tanksutil.MakeBehaviors(
			s.reflectInBounds,
			s.behaviorRedHitsBlueBullet,
		),
	}
	s.redPlayer = redPlayer
	s.layers[layerTanks].Add(redPlayer)

	dirt := &tanksutil.Object{
		Tag:      tagBackground,
		Size:     tanksutil.V(ScreenWidth, ScreenHeight),
		Drawable: tanksutil.NewImageDrawable(dirtImage),
	}
	s.layers[layerBackground].Add(dirt)

	return s, nil
}

func (s *gameScene) Update(dt float64) error {
	s.time += dt

	switch s.phase {
	case phaseCountdown:

		countdownTime := s.time * 2
		if countdownTime >= 3 {
			s.phase = phaseBattle
			break
		}
		seconds := 3 - int(countdownTime)

		countdownColorIndex := 3 - seconds
		if countdownColorIndex < 0 {
			countdownColorIndex = 0
		}
		text := tanksutil.NewText(s.messageFace, countdownColors[countdownColorIndex], strconv.Itoa(seconds))
		s.message = &text
	case phaseBattle:
		s.blueShotDelay += dt
		s.redShotDelay += dt
		s.layers.Update(dt)
	case phaseBlueVictory:
		fallthrough
	case phaseRedVictory:
		s.victoryTime -= dt
		if s.victoryTime <= 0 {
			return s.g.SetNewScene(NewTitleScene)
		}
	}

	return nil
}

func (s *gameScene) Draw(image *ebiten.Image) {
	s.layers.Draw(image)

	switch s.phase {
	case phaseBattle:
	case phaseBlueVictory:
		fallthrough
	case phaseRedVictory:
		fallthrough
	case phaseCountdown:
		if s.message == nil {
			return
		}
		s.message.Draw(image, ScreenWidth/2, ScreenHeight/2+s.message.H/2, tanksutil.AlignCenter)
	}
}

func (s *gameScene) reflectInBounds(source *tanksutil.Object, dt float64) {
	objBounds := source.Bounds()
	switch {
	case objBounds.Min.X <= 0:
		source.Velocity = tanksutil.V(-source.Velocity.X, source.Velocity.Y)
		source.Rot = source.Velocity.Angle()
		source.Pos = tanksutil.V(0, source.Pos.Y)
	case objBounds.Max.X >= ScreenWidth:
		source.Velocity = tanksutil.V(-source.Velocity.X, source.Velocity.Y)
		source.Rot = source.Velocity.Angle()
		source.Pos = tanksutil.V(ScreenWidth-source.Size.X, source.Pos.Y)
	}
	switch {
	case objBounds.Min.Y <= 0:
		source.Velocity = tanksutil.V(source.Velocity.X, -source.Velocity.Y)
		source.Rot = source.Velocity.Angle()
		source.Pos = tanksutil.V(source.Pos.X, 0)
	case objBounds.Max.Y >= ScreenHeight:
		source.Velocity = tanksutil.V(source.Velocity.X, -source.Velocity.Y)
		source.Rot = source.Velocity.Angle()
		source.Pos = tanksutil.V(source.Pos.X, ScreenHeight-source.Size.Y)
	}
}

func (s *gameScene) behaviorBlueRotateOnButton(source *tanksutil.Object, dt float64) {
	if s.g.input.BlueRotate() {
		// rotate
		source.Rot += tanksutil.DegToRad(-tankRotatesPerSecond*360) * dt
		s.blueShotDelay = 0
	} else {
		source.Velocity = tanksutil.V(tankSpeed, 0).Rotated(source.Rot)
		tanksutil.Movement(source, dt)
		if s.blueShotDelay > 1.0/autoShotPerSecond {
			s.spawnBlueShots()
			s.blueShotDelay = 0
		}
	}
}

func (s *gameScene) behaviorRedRotateOnButton(source *tanksutil.Object, dt float64) {
	if s.g.input.RedRotate() {
		// rotate
		source.Rot += tanksutil.DegToRad(-tankRotatesPerSecond*360) * dt
		s.redShotDelay = 0
	} else {
		source.Velocity = tanksutil.V(tankSpeed, 0).Rotated(source.Rot)
		tanksutil.Movement(source, dt)
		if s.redShotDelay > 1.0/autoShotPerSecond {
			s.spawnRedShots()
			s.redShotDelay = 0
		}
	}
}

func (s *gameScene) spawnBlueShots() {

	bounds := s.bluePlayer.Bounds()
	pos1 := bounds.Center().Add(tanksutil.V(bounds.W()/2, 2).Rotated(s.bluePlayer.Rot))
	pos2 := bounds.Center().Add(tanksutil.V(bounds.W()/2, -8).Rotated(s.bluePlayer.Rot))

	blueBullet1 := &tanksutil.Object{
		Tag:      tagBlueBullet,
		Pos:      pos1,
		Size:     tanksutil.V(8, 8),
		Drawable: s.shot,
		Velocity: tanksutil.V(bulletSpeed, 0).Rotated(s.bluePlayer.Rot),
		Steps: tanksutil.MakeBehaviors(
			tanksutil.Movement,
		),
		PostSteps: tanksutil.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	blueBullet2 := &tanksutil.Object{
		Tag:      tagBlueBullet,
		Pos:      pos2,
		Size:     tanksutil.V(8, 8),
		Drawable: s.shot,
		Velocity: tanksutil.V(bulletSpeed, 0).Rotated(s.bluePlayer.Rot),
		Steps: tanksutil.MakeBehaviors(
			tanksutil.Movement,
		),
		PostSteps: tanksutil.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	s.layers[layerBullets].Add(blueBullet1)
	s.layers[layerBullets].Add(blueBullet2)

	s.cannonSFX.Play()
}

func (s *gameScene) spawnRedShots() {

	bounds := s.redPlayer.Bounds()
	offset := tanksutil.V(bounds.H()/2, -8).Rotated(s.redPlayer.Rot)
	pos := bounds.Center().Add(offset)

	redBullet := &tanksutil.Object{
		Tag:      tagRedBullet,
		Pos:      pos,
		Size:     tanksutil.V(14, 14),
		Drawable: s.shot,
		Velocity: tanksutil.V(bulletSpeed, 0).Rotated(s.redPlayer.Rot),
		Steps: tanksutil.MakeBehaviors(
			tanksutil.Movement,
		),
		PostSteps: tanksutil.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	s.layers[layerBullets].Add(redBullet)

	s.cannonSFX.Play()
}

func (s *gameScene) behaviorRemoveOutOfBounds(source *tanksutil.Object, dt float64) {
	if !tanksutil.Collision(source.Bounds(), ScreenBounds) {
		s.layers[layerBullets].Remove(source)
	}
}

func (s *gameScene) behaviorRedHitsBlueBullet(source *tanksutil.Object, dt float64) {
	if s.phase != phaseBattle {
		return
	}
	sourceBounds := source.Bounds().ScaledAtCenter(tankCollisionScale)
	iter := s.layers.TagIterator(tagBlueBullet)
	for bullet, ok := iter(); ok; bullet, ok = iter() {
		if tanksutil.Collision(sourceBounds, bullet.Bounds()) {
			s.g.blueScore++
			s.phase = phaseBlueVictory
			s.onVictory("Blue", colornames.Cadetblue)
		}
	}
}

func (s *gameScene) behaviorBlueHitsRedBullet(source *tanksutil.Object, dt float64) {
	if s.phase != phaseBattle {
		return
	}
	sourceBounds := source.Bounds().ScaledAtCenter(tankCollisionScale)
	iter := s.layers.TagIterator(tagRedBullet)
	for bullet, ok := iter(); ok; bullet, ok = iter() {
		if tanksutil.Collision(sourceBounds, bullet.Bounds()) {
			s.g.redScore++
			s.phase = phaseRedVictory
			s.onVictory("Red", colornames.Indianred)
		}
	}
}

func (s *gameScene) onVictory(winner string, textColor color.Color) {
	s.victoryTime = victoryMessageDuration

	saying := winningMessages[rand.Intn(len(winningMessages))]
	victoryMessage := fmt.Sprintf(saying, winner)

	text := tanksutil.NewText(s.messageFace, textColor, victoryMessage)
	s.message = &text
}
