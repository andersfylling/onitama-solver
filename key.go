package onitamago

import (
	"fmt"
	"math/bits"
)

const (
	MaskKeyBoards uint64 = 0x1ffffff
)

type StateEncoder interface {
	Encode(*State)
}

type StateDecoder interface {
	Decode(*State)
}

// Key represents all the pieces, current player, relative cards.
// move history is discarded.
//
// Note! this will not function for states with missing masters and as such should not
//  be calculated on leaf nodes.
type Key uint64

// interface implementations
var _ fmt.Stringer = Key(0)      // immutable, reads the key
var _ StateEncoder = (*Key)(nil) // mutable, changes the key
var _ StateDecoder = Key(0)      // immutable, reads the key

// String pretty print the binary version where segments and card indexes are separated.
//  eg. the key 8093039719199539215 is pretty printed as
//      011.100.000|1|01000|00100|0000101111|1011110000000001000000000000001111
func (k Key) String() string {
	binary := []byte(fmt.Sprintf("%064b", k))

	merge := func(slices [][]byte, delim byte) (b []byte) {
		b = make([]byte, 0, 94)
		for i := range slices {
			b = append(b, slices[i]...)
			b = append(b, delim)
		}

		return b[:len(b)-1]
	}

	// cards
	cards := [][]byte{
		binary[0:3], // suspended
		binary[3:6], // blue 2
		binary[6:9], // blue 1
	}
	segments := [][]byte{
		merge(cards, '.'),                     // player cards
		binary[64-34-10-5-5-1 : 64-34-10-5-5], // active player
		binary[64-34-10-5-5 : 64-34-10-5],     // brown master
		binary[64-34-10-5 : 64-34-10],         // blue master
		binary[64-34-10 : 64-34],              // blue students, relative positions
		binary[64-34:],                        // pieces
	}

	return string(merge(segments, '|'))
}

// Encode converts a State into a uint64 key that can be used in caching and other
// situations to compare states. Note that the key is not unique across card configurations
// and assumes that interacting keys uses the same card configuration as its base.
//
// Assumptions/requirements:
//  - The state has both masters
//  - At least one blue piece exists
func (k *Key) Encode(st *State) {
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
	findPieceOffset := func(master Bitboard, students Bitboard) uint64 {
		studentsBeforeMaster := (master - 1) & students
		// github.com/tmthrgd/go-popcount.Count64(..) was slower: 69s -> 78s
		return uint64(bits.OnesCount64(studentsBeforeMaster))
	}

	// hardcode the exact bitboards to use, instead of looping and merging them
	//  - old: benchmarks/cache/uint64-blue-boards-by-slice-merge.txt
	//  - new: benchmarks/cache/uint64-blue-boards-hardcoded.txt
	//
	// This also removed the need to handle the trash bitboard
	//  - old: benchmarks/cache/uint64-clean-trash*
	//  - new: benchmarks/cache/uint64-without-trash*
	blueBoards := st.board[bsi] | st.board[bmi]
	allBoards := blueBoards | st.board[brsi] | st.board[brmi]
	compact := MakeSemiCompactBoard(allBoards)

	// Relative positions for blue pieces.
	// Iterate over only the section of the board that contains blue pieces
	// using "popcount" or "leading zeros".
	//
	//  a=blue, b=brown   X=discarded bits
	//                    after applying "blueMask"
	//  b  _  _  _  b     X  X  X  X  X
	//  _  a  b  _  _     X  a  b  _  _
	//  _  _  b  _  _  => _  _  b  _  _
	//  _  a  _  a  _     _  a  _  a  _
	//  a  _  _  _  a     a  _  _  _  a
	//
	//  b  _  _  _  b     X  X  X  X  X
	//  _  _  b  _  _     X  X  X  X  X
	//  _  _  b  _  _  => X  X  X  X  X
	//  _  a  _  _  b     X  a  _  _  b
	//  a  a  _  a  a     a  a  _  a  a
	//
	//  Note! This only shows a 5x5 board, instead of the 8x8 bitboard.
	//
	// performance improvement:
	//  - old: benchmarks/cache/uint64-without-blue-mask.txt
	//  - new: benchmarks/cache/uint64-blue-mask.txt
	highestBlue := uint64(1) << uint64(63-bits.LeadingZeros64(blueBoards))
	blueMask := highestBlue | (highestBlue - 1)
	piecesOfInterest := allBoards & blueMask
	var pos uint64
	var bluePieces Bitboard // 10 bits, each bit represents the sequence of blue pos in compact
	for i := LSB(piecesOfInterest); i != 64; i = NLSB(&piecesOfInterest, i) {
		if pieceAtBoardIndex(blueBoards, i) {
			bluePieces |= 1 << pos
		}
		pos++
	}

	// blue and brown masters relative position
	blueMaster := uint64(1) << findPieceOffset(st.board[bmi], st.board[bsi])
	brownMaster := uint64(1) << findPieceOffset(st.board[brmi], st.board[brsi])

	// ordering of the players cards does not matter
	// so use the lowest index first to get more cache hits
	// use blue cards + suspended card, instead of all 4 player cards
	//  - old: benchmarks/cache/uint64-cardIDs-loop-suspended*
	//  - new: benchmarks/cache/uint64-cardIDs-hardcoded-suspended*
	cardIDs := [...]uint64{
		findCardID(st.playerCards[0]),
		findCardID(st.playerCards[1]),
		findCardID(st.suspendedCard),
	}
	if cardIDs[0] > cardIDs[1] {
		cardIDs[0], cardIDs[1] = cardIDs[1], cardIDs[0]
	}
	// loop unrolling:
	//  - old: benchmarks/cache/uint64-cardIDs-loop.txt
	//  - new: benchmarks/cache/uint64-cardIDs-hardcoded.txt
	var cards uint64
	cards |= cardIDs[0]
	cards |= cardIDs[1] << 3
	cards |= cardIDs[2] << 6

	var offset uint64
	holder := compact //<< offset
	offset += 34
	holder |= bluePieces << offset
	offset += 10
	holder |= blueMaster << offset
	offset += 5
	holder |= brownMaster << offset
	offset += 5
	holder |= st.activePlayer << offset
	offset += 1
	holder |= cards << offset
	//offset += uint64(3 * len(st.playerCards))

	*k = Key(holder)
	st.setCacheKey(*k)
}

// Deprecated
// Does not work with the new encoder
func (k Key) Decode(st *State) {
	const bsi = (BluePlayer * NrOfPieceTypes) + StudentsIndex
	const bmi = (BluePlayer * NrOfPieceTypes) + MasterIndex
	const brsi = (BrownPlayer * NrOfPieceTypes) + StudentsIndex
	const brmi = (BrownPlayer * NrOfPieceTypes) + MasterIndex

	key := uint64(k)
	board := CompactBoardToBitBoard(key)
	bluesPos := (key >> 25) & 0x3ff

	cp := board
	var shift uint64
	for i := LSB(cp); i != 64; i = NLSB(&cp, i) {
		blue := 1 & (bluesPos >> shift)
		st.board[bsi] |= blue << i
		shift++
	}
	st.board[brsi] = board ^ st.board[bsi]

	// blue master
	bm := (key >> 35) & 0xf
	var rounds BitboardPos
	if bm > 0 {
		rounds = LSB(bm) + 1
	}

	cp = st.board[bsi]
	var p BitboardPos
	var i BitboardPos
	for p = LSB(cp); i < rounds && p != 64; p = NLSB(&cp, p) {
		i++
	}
	st.board[bmi] = Bitboard(1 << p)
	st.board[bsi] ^= st.board[bmi]

	// brown master
	bm = (key >> 39) & 0xf
	if bm > 0 {
		rounds = LSB(bm) + 1
	}

	cp = st.board[brsi]
	i = 0
	for p = LSB(cp); i < rounds && p != 64; p = NLSB(&cp, p) {
		i++
	}
	st.board[brmi] = Bitboard(1 << p)
	st.board[brsi] ^= st.board[brmi]

	st.activePlayer = (key >> 43) & 1

	cards := (key >> 44) & 0xffff
	for i := range st.playerCards {
		id := cards >> (Bitboard(i) * 4)
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
