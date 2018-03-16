package tanks

type Input interface {
	GlobalInput
	TitleInput
	GameInput
}

type GlobalInput interface {
	ToggleFullscreen() bool
	Exit() bool
}

type TitleInput interface {
	Begin() bool
}

type GameInput interface {
	BlueRotate() bool
	RedRotate() bool
}

var _ Input = (*inputImpl)(nil)

type inputImpl struct {
}

func NewInput() Input {
	return &inputImpl{}
}
