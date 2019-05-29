package onitamago

import (
	"fmt"
	"testing"
)

func TestCacheKey_Decode_Encode(t *testing.T) {
	st := NewState()
	st.CreateGame([]Card{
		Frog, Eel,
		Dragon, Crab,
		Tiger,
	})

	moves := []Move{
		50005, // student{A1 => A2, Crab}
		52003, // student{B5 => C4, Eel}
		50533, // student{A2 => A4, Tiger}
		51419, // student{C4 => C3, Crab}
	}
	for _, move := range moves {
		st.ApplyMove(move)
	}

	st.GenerateMoves()
	if st.generatedMovesLen == 0 {
		t.Fatal("expected >0 moves to be generated")
	}

	var key Key
	key.Encode(&st)
	if key == 0 {
		t.Fatal("key should not be 0")
	}

	fmt.Println(key, key.String())

	// duplicates on 234468657779965979

	fmt.Println(st)
	st.Reset()
	fmt.Println(st)
	//	key.Decode(&st)
	//	fmt.Println(st)
}

func TestCacheKey_Encode(t *testing.T) {
	key := Key(7534804192330777627)
	expects := "011.010.001|0|01000|10000|0000011111|0111101000000000000000010000011011"
	if key.String() != expects {
		t.Error("\n :::" + key.String() + "\n  =>" + expects)
	}
}
