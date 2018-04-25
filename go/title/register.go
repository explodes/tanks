package title

import "github.com/explodes/tanks/go/games"

func init() {
	games.RegisterGameFactory("title", NewGame)
}
