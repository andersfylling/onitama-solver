package main

import (
	"fmt"
)

func concate(t1 IntTuple, t2 IntTuple, n int) [5]int {
	return [5]int{
		t1[0], t1[1], t2[0], t2[1], n,
	}
}

type IntTuple [2]int

func (i *IntTuple) Contains(n int) bool {
	return i[0] == n || i[1] == n
}

func findNrOfCardConfigurations() {
	S := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	inc := func(t *IntTuple) (max bool) {
		t[1]++
		if t[1] == 17 {
			t[1] = 1
		}

		if t[1] == t[0] {
			t[0]++
			t[1] = t[0] + 1
			if t[1] == 17 {
				t[1] = 1
			}
		} else {
			if t[0] == 16 && t[1] == 15 {
				return true
			}
		}

		return false
	}

	tmp := make([]IntTuple, 240)
	tmp[0][0] = 1
	tmp[0][1] = 2

	var i int
	var done bool
	for {
		i++
		tmp[i][0], tmp[i][1] = tmp[i-1][0], tmp[i-1][1]

		if done = inc(&tmp[i]); done {
			break
		}
	}

	tuples := make([]IntTuple, 0, 133770)
	for i := range tmp {
		var duplicate bool
		for j := i - 2; j >= 0; j-- {
			if duplicate = tmp[i][0] == tmp[j][1] && tmp[i][1] == tmp[j][0]; duplicate {
				break
			}
		}

		if !duplicate {
			tuples = append(tuples, tmp[i])
		}
	}

	configs := make([][5]int, 0, 133770)
	for i := range tuples {
		for j := range tuples {
			if tuples[i].Contains(tuples[j][0]) || tuples[i].Contains(tuples[j][1]) {
				continue
			}

			for n := range S {
				if tuples[i].Contains(n) || tuples[j].Contains(n) {
					continue
				}

				configs = append(configs, concate(tuples[i], tuples[j], n))
			}
		}
	}

	// look for duplicates
	for i := range configs {
		cards := configs[i]
		for x := range cards {
			for y := x + 1; y < len(cards); y++ {
				if cards[x] == cards[y] {
					panic(cards)
				}
			}
		}
	}

	fmt.Println("total card configurations", len(configs))
}

func main() {
	findNrOfCardConfigurations()
}
