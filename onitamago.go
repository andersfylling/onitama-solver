package onitamago

import (
	"github.com/sirupsen/logrus"
)

type Number = uint64

func init() {
	// TODO: cleanup
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func NewGame() (st State, game Game) {
	st = State{}
	st.CreateGame(nil)

	return st, Game{}
}
