package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/onitamago"
	ikea "github.com/ikkerens/ikeapack"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type GameMetrics struct {
	ForcedWins [][]onitamago.Move
	Cards      []onitamago.Card
	CardNames  []string
	Metrics    []onitamago.DepthMetric
	Duration   time.Duration
	Depth      uint8
}

func (m *GameMetrics) String() string {
	st := onitamago.State{}
	st.CreateGame(m.Cards)

	state := st.String() + "\n"
	state += "Search depth: " + fmt.Sprint(m.Depth) + "\n"
	state += "Duration: " + fmt.Sprint(m.Duration) + "\n"
	state += "Forced wins (-duplicates): " + fmt.Sprint(len(m.ForcedWins)) + "\n"

	state += fmt.Sprint(m.Metrics)


	return state
}


var cmdSearch = cli.Command{
	Name:  "search",
	Usage: "search runs through the search space and collects data",
	Action: func(c *cli.Context) error {
		depth := c.GlobalInt("depth")
		workers := c.GlobalInt("workers")
		cardsStr := c.String("cards")
		if cardsStr == "" {
			panic("no cards given")
		}
		cardsStrSlice := strings.Split(cardsStr, ",")

		var cards []onitamago.Card
		for i := range cardsStrSlice {
			v, err := strconv.ParseInt(cardsStrSlice[i], 10, 64)
			if err != nil {
				panic(err)
			}

			cards = append(cards, onitamago.Card(uint64(v)))
		}

		fmt.Println(cardsStr)

		startWork(workers, depth, [][]onitamago.Card{cards})
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "cards",
			Value: "",
			Usage: "card configuration",
		},
	},
}

func startWork(workers, depth int, configs [][]onitamago.Card) {
	if len(configs) == 1 {
		results := work(configs[0], uint64(depth))
		saveToFile(results, new(bytes.Buffer))
		return
	}

	storeChan := make(chan *GameMetrics, workers*2)
	go store(storeChan)

	chans := make([]chan []onitamago.Card, workers)
	for i := range chans {
		chans[i] = make(chan []onitamago.Card)
		go worker(uint64(depth), storeChan, chans[i])
	}

	start := time.Now()
	prev := -1
	for i, config := range configs {
		chans[i%workers] <- config

		progress := (float32(i) / float32(len(configs))) * 100.0
		simpleProgress := int(progress * 100)
		if prev != simpleProgress {
			prev = simpleProgress
			fmt.Printf("%d/%d (%.2f%%), %s\n", i, len(configs), progress, time.Now().Sub(start).String())
		}
	}
}

func cardNames(cards []onitamago.Card) (names []string) {
	for i := range cards {
		names = append(names, cards[i].Name())
	}
	return
}

func writeFile(buffer *bytes.Buffer, filename string) error {
	f, err := os.Create("./" + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	f.Sync()
	logrus.Info("saved file " + filename)
	return nil
}

func saveToFile(metric *GameMetrics, b *bytes.Buffer) {
	if b == nil {
		b = new(bytes.Buffer)
	}

	// TODO: proper error handling
	var filename string
	filename += "onitamago."
	filename += "results."
	filename += join(metric.Cards, ".", true)
	filename += ".d" + strconv.FormatUint(uint64(metric.Depth), 10)
	filename += ".ikeapack"

	maxTries := 10
	for i := 1; i <= maxTries; i++ {
		var err error
		if _, err = os.Stat("./" + filename); err == nil {
			filename = strconv.FormatUint(rand.Uint64(), 10) + "." + filename
		}
		if err != nil {
			break
		}
		if i == maxTries {
			panic("max tries hit")
		}
	}

	if err := ikea.Pack(b, metric); err != nil {
		log.Fatalln(err)
		panic("unable to pack file content")
	}

	if err := writeFile(b, filename); err != nil {
		panic("unable to write content to file" + err.Error())
	}
}

func work(cards []onitamago.Card,depth uint64) *GameMetrics {
	m, w, d := onitamago.SearchExhaustiveForForcedWins(cards, depth)
	return &GameMetrics{
		Depth:      uint8(depth),
		Duration:   d,
		ForcedWins: w,
		Metrics:    m,
		Cards:      cards,
		CardNames:  cardNames(cards),
	}
}



func store(metricsChan <-chan *GameMetrics) {
	b := new(bytes.Buffer)
	for {
		metric, ok := <-metricsChan
		if !ok {
			fmt.Println("closing metrics store")
			break
		}

		saveToFile(metric, b)
		b.Reset()
	}
}

func worker(depth uint64, metricsChan chan<- *GameMetrics, c <-chan []onitamago.Card) {
	for {
		cards, ok := <-c
		if !ok {
			fmt.Println("stops worker")
			break
		}
		metrics := work(cards, depth)
		if len(metrics.ForcedWins) > 200 {
			fmt.Println(metrics.CardNames)
		}
		//metricsChan <- work(cards, depth)
	}
}
