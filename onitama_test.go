package onitamago

import (
	"fmt"
	"testing"
)

func TestOnitama(t *testing.T) {
	st := NewGame()
	fmt.Println(st.String())
}