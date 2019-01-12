package onitamago

// Board type is based off on bitboards. But uses a mask to identify the 5x5 board
//  63	62	61	60	59	58	57	56
//  55	54	53	52	51	50	49	48
//  47	46	45	44	43	42	41	40
//  39	38	37	36	35	34	33	32
//  31	30	29	28	27	26	25	24
//  23	22	21	20	19	18	17	16
//  15	14	13	12	11	10	09	08
//  07	06	05	04	03	02	01	00
type Board = uint64

const (
	// Board mask is the actual area of the Board type that is used by the players
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	46	45	44	43	42	_	_
	//  _	38	37	36	35	34	_	_
	//  _	30	29	28	27	26	_	_
	//  _	22	21	20	19	18	_	_
	//  _	14	13	12	11	10	_	_
	//  _	_	_	_	_	_	_	_
	BoardMask Board = 0x7c7c7c7c7c00

	// BoardMaskOffset ...
	BoardMaskOffset Board = 0xa

	// CardOffset is how many bit position the initial card masks are shifted
	// remember that offset is number of bit positions
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	42	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	CardOffset Board = 0x2a
)

// FlipVertical Flip a bitboard vertically
func FlipVertical(i Board) Board {
	k1 := Board(0x00FF00FF00FF00FF)
	k2 := Board(0x0000FFFF0000FFFF)
	i = ((i >> 8) & k1) | ((i & k1) << 8)
	i = ((i >> 16) & k2) | ((i & k2) << 16)
	i = (i >> 32) | (i << 32)
	return i
}

// FlipHorizontal Flip a bitboard horizontally
func FlipHorizontal(i Board) Board {
	k1 := Board(0x5555555555555555)
	k2 := Board(0x3333333333333333)
	k4 := Board(0x0f0f0f0f0f0f0f0f)
	i = ((i >> 1) & k1) + 2*(i&k1)
	i = ((i >> 2) & k2) + 4*(i&k2)
	i = ((i >> 4) & k4) + 16*(i&k4)
	return i
}

func RotateBoard(b Board) Board {
	b = FlipVertical(b)
	b = FlipHorizontal(b)

	// move it to match the BoardMask
	b = b >> 8 // one row down
	b = b << 1 // one column left

	return b
}
