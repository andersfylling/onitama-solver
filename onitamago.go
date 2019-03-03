package onitamago

func NewGame() (st State, game Game) {
	st = State{}
	st.CreateGame(nil)

	return st, Game{}
}
