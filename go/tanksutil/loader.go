package tanksutil

import (
	"bytes"
	"image"
	"io"

	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

type AssetFunc func(name string) ([]byte, error)

type Loader struct {
	bytes AssetFunc
}

func NewLoader(assets AssetFunc) *Loader {
	return &Loader{
		bytes: assets,
	}
}

func (ld *Loader) Reader(name string) (io.Reader, error) {
	b, err := ld.bytes(name)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

type readSeekCloserImpl struct {
	*bytes.Reader
}

func (readSeekCloserImpl) Close() error {
	return nil
}

func (ld *Loader) ReadCloser(name string) (audio.ReadSeekCloser, error) {
	b, err := ld.bytes(name)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b)
	return &readSeekCloserImpl{Reader: reader}, nil
}

func (ld *Loader) Image(name string) (image.Image, error) {
	r, err := ld.Reader(name)
	if err != nil {
		return nil, err
	}
	src, _, err := image.Decode(r)
	return src, err
}

func (ld *Loader) EbitenImage(name string, filter ebiten.Filter) (*ebiten.Image, error) {
	src, err := ld.Image(name)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(src, filter)
}

func (ld *Loader) Font(name string) (*truetype.Font, error) {
	b, err := ld.bytes(name)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(b)
}

func (ld *Loader) Face(name string, size float64) (font.Face, error) {
	b, err := ld.bytes(name)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: size,
	})
	return face, nil
}

type readSeekCloseLengther interface {
	audio.ReadSeekCloser
	Length() int64
}

func (ld *Loader) audioStream(context *audio.Context, fmt, name string) (readSeekCloseLengther, error) {
	rsc, err := ld.ReadCloser(name)
	if err != nil {
		return nil, err
	}
	switch fmt {
	case "mp3":
		return mp3.Decode(context, rsc)
	case "wav":
		return wav.Decode(context, rsc)
	default:
		return nil, errors.Errorf("format not supported: %s", fmt)
	}
}

type AudioPlayer interface {
	Play()
	Pause()
	SetVolume(v float64)
}

type audioReplayer struct {
	stream     []byte
	context    *audio.Context
	lastPlayer *audio.Player
	volume     float64
}

func NewPlayer(context *audio.Context, stream []byte) AudioPlayer {
	return &audioReplayer{
		stream:  stream,
		context: context,
		volume:  1,
	}
}

func (a *audioReplayer) Play() {
	p, err := audio.NewPlayerFromBytes(a.context, a.stream)
	if err != nil {
		return
	}
	p.SetVolume(a.volume)
	p.Play()
	a.lastPlayer = p
}

func (a *audioReplayer) Pause() {
	if a.lastPlayer != nil {
		a.lastPlayer.Pause()
	}
}

func (a *audioReplayer) SetVolume(v float64) {
	if v > 1 {
		v = 1
	} else if v < 0 {
		v = 0
	}
	a.volume = v
}

func (ld *Loader) SFX(context *audio.Context, fmt, name string) (AudioPlayer, error) {
	stream, err := ld.audioStream(context, fmt, name)
	if err != nil {
		return nil, err
	}
	streamBytes, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	return NewPlayer(context, streamBytes), nil
}

func (ld *Loader) AudioLoop(context *audio.Context, fmt, name string) (*audio.Player, error) {
	stream, err := ld.audioStream(context, fmt, name)
	if err != nil {
		return nil, err
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	return audio.NewPlayer(context, loop)
}
