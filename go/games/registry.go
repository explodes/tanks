package games

import (
	"github.com/explodes/tanks/go/core"
	"github.com/pkg/errors"
)

type GameFactory func(context core.Context) (game core.Game, err error)

var games = make(map[string]GameFactory)

func RegisterGameFactory(name string, factory GameFactory) {
	if _, exists := games[name]; exists {
		panic(errors.Errorf("game %s already registered", name))
	}
	games[name] = factory
}

func GetGameFactory(name string) GameFactory {
	return games[name]
}
