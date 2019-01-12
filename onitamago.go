package onitamago

func NewGame() (st *State) {
	st = &State{}
	st.CreateGame(nil)

	return
}
