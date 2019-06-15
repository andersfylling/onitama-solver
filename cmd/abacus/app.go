package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "cores",
			Value: 1,
			Usage: "number of cores to use",
		},
		cli.IntFlag{
			Name:  "workers",
			Value: 1,
			Usage: "number of workers. Usually cores-1",
		},
		cli.IntFlag{
			Name:  "depth",
			Value: 1,
			Usage: "game tree height",
		},
	}

	app.Name = "onitamago abacus"
	app.Usage = "for creating and start onitama explorations in the abacus 2.0 supercomputer"
	app.Commands = []cli.Command{
		cmdCreateJobs,
		cmdSearch,
		cmdRead,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
