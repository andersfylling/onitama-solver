package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"

	ikea "github.com/ikkerens/ikeapack"
	"github.com/sirupsen/logrus"
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

	zr, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		logrus.Error("unable to decompress", err)
		zr.Close()
	} else {
		b, err = ioutil.ReadAll(zr)
		zr.Close()
		if err != nil {
			logrus.Error("unable to read file content")
			return
		}
	}


	newBlob := new(GameMetrics)
	r := bytes.NewReader(b)
	if err := ikea.Unpack(r, newBlob); err != nil { // Read *needs* a pointer, or it will panic
		log.Fatalln(err)
	}

	fmt.Println(newBlob)
}