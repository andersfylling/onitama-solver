package onitamago

const (
	BrownPlayer Index = iota
	BluePlayer
	NrOfPlayers Amount = BluePlayer + 1

	OppositePlayer = BrownPlayer

	MaxDepth = 60
)

type Game struct {
	// game tree with depth of 40
	// one depth at the time... same as chess
	Tree [HighestNrOfMoves * NrOfPlayerPieces * MaxDepth]Move
}

func NewState() State {
	return State{
		otherPlayer:  1,
		activePlayer: 0,
	}
}

type State struct {
	suspendedCard Card

	// Player[0] is at the bottom with the dark pieces. So if the he is not the first player, the board must rotate.
	playerCards  [NrOfPlayers * NrOfPlayerCards]Card
	activePlayer Index
	otherPlayer  Index

	board   [NrOfPlayers * NrOfPieceTypes]Board
	temples [NrOfPlayers]Board

	generatedMoves    [HighestNrOfMoves * NrOfPlayerPieces * NrOfPlayerCards]Move
	generatedMovesLen int

	hasWon bool

	// the first move is 0, as there is no actual move. Think that every
	// index represents the actual depth.
	previousMoves [MaxDepth]Move
	currentDepth  Index
}

func (st *State) GenerateMoves() {
	if st.hasWon {
		return
	}

	// WARNING: remember to add the generated moves to your game tree as these will be overwritten at the next depth.
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
	if !st.hasWon {
		// if it is a winning node, the previously generated moves have not been overwritten
		st.generatedMovesLen = 0
	}
	st.hasWon = false // you can never go beyond a winning node
	st.changePlayer() // we need to make changes to the previous player, not the current

	move := st.previousMoves[st.currentDepth]

	// adjust for the player offset
	cardIndex := NrOfPlayerCards*st.activePlayer + getMoveCardIndex(move)
	st.swapCard(cardIndex)
	st.currentDepth--
}

func (st *State) ApplyMove(move Move) {
	from := getMoveFrom(move)
	to := getMoveTo(move)

	st.hasWon = getMoveWin(move) == 0x1
	friendlyBoardIndex := getMoveFriendlyBoardIndex(move)
	hostileBoardIndex := getMoveHostileBoardIndex(move)

	st.board[st.activePlayer*NrOfPieceTypes+friendlyBoardIndex] ^= from | to
	st.board[st.otherPlayer*NrOfPieceTypes+hostileBoardIndex] ^= to

	moveBoard := getMoveFrom(move) | getMoveTo(move)
	var offset Amount // piece type. TODO: remove if sentence
	if (st.board[NrOfPlayers*st.activePlayer+1] & moveBoard) > 0 {
		offset++
	}
	st.board[NrOfPlayers*st.activePlayer+offset] ^= moveBoard

	// update opponent
	st.board[NrOfPlayers*st.otherPlayer+0] = (st.board[NrOfPlayers*st.otherPlayer+0] | moveBoard) ^ moveBoard
	st.board[NrOfPlayers*st.otherPlayer+1] = (st.board[NrOfPlayers*st.otherPlayer+1] | moveBoard) ^ moveBoard

	// TODO: remove comparison
	st.hasWon = getMoveWin(move) == 1

	st.generatedMovesLen = 0

	// the move represents the change needed to be done, to reach this depth...
	st.currentDepth++
	st.previousMoves[st.currentDepth] = move

	// adjust for the player offset
	cardIndex := NrOfPlayerCards*st.activePlayer + getMoveCardIndex(move)
	st.swapCard(cardIndex)
	st.changePlayer()
}

func (st *State) changePlayer() {
	st.otherPlayer = st.activePlayer
	st.activePlayer = (st.activePlayer + 1) % NrOfPlayers
}

func (st *State) swapCard(cardIndex Index) {
	tmp := st.suspendedCard
	st.suspendedCard = st.playerCards[cardIndex]
	st.playerCards[cardIndex] = tmp
}
