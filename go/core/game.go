package core

type Game interface {
	Scene
	OnMuted(muted bool)
	Close() error
}
