package onitamago

import (
	"strconv"
	"strings"
)

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

func (st *State) MovesLen() int {
	return st.generatedMovesLen
}

func (st *State) Moves() []Move {
	return st.generatedMoves[:st.MovesLen()]
}

func (st *State) Depth() Index {
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
			indexed += "  " + CardName(st.suspendedCard)
		}
		indexed += "\n"
	}
	indexed += "   A B C D E\n"

	// add cards for both players
	var formatted string
	brownIndex = BrownPlayer * NrOfPlayerCards
	for i := 0; i < NrOfPlayerCards; i++ {
		formatted += CardName(st.playerCards[int(brownIndex)+i]) + ", "
	}

	// add the pieces
	formatted += "\n" + indexed + "\n"

	blueIndex = BluePlayer * NrOfPlayerCards
	for i := 0; i < NrOfPlayerCards; i++ {
		formatted += CardName(st.playerCards[int(blueIndex)+i]) + ", "
	}

	return formatted + "\n----------------\n"
}

func (st *State) GenerateMoves() {
	st.generatedMovesLen = 0
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

	st.otherPlayer = BrownPlayer
	st.activePlayer = BluePlayer

	for i := Index(0); i < 2; i++ {
		brownIndex := NrOfPlayerCards*st.otherPlayer + i
		blueIndex := NrOfPlayerCards*st.activePlayer + i
		st.playerCards[brownIndex] = cards[brownIndex]
		st.playerCards[blueIndex] = cards[blueIndex]
	}
	st.suspendedCard = cards[len(cards)-1]

	st.temples[st.otherPlayer] = TempleTop
	st.temples[st.activePlayer] = TempleBottom

	// populate board with pieces
	st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex] = MasterTop
	st.board[st.otherPlayer*NrOfPieceTypes+StudentsIndex] = StudentsTop
	st.board[st.activePlayer*NrOfPieceTypes+MasterIndex] = MasterBottom
	st.board[st.activePlayer*NrOfPieceTypes+StudentsIndex] = StudentsBottom
}

func (st *State) UndoMove() {
	if st.currentDepth == 0 {
		return
	}
	// TODO: undo move
	if !st.hasWon {
		// if it is a winning node, the previously generated moves have not been overwritten
		st.generatedMovesLen = 0
	}
	st.hasWon = false // you can never go beyond a winning node
	st.changePlayer() // we need to make changes to the previous player, not the current

	move := st.previousMoves[st.currentDepth]
	//if st.previousMoves[st.currentDepth] != move {
	//	fmt.Println(st)
	//	fmt.Printf("%+v\n", st.previousMoves)
	//	panic(fmt.Sprintln(st.previousMoves[st.currentDepth], move))
	//}
	st.currentDepth--

	// adjust for the player offset
	cardIndex := NrOfPlayerCards*st.activePlayer + getMoveCardIndex(move)
	st.swapCard(cardIndex)

	// update boards
	from := getMoveFrom(move)
	to := getMoveTo(move)

	friendlyBoardIndex := getMoveFriendlyBoardIndex(move)
	hostileBoardIndex := getMoveHostileBoardIndex(move)

	st.board[st.activePlayer*NrOfPieceTypes+friendlyBoardIndex] ^= boardIndexToBoard(from) | boardIndexToBoard(to)
	st.board[st.otherPlayer*NrOfPieceTypes+hostileBoardIndex] |= boardIndexToBoard(to)
}

func (st *State) ApplyMove(move Move) {
	from := getMoveFrom(move)
	to := getMoveTo(move)

	st.hasWon = getMoveWin(move) == 0x1
	friendlyBoardIndex := getMoveFriendlyBoardIndex(move)
	hostileBoardIndex := getMoveHostileBoardIndex(move)

	// update boards
	st.board[st.activePlayer*NrOfPieceTypes+friendlyBoardIndex] ^= boardIndexToBoard(from) | boardIndexToBoard(to)
	st.board[st.otherPlayer*NrOfPieceTypes+hostileBoardIndex] ^= boardIndexToBoard(to)

	// the move represents the change needed to be done, to reach this depth...
	st.currentDepth++ // TODO: decrement after?
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
