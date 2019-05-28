package bits

// BenchmarkRotateCard
//func BenchmarkRotateCard(b *testing.B) {
//	cards := [...]Card{
//		Rooster, Rabbit, Ox, Cobra,
//		Horse, Goose, Frog, Eel,
//		Tiger, Dragon, Crab, Elephant,
//		Monkey, Mantis, Crane, Boar,
//	}
//	highestShift := 24 // end of board
//
//	var safeRotate = func(card Card) Card {
//		// Rotate the move card and shift it into the original position
//		card = FlipVertical(card)
//		card = FlipHorizontal(card)
//		card = card << 3
//		card = card << (8 * 3)
//
//		return card
//	}
//
//	// tuples of [normal, rotated]
//	var inputs [][2]Card
//	for _, card := range cards {
//		for i := 0; i <= highestShift; i++ {
//			normal := card >> Card(i)
//			rotated := safeRotate(normal)
//			inputs = append(inputs, [2]Card{normal, rotated})
//		}
//	}
//
//	var fastRotate = func(card Card) Card {
//		return safeRotate(card)
//	}
//
//	// verify success
//	for j := range inputs {
//		a, b := inputs[j][0], inputs[j][1]
//		if fastRotate(a) != b {
//			panic("fast rotate does not equal safe rotate")
//		}
//	}
//
//	b.Run("rotate", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			for j := range inputs {
//				_ = fastRotate(inputs[j][0])
//			}
//		}
//	})
//}
