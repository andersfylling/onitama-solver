package main

import (
	"fmt"
	. "github.com/andersfylling/onitamago"
	"github.com/disintegration/imaging"
)

func vis(cards []Card, id int) {
	heatmap := BitboardHeatmap{}
	for _, card := range cards {
		heatmap.AddCard(card)
	}

	imaging.Save(heatmap.Render(), fmt.Sprintf("cards-heatmap-%d.png", id))
}

func main() {
	cardsSlice := [...][]Card{
		{4627019807588352, 9191917208207360, 18067449945522176, 4609221463113728, 18137612531269632},
		{18190389089402880, 9191917208207360, 4609221463113728, 9112889809960960, 22553526106324992},
		{9060113251827712, 4609221463113728, 4547854970388480, 38316124802121728, 18067449945522176},
		{22641143439163392, 4547854970388480, 9112889809960960, 2305878331024736256, 18137612531269632},
		{38316124802121728, 9060113251827712, 9130344557051904, 18067449945522176, 4627019807588352},
		/*5*/ {4627019807588352, 18137612531269632, 22553319947894784, 9130344557051904, 4547854970388480},
	}
	for id, cards := range cardsSlice {
		vis(cards, id)
	}

	// most variance in move distance
	vis([]Card{
		Dragon, Tiger, Frog, Rabbit, Crab,
	}, 6)

	// least variance in move distance
	vis([]Card{
		Cobra, Eel, Boar, Ox, Horse,
	}, 7)

	// most number of moves per card
	vis([]Card{
		Dragon, Tiger, Frog, Rabbit, Crab,
	}, 6)

	// lowest number of moves per card
	vis([]Card{
		Dragon, Tiger, Frog, Rabbit, Crab,
	}, 6)
}
