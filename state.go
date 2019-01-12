package onitamago

const (
	BrownPlayer Index = iota
	BluePlayer
	NrOfPlayers Amount = BluePlayer + 1

	OppositePlayer = BrownPlayer
)

type Game struct {
	// game tree with depth of 40
	// one depth at the time... same as chess
	Tree [HighestNrOfMoves * NrOfPlayerPieces * 40]Move
}

func NewState() State {
	return State{
		otherPlayer: 1,
		activePlayer: 0,
	}
}

type State struct {
	suspendedCard Card

	// Player[0] is at the bottom with the dark pieces. So if the he is not the first player, the board must rotate.
	playerCards  [NrOfPlayers * NrOfPlayerCards]Card
	activePlayer Index
	otherPlayer Index

	board [NrOfPlayers * NrOfPieceTypes]Board
	temples [NrOfPlayers]Board

	generatedMoves [HighestNrOfMoves * NrOfPlayerPieces * NrOfPlayerCards]Move
	generatedMovesLen int

	hasWon bool
	previousMove Move
}

func (st *State) GenerateMoves() {
	if st.hasWon {
		return
	}

	generateMoves(st)
}

func (st *State) CreateGame(cards []Card) {
	if len(cards) == 0 {
		cards = DrawCards() // random cards
	}
	for i := Index(0); i < 2; i++ {
		brownIndex := NrOfPlayerCards*BrownPlayer + i
		blueIndex := NrOfPlayerCards*BluePlayer + i
		st.playerCards[brownIndex] = cards[brownIndex]
		st.playerCards[blueIndex] = cards[blueIndex]
	}
	st.suspendedCard = cards[len(cards)-1]

	st.otherPlayer = BluePlayer
	st.activePlayer = BrownPlayer

	st.temples[BrownPlayer] = TempleBottom
	st.temples[BluePlayer] = TempleTop
}

func (st *State) UndoMove() {
	// TODO: undo move
}

func (st *State) ApplyMove(move Move) {

	// TODO: apply the move

	// adjust for the player offset
	cardIndex := NrOfPlayerCards*st.activePlayer + getMoveCardIndex(move)
	st.swapCard(cardIndex)
	st.nextPlayer()
}

func (st *State) nextPlayer() {
	st.otherPlayer = st.activePlayer
	st.activePlayer = (st.activePlayer + 1) % NrOfPlayers
}

func (st *State) swapCard(cardIndex Index) {
	tmp := st.suspendedCard
	st.suspendedCard = st.playerCards[cardIndex]
	st.playerCards[cardIndex] = tmp
}