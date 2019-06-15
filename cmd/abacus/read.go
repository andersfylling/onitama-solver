package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	ikea "github.com/ikkerens/ikeapack"
	"github.com/urfave/cli"
)


var cmdRead = cli.Command{
	Name:  "read",
	Usage: "read a results file - only gives basic information",
	Action: func(c *cli.Context) error {
		file := c.String("file")
		if file == "" {
			panic("no filename given")
		}

		readFile(file)
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "results file",
		},
	},
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

	fmt.Println(newBlob)
}