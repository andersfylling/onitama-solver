package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/andersfylling/onitamago"
	ikea "github.com/ikkerens/ikeapack"
	"github.com/sirupsen/logrus"
)

var cores = 3

const depth = 8

type GameMetrics struct {
	ForcedWins [][]onitamago.Move
	Cards      []onitamago.Card
	CardNames  []string
	Metrics    []onitamago.DepthMetric
	Duration   time.Duration
	Depth      uint8
}

func readFile(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	newBlob := new(GameMetrics)
	r := bytes.NewReader(b)
	if err := ikea.Unpack(r, newBlob); err != nil { // Read *needs* a pointer, or it will panic
		log.Fatalln(err)
	}


	fmt.Println(*newBlob)
}

func main() {
	//readFile("Rooster.Rabbit.Ox.Cobra.Frog.d8.ikeapack")
	//return
	allCards := onitamago.Deck()
	configs := onitamago.GenCardConfigs(allCards)

	storeChan := make(chan *GameMetrics, cores*2)
	go store(storeChan)

	chans := make([]chan []onitamago.Card, cores)
	for i := range chans {
		chans[i] = make(chan []onitamago.Card)
		go worker(storeChan, chans[i])
	}

	start := time.Now()
	prev := -1
	for i, config := range configs {
		chans[i%cores] <- config

		progress := (float32(i) / float32(len(configs))) * 100.0
		simpleProgress := int(progress * 100)
		if prev != simpleProgress {
			prev = simpleProgress
			fmt.Printf("%d/%d (%.2f%%), %s\n", i, len(configs), progress, time.Now().Sub(start).String())
		}
	}
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

func store(metricsChan <-chan *GameMetrics) {
	b := new(bytes.Buffer)
	for {
		metric, ok := <-metricsChan
		if !ok {
			fmt.Println("closing metrics store")
			break
		}

		// TODO: proper error handling
		var filename string
		for i := range metric.CardNames {
			filename += metric.CardNames[i] + "."
		}
		filename += "d" + strconv.FormatUint(uint64(metric.Depth), 10)
		filename += ".ikeapack"

		maxTries := 10
		for i:=1; i <= maxTries; i++{
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
		b.Reset()

	}
}

func worker(metricsChan chan<- *GameMetrics, c <-chan []onitamago.Card) {
	cardNames := func(cards []onitamago.Card) (names []string) {
		for i := range cards {
			names = append(names, cards[i].Name())
		}
		return
	}

	for {
		cards, ok := <-c
		if !ok {
			fmt.Println("stops worker")
			break
		}

		m, w, d := onitamago.SearchForTempleWins(cards, depth)
		metricsChan <- &GameMetrics{
			Depth:      depth,
			Duration:   d,
			ForcedWins: w,
			Metrics:    m,
			Cards:      cards,
			CardNames:  cardNames(cards),
		}
	}
}
