package onitamago

// Board type is based off on bitboards. But uses a sub-mask to identify the 5x5 board (range <10, 46>).
//  63	62	61	60	59	58	57	56
//  55	54	53	52	51	50	49	48
//  47	46	45	44	43	42	41	40
//  39	38	37	36	35	34	33	32
//  31	30	29	28	27	26	25	24
//  23	22	21	20	19	18	17	16
//  15	14	13	12	11	10	09	08
//  07	06	05	04	03	02	01	00
type Board = uint64
type BoardIndex = Board
type Index = Board
type Amount = Board

const (
	// Board mask is the actual area of the Board type that is used by the players
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	45	44	43	42	41	_  R5
	//  _	_	37	36	35	34	33	_  R4
	//  _	_	29	28	27	26	35	_  R3
	//  _	_	21	20	19	18	17	_  R2
	//  _	_	13	12	11	10	09	_  R1
	//  _	_	_	_	_	_	_	_
	BoardMask Board = 0x3e3e3e3e3e00

	// rows and columns
	// rows starts at the bottom, going up, just like chess
	// and columns starts left, going right, just like chess
	R1Mask Board = 0x3e00
	R2Mask Board = R1Mask << 8
	R3Mask Board = R2Mask << 8
	R4Mask Board = R3Mask << 8
	R5Mask Board = R4Mask << 8

	// BoardMaskOffset ...
	BoardMaskOffset BoardIndex = 0x9

	TempleTop    Board = 0x80000000000
	TempleBottom Board = 0x800

	StudentsTop Board = 0x360000000000
	StudentsBottom Board =0x3600

	MasterTop Board = TempleTop
	MasterBottom Board = TempleBottom
)

// FlipVertical Flip a bitboard vertically
func FlipVertical(b Board) Board {
	k1 := Board(0x00FF00FF00FF00FF)
	k2 := Board(0x0000FFFF0000FFFF)
	b = ((b >> 8) & k1) | ((b & k1) << 8)
	b = ((b >> 16) & k2) | ((b & k2) << 16)
	b = (b >> 32) | (b << 32)
	return b
}

// FlipHorizontal Flip a bitboard horizontally
func FlipHorizontal(b Board) Board {
	k1 := Board(0x5555555555555555)
	k2 := Board(0x3333333333333333)
	k4 := Board(0x0f0f0f0f0f0f0f0f)
	b = ((b >> 1) & k1) + 2*(b&k1)
	b = ((b >> 2) & k2) + 4*(b&k2)
	b = ((b >> 4) & k4) + 16*(b&k4)
	return b
}

func RotateBoard(b Board) Board {
	b = FlipVertical(b)
	b = FlipHorizontal(b)

	// move it to match the BoardMask
	b = b >> 8 // one row down
	b = b >> 1 // one column right

	return b
}

// 8x8 board index [0, 63] to 5x5 index [0, 24]
func bitboardIndexToOnitamaIndex(i BoardIndex) BoardIndex {
	var diff BoardIndex // i - x = y, where y is the highest pos on the row
	board := boardIndexToBoard(i)
	if (board & R5Mask) > 0 {
		diff = 21
	} else if (board & R4Mask) > 0 {
		diff = 18
	} else if (board & R3Mask) > 0 {
		diff = 15
	} else if (board & R2Mask) > 0 {
		diff = 12
	} else if (board & R1Mask) > 0 {
		diff = 9
	} else {
		panic("no matching row for index")
	}

	return i - diff
}