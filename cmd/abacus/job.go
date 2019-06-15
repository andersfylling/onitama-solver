package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andersfylling/onitamago"
	"github.com/urfave/cli"
)

var cmdCreateJobs = cli.Command{
	Name:  "create-jobs",
	Usage: "create sbatch jobs for ABACUS 2.0",
	Action: func(c *cli.Context) error {
		depth := c.GlobalInt("depth")
		cores := c.GlobalInt("cores")
		workers := c.GlobalInt("workers")
		createJobs(cores, workers, depth)
		return nil
	},
}

func createJobs(cores, workers, depth int) {
	allCards := onitamago.Deck(onitamago.DeckOriginal)
	configs := onitamago.GenCardConfigs(allCards)

	prev := -1
	for i, cards := range configs {
		createJob(cores, workers, depth, cards)
		progress := (i / len(configs)) * 100
		if prev != progress {
			prev = progress
			fmt.Println(progress, "% (", i, "/", len(configs), ")")
		}
	}
}

func createJob(cores, workers, depth int, cards []onitamago.Card) {
	//coresStr := strconv.FormatInt(int64(cores), 10)
	workersStr := strconv.FormatInt(int64(workers), 10)
	depthStr := strconv.FormatInt(int64(depth), 10)


	// #SBATCH --ntasks-per-node ` + workersStr + `      # number of workers
	template := `#! /bin/bash
#
#SBATCH --account anders          # account
#SBATCH --nodes 1                 # number of nodes
#SBATCH --time 5:00:00            # max time (HH:MM:SS)

#onitamago:cards=[` + join(cards, ", ", true) + `]
./oniabacus -workers=` + workersStr + ` -depth=` + depthStr + ` search -cards="` + join(cards, ",", false) + `" > ` + join(cards, ".", true) + `.log
`

	f, err := os.Create("./onitamago.job." + join(cards, ".", true) + ".sh")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(template))
	if err != nil {
		panic(err)
	}

	f.Sync()
}
