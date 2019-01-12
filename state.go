package onitamago

const (
	BrownPlayer = iota
	BluePlayer
	NrOfPlayers

	OppositePlayer = BrownPlayer
)

type State struct {
	suspendedCard Card

	// Player[0] is at the bottom with the dark pieces. So if the he is not the first player, the board must rotate.
	playerCards  [NrOfPlayers * NrOfPlayerCards]Card
	activePlayer int

	board [NrOfPlayers * NrOfPieceTypes]Board
}

func (st *State) CreateGame(cards []Card) {
	if len(cards) == 0 {
		cards = DrawCards() // random cards
	}
	for i := 0; i < 2; i++ {
		brownIndex := NrOfPlayerCards*BrownPlayer + i
		blueIndex := NrOfPlayerCards*BluePlayer + i
		st.playerCards[brownIndex] = cards[brownIndex]
		st.playerCards[blueIndex] = cards[blueIndex]
	}
	st.suspendedCard = cards[len(cards)-1]

	//st.board =
}

func (st *State) Move(cardIndex, pieceIndex int) {
	// adjust for the player offset
	cardIndex = NrOfPlayerCards*st.activePlayer + cardIndex

	// get the card and swap it with the suspended card
	card := st.playerCards[cardIndex]
	st.swapCard(cardIndex)

	// execute the move
	st.move(card)
	st.nextPlayer()
}

func (st *State) nextPlayer() {
	st.activePlayer = (st.activePlayer + 1) % NrOfPlayers
}

func (st *State) swapCard(cardIndex int) {
	tmp := st.suspendedCard
	st.suspendedCard = st.playerCards[cardIndex]
	st.playerCards[cardIndex] = tmp
}

func (st *State) move(card Card) {
	// TODO: remove if sentence
	if st.activePlayer == OppositePlayer {
		// Rotate the move card and shift it into the original position
		card = RotateBoard(card)
		card = card >> 5     // columns to the right
		card = card << 8 * 5 // rows up
	}

	// TODO: movement...
}
