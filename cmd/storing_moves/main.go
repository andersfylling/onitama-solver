package storing_moves

import (
	"bytes"
	"github.com/andersfylling/onitamago"
	"github.com/ikkerens/ikeapack"
	"io/ioutil"
	"strconv"
)

//   50,000,000
const max = 50000000

type Moves struct {
	cards [5]onitamago.Card

	moves [max][onitamago.MaxDepth]onitamago.Move
	index uint

	buffer    *bytes.Buffer `ikea:"-"`
	path      string        `ikea:"-"`
	filename  string        `ikea:"-"`
	iteration uint64        `ikea:"-"`
}

func (m *Moves) Init() {
	m.buffer = new(bytes.Buffer)
}

func (m *Moves) compress() {
	if err := ikea.Pack(m.buffer, m); err != nil {
		// TODO: backup compression. data must be written to the buffer(!)
		panic(err)
	}
}

func (m *Moves) createFilePath() (path string) {
	if m.filename == "" {
		for _, card := range m.cards {
			m.filename += card.Name() + "_"
		}
		// remove last underscore
		m.filename = m.filename[:len(m.filename)-1]
	}
	if m.path == "" {
		m.path = "."
	}
	path = m.path + "/" + m.filename + "." + strconv.FormatUint(m.iteration, 10) + ".ikeapack"
	m.iteration++

	return path
}
func (m *Moves) save() {
	if err := ioutil.WriteFile(m.createFilePath(), m.buffer.Bytes(), 0644); err != nil {
		// TODO: backup save(!)
		panic(err)
	}
}

func (m *Moves) reset() {
	m.index = 0
	m.buffer.Reset()
}

func (m *Moves) Add(moves [onitamago.MaxDepth]onitamago.Move) {
	if m.index == max {
		m.compress()
		m.save()
		m.reset()
	}
	m.moves[m.index] = moves
	m.index++
}

func (m *Moves) Depth(d int) (moves [][onitamago.MaxDepth]onitamago.Move) {
	// find allocation size
	var size uint
	for i := range m.moves {
		if m.moves[i][d] != 0 {
			continue
		}

		// find depth
		var j int
		for j = range m.moves[i] {
			if m.moves[i][j] == 0 {
				break
			}
		}

		if j == d {
			size++
		}
	}

	moves = make([][onitamago.MaxDepth]onitamago.Move, 0, size)
	for i := range m.moves {
		if m.moves[i][d] != 0 {
			continue
		}

		// find depth
		var j int
		for j = range m.moves[i] {
			if m.moves[i][j] == 0 {
				break
			}
		}

		if j == d {
			moves = append(moves, m.moves[i])
		}
	}

	return moves
}
