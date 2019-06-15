package main

import (
	"strconv"

	"github.com/andersfylling/onitamago"
)

// [a,b,c,d,e] => "a.b.c.d.e"
func join(cards []onitamago.Card, delimeter string, name bool) (filename string) {
	segment := func(c onitamago.Card) string {
		if name {
			return c.Name()
		} else {
			return strconv.FormatUint(uint64(c), 10)
		}
	}
	for i := range cards {
		filename += segment(cards[i]) + delimeter
	}

	return filename[:len(filename)-len(delimeter)]
}
