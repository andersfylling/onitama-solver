package onitamago

const (
	MaskKeyBoards uint64 = 0x1ffffff
)

// CacheKey represents all the pieces, current player, relative cards.
// move history is discarded.
//
// Note! this will not function for states with missing masters and as such should not
//  be calculated on leaf nodes.
type CacheKey uint64

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

	st.cleanTrashBoards() // the trash boards hold "random" set bits
	allBoards := Merge(st.board[:])
	compact := MakeCompactBoard(allBoards)

	var bluePieces Board // 10 bits, each bit represents the sequence of blue pos in compact
	blueBoards := Merge(st.board[bi : bi+NrOfPieceTypes])
	brownBoards := allBoards ^ blueBoards
	var pos uint64
	for i := LSB(allBoards); i != 64; i = NLSB(&allBoards, i) {
		if pieceAtBoardIndex(blueBoards, i) {
			bluePieces |= 1 << pos
		}
		pos++
	}

	var blueMaster Board // <0, 4>, bb
	pos = 0
	for i := LSB(blueBoards); i != 64; i = NLSB(&blueBoards, i) {
		if pieceAtBoardIndex(st.board[bmi], i) {
			if pos > 0 {
				blueMaster = 1 << (pos - 1)
			}
			break
		}
		pos++
	}

	var brownMaster Board // <0, 4>
	pos = 0
	for i := LSB(brownBoards); i != 64; i = NLSB(&brownBoards, i) {
		if pieceAtBoardIndex(st.board[brmi], i) {
			if pos > 0 {
				brownMaster = 1 << (pos - 1)
			}
			break
		}
		pos++
	}

	var cards uint64

	// ordering of the players cards does not matter
	// so use the lowest index first
	cardIDs := []uint64{
		findCardID(st.playerCards[0]),
		findCardID(st.playerCards[1]),
		findCardID(st.playerCards[2]),
		findCardID(st.playerCards[3]),
	}
	if cardIDs[0] > cardIDs[1] {
		cardIDs[0], cardIDs[1] = cardIDs[1], cardIDs[0]
	}
	if cardIDs[2] > cardIDs[3] {
		cardIDs[2], cardIDs[3] = cardIDs[3], cardIDs[2]
	}
	for i := range cardIDs {
		cards |= cardIDs[i] << 4 * uint64(i)
	}

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
	offset += uint64(4 * len(st.playerCards))

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
