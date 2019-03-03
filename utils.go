package onitamago


func CreateRunes(char rune, length int) []rune {
	runes := make([]rune, 0, length)
	for i := 0; i < length; i++ {
		runes = append(runes, char)
	}

	return runes
}

func Merge(boards []Board) (b Board) {
	for i := range boards {
		b |= boards[i]
	}
	return
}