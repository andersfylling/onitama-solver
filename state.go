package onitamago

import (
	"strconv"
	"strings"
)

const (
	BrownPlayer Number = iota
	BluePlayer
	NrOfPlayers Number = BluePlayer + 1

	OppositePlayer = BrownPlayer
)

type Game struct {
	// game tree with Depth of 40
	// one Depth at the time... same as chess
	Tree [HighestNrOfMoves * NrOfPlayerPieces * MaxDepth]Move
}

func NewState() State {
	return State{
		otherPlayer:  1,
		activePlayer: 0,
		cards:        [5]Card{},
	}
}

type State struct {
	suspendedCard Card

	// Player[0] is at the bottom with the dark pieces. So if the he is not the first player, the board must rotate.
	cards        [NrOfPlayers*NrOfPlayerCards + 1]Card
	playerCards  [NrOfPlayers * NrOfPlayerCards]Card
	activePlayer Number
	otherPlayer  Number

	board   [NrOfPlayers * NrOfPieceTypes]Bitboard
	temples [NrOfPlayers]Bitboard

	generatedMoves    [HighestNrOfMoves * NrOfPlayerPieces * NrOfPlayerCards]Move
	generatedMovesLen int

	hasWon bool

	// the first move is 0, as there is no actual move. Think that every
	// index represents the actual Depth.
	previousMoves [MaxDepth]Move
	currentDepth  Number // NOTE! this must be handled during caching (key decoding)

	// this is only activated when using the build tag "onitama_cache"
	previousCacheKeys previousCacheKeys
}

func (st *State) Reset() {
	st.suspendedCard = 0
	//for i := range st.cards {
	//	st.cards[i] = 0
	//}
	for i := range st.playerCards {
		st.playerCards[i] = 0
	}
	st.activePlayer = 0
	st.otherPlayer = 0
	for i := range st.board {
		st.board[i] = 0
	}
	for i := range st.temples {
		st.temples[i] = 0
	}
	st.generatedMovesLen = 0
	st.hasWon = false
	st.currentDepth = 0
}

func (st *State) MoveHistory() []Move {
	return st.previousMoves[1 : st.Depth()+1]
}

func (st *State) cleanTrashBoards() {
	st.board[(NrOfPieceTypes*BluePlayer)+TrashIndex] = 0
	st.board[(NrOfPieceTypes*BrownPlayer)+TrashIndex] = 0
}

func (st *State) MovesLen() int {
	return st.generatedMovesLen
}

func (st *State) Moves() []Move {
	return st.generatedMoves[:st.MovesLen()]
}

func (st *State) Depth() Number {
	return st.currentDepth
}

//var _ fmt.Stringer = (*State)(nil)

func (st State) String() string {
	// white = who ever moves first
	// this can be determined using State.currentDepth
	const blueStudent = '♟'
	const blueMaster = '♚'
	const blueTemple = '⚑'
	const brownStudent = '♙'
	const brownMaster = '♔'
	const brownTemple = '⚐'
	const empty = '-'

	//piece := func(master bool) rune {
	//	if st.currentDepth % 2 == 0 {
	//		if master {
	//			return blueMaster
	//		}
	//		return blueStudent
	//	}
	//	if master {
	//		return brownMaster
	//	}
	//	return brownStudent
	//}
	//temple := func(white bool) string {
	//	if white {
	//		return blueTemple
	//	}
	//	return brownTemple
	//}

	// create a 5x5 board representation
	board := CreateRunes(empty, 25)

	blueIndex := BluePlayer * NrOfPieceTypes
	brownIndex := BrownPlayer * NrOfPieceTypes
	blueTempleIndex := 24 - 2   // reversed
	brownTempleIndex := 24 - 22 // reversed
	bluePieces := Merge(st.board[blueIndex : blueIndex+NrOfPieceTypes-1])
	brownPieces := Merge(st.board[brownIndex : brownIndex+NrOfPieceTypes-1])

	for i := LSB(bluePieces); i != 64; i = NLSB(&bluePieces, i) {
		index := bitboardIndexToOnitamaIndex(i)
		// reverse it, since string and LSB starts at different ends
		index = 24 - index

		if pieceAtBoardIndex(st.board[blueIndex+MasterIndex], i) {
			board[index] = blueMaster
		} else {
			board[index] = blueStudent
		}
	}
	if board[blueTempleIndex] == empty {
		board[blueTempleIndex] = blueTemple
	}

	for i := LSB(brownPieces); i != 64; i = NLSB(&brownPieces, i) {
		index := bitboardIndexToOnitamaIndex(i)
		// reverse it, since string and LSB starts at different ends
		index = 24 - index

		if pieceAtBoardIndex(st.board[brownIndex+MasterIndex], i) {
			board[index] = brownMaster
		} else {
			board[index] = brownStudent
		}
	}
	if board[brownTempleIndex] == empty {
		board[brownTempleIndex] = brownTemple
	}

	// add new lines, spacing, rows and col identifiers
	// and the idle card
	var indexed string = "\n"
	rows := []int64{5, 4, 3, 2, 1}
	for i := range rows {
		tmp := strconv.FormatInt(rows[i], 10) + string(board[5*i:5*i+5])
		split := strings.Split(tmp, "")
		indexed += " " + strings.Join(split, " ")
		if i == 2 {
			indexed += "  " + st.suspendedCard.Name()
		}
		indexed += "\n"
	}
	indexed += "   A B C D E\n"

	// add cards for both players
	var formatted string
	brownIndex = BrownPlayer * NrOfPlayerCards
	for i := 0; i < NrOfPlayerCards; i++ {
		formatted += st.playerCards[int(brownIndex)+i].Name() + ", "
	}

	// add the pieces
	formatted += "\n" + indexed + "\n"

	blueIndex = BluePlayer * NrOfPlayerCards
	for i := 0; i < NrOfPlayerCards; i++ {
		formatted += st.playerCards[int(blueIndex)+i].Name() + ", "
	}

	return formatted + "\n----------------\n"
}

func (st *State) GenerateMoves() {
	st.generatedMovesLen = 0
	if st.hasWon {
		return
	}

	// WARNING: remember to add the generated moves to your game tree as these will be overwritten at the next Depth.
	generateMoves(st)
}

// CreateGame
//  cards[0-1]: brown player (opponent)
//  cards[2-3]: blue player (current)
//  cards[4-4]: idle card
func (st *State) CreateGame(cards []Card) {
	if len(cards) == 0 {
		cards = DrawCards() // random cards
	}

	for i := range cards {
		st.cards[i] = cards[i]
	}

	// organize cards
	st.otherPlayer = BrownPlayer
	st.activePlayer = BluePlayer
	for i := Number(0); i < 2; i++ {
		brownIndex := NrOfPlayerCards*st.otherPlayer + i
		blueIndex := NrOfPlayerCards*st.activePlayer + i
		st.playerCards[brownIndex] = cards[brownIndex]
		st.playerCards[blueIndex] = cards[blueIndex]
	}
	st.suspendedCard = cards[len(cards)-1]

	// clear boards
	for i := range st.temples {
		st.temples[i] = 0
	}
	for i := range st.board {
		st.board[i] = 0
	}

	// set up pieces and temples
	st.temples[st.otherPlayer] = TempleTop
	st.temples[st.activePlayer] = TempleBottom
	st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex] = MasterTop
	st.board[st.otherPlayer*NrOfPieceTypes+StudentsIndex] = StudentsTop
	st.board[st.activePlayer*NrOfPieceTypes+MasterIndex] = MasterBottom
	st.board[st.activePlayer*NrOfPieceTypes+StudentsIndex] = StudentsBottom
}

func (st *State) UndoMove() {
	if st.currentDepth == 0 {
		return
	}

	// Note: not setting generatedMovesLen to 0 when the state has .hasWon,
	// the leafs will still be in memory, and there is no reason to regenerate these.
	// However, if there's a game tree used; there should be no reason care about the already stored moves.
	st.generatedMovesLen = 0
	st.hasWon = false // you can never go beyond a winning node

	// we need to make changes to the previous player, not the current
	st.changePlayer()
	move := st.previousMoves[st.currentDepth]
	cardIndex := NrOfPlayerCards*st.activePlayer + move.CardIndex() // adjust for the player offset
	st.swapCard(cardIndex)
	st.currentDepth--

	// update boards
	from := move.From()
	to := move.To()

	friendlyBoardIndex := move.BoardIndex()
	hostileBoardIndex := move.HostileBoardIndex()

	st.board[st.activePlayer*NrOfPieceTypes+friendlyBoardIndex] ^= boardIndexToBoard(from) | boardIndexToBoard(to)
	st.board[st.otherPlayer*NrOfPieceTypes+hostileBoardIndex] |= boardIndexToBoard(to)

	st.removeLastCacheKey()
}

func (st *State) ApplyMove(move Move) {
	from := move.From()
	to := move.To()

	st.hasWon = move.Win()
	friendlyBoardIndex := move.BoardIndex()
	hostileBoardIndex := move.HostileBoardIndex()

	// update boards
	st.board[st.activePlayer*NrOfPieceTypes+friendlyBoardIndex] ^= boardIndexToBoard(from) | boardIndexToBoard(to)
	st.board[st.otherPlayer*NrOfPieceTypes+hostileBoardIndex] ^= boardIndexToBoard(to)

	// the move represents the change needed to be done, to reach this Depth...
	st.currentDepth++ // TODO: decrement after?
	st.previousMoves[st.currentDepth] = move

	// adjust for the player offset
	cardIndex := NrOfPlayerCards*st.activePlayer + move.CardIndex()
	st.swapCard(cardIndex)
	st.changePlayer()
}

func (st *State) NextPlayer() int {
	return int((st.activePlayer + 1) % NrOfPlayers)
}

func (st *State) changePlayer() {
	st.otherPlayer = st.activePlayer
	st.activePlayer = (st.activePlayer + 1) % NrOfPlayers
}

func (st *State) swapCard(cardIndex Number) {
	st.suspendedCard, st.playerCards[cardIndex] = st.playerCards[cardIndex], st.suspendedCard
}
