package tanks

type Input interface {
	TitleInput
	GameInput
}

type TitleInput interface {
	Begin() bool
}

type GameInput interface {
	BlueRotate() bool
	RedRotate() bool
}
