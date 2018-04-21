package tanks

import "github.com/explodes/tanks/go/games"

func init() {
	games.RegisterGameFactory("tanks", NewGame)
}
