package onitamago

import (
	"math/bits"
	"strconv"
)

const (
	MaskKeyBoards uint64 = 0x1ffffff
)

var (
	zeros = []byte("0000000000000000000000000000000000000000000000000000000000000000")
)

// CacheKey represents all the pieces, current player, relative cards.
// move history is discarded.
//
// Note! this will not function for states with missing masters and as such should not
//  be calculated on leaf nodes.
type CacheKey uint64

func (c *CacheKey) String() string {
	binary := []byte(strconv.FormatUint(uint64(*c), 2))
	binary = append(zeros[:64-len(binary)], binary...)

	merge := func(slices [][]byte, delim byte) (b []byte) {
		for i := range slices {
			b = append(b, slices[i]...)
			b = append(b, delim)
		}

		return b[:len(b)-1]
	}

	// cards
	cards := [][]byte{
		binary[11:14], // suspended
		binary[14:17], // blue 2
		binary[17:20], // blue 1
	}
	segments := [][]byte{
		binary[0:11],                          // unused
		merge(cards, '.'),                     // player cards
		binary[64-25-10-4-4-1 : 64-25-10-4-4], // active player
		binary[64-25-10-4-4 : 64-25-10-4],     // brown master
		binary[64-25-10-4 : 64-25-10],         // blue master
		binary[64-25-10 : 64-25],              // blue students, relative positions
		binary[64-25:],                        // pieces
	}

	return string(merge(segments, '|'))
}

func (c *CacheKey) Encode(st *State) {
	// only call this method after ApplyMove or the current depth of a state
	// has been correctly set. This method sets the cache key for a state at
	// the current depth.
	const bi = BluePlayer * NrOfPieceTypes
	const bsi = (BluePlayer * NrOfPieceTypes) + StudentsIndex
	const bmi = (BluePlayer * NrOfPieceTypes) + MasterIndex
	const bri = BrownPlayer * NrOfPieceTypes
	const brsi = (BrownPlayer * NrOfPieceTypes) + StudentsIndex
	const brmi = (BrownPlayer * NrOfPieceTypes) + MasterIndex

	findCardID := func(card Card) uint64 {
		for i := range st.cards {
			if st.cards[i] == card {
				return uint64(i)
			}
		}

		panic("missing card")
	}

	allBoards := st.board[bsi] | st.board[bmi] | st.board[brsi] | st.board[brmi]
	compact := MakeCompactBoard(allBoards)

	var bluePieces Board // 10 bits, each bit represents the sequence of blue pos in compact
	blueBoards := st.board[bsi] | st.board[bmi]
	//brownBoards := allBoards ^ blueBoards

	highestBlue := uint64(1) << uint64(63-bits.LeadingZeros64(blueBoards))
	blueMask := highestBlue | (highestBlue - 1)
	piecesOfInterest := allBoards & blueMask
	var pos uint64
	for i := LSB(piecesOfInterest); i != 64; i = NLSB(&piecesOfInterest, i) {
		if pieceAtBoardIndex(blueBoards, i) {
			bluePieces |= 1 << pos
		}
		pos++
	}

	findOffset := func(master Board, students Board) uint64 {
		studentsBeforeMaster := (master - 1) & students
		// github.com/tmthrgd/go-popcount.Count64(..) was slower: 69s -> 78s
		return uint64(bits.OnesCount64(studentsBeforeMaster))
	}

	blueMaster := uint64(1) << findOffset(st.board[bmi], st.board[bsi])
	brownMaster := uint64(1) << findOffset(st.board[brmi], st.board[brsi])

	var cards uint64

	// ordering of the players cards does not matter
	// so use the lowest index first to get more cache hits
	cardIDs := [...]uint64{
		findCardID(st.playerCards[0]),
		findCardID(st.playerCards[1]),
		findCardID(st.suspendedCard),
	}
	if cardIDs[0] > cardIDs[1] {
		cardIDs[0], cardIDs[1] = cardIDs[1], cardIDs[0]
	}
	cards |= cardIDs[0]
	cards |= cardIDs[1] << 3
	cards |= cardIDs[2] << 6

	var offset uint64
	holder := compact << offset
	offset += 25
	holder |= bluePieces << offset
	offset += 10
	holder |= blueMaster << offset
	offset += 4
	holder |= brownMaster << offset
	offset += 4
	if st.activePlayer == BluePlayer {
		holder |= 1 << offset
	}
	offset += 1
	holder |= cards << offset
	offset += uint64(3 * len(st.playerCards))

	*c = CacheKey(holder)
	st.setCacheKey(*c)
}

func (c *CacheKey) Decode(st *State) {
	const bsi = (BluePlayer * NrOfPieceTypes) + StudentsIndex
	const bmi = (BluePlayer * NrOfPieceTypes) + MasterIndex
	const brsi = (BrownPlayer * NrOfPieceTypes) + StudentsIndex
	const brmi = (BrownPlayer * NrOfPieceTypes) + MasterIndex

	k := uint64(*c)
	board := CompactBoardToBitBoard(k)
	bluesPos := (k >> 25) & 0x3ff

	cp := board
	var shift uint64
	for i := LSB(cp); i != 64; i = NLSB(&cp, i) {
		blue := 1 & (bluesPos >> shift)
		st.board[bsi] |= blue << i
		shift++
	}
	st.board[brsi] = board ^ st.board[bsi]

	// blue master
	bm := (k >> 35) & 0xf
	var rounds BoardIndex
	if bm > 0 {
		rounds = LSB(bm) + 1
	}

	cp = st.board[bsi]
	var p BoardIndex
	var i BoardIndex
	for p = LSB(cp); i < rounds && p != 64; p = NLSB(&cp, p) {
		i++
	}
	st.board[bmi] = Board(1 << p)
	st.board[bsi] ^= st.board[bmi]

	// brown master
	bm = (k >> 39) & 0xf
	if bm > 0 {
		rounds = LSB(bm) + 1
	}

	cp = st.board[brsi]
	i = 0
	for p = LSB(cp); i < rounds && p != 64; p = NLSB(&cp, p) {
		i++
	}
	st.board[brmi] = Board(1 << p)
	st.board[brsi] ^= st.board[brmi]

	st.activePlayer = (k >> 43) & 1

	cards := (k >> 44) & 0xffff
	for i := range st.playerCards {
		id := cards >> (Board(i) * 4)
		id &= 0xf
		st.playerCards[i] = st.cards[id]
	}

	for i := range st.cards {
		var used bool
		for j := range st.playerCards {
			if used = st.cards[i] == st.playerCards[j]; used {
				break
			}
		}

		if !used {
			st.suspendedCard = st.cards[i]
		}
	}
}
