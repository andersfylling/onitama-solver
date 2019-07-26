package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/urfave/cli"
)

var cmdDeploy = cli.Command{
	Name:  "show-undeployed",
	Usage: "shows up to 40 undeployed sbatch jobs",
	Action: func(c *cli.Context) error {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		filenames := make([]string, 0, len(files))
		for i := range files {
			filenames = append(filenames, files[i].Name())
		}

		// print the filenames to the terminal, instead of running a sh task from here
		undeployed := rmDeployedFiles(filenames)
		undeployed = rmNonJobFiles(undeployed)
		var size int
		if len(undeployed) >= 15 {
			size = 15
		} else {
			size = len(undeployed)
		}
		for i := range undeployed[:size] {
			fmt.Println(undeployed[i])
		}

		return nil
	},
}

func rmNonJobFiles(files []string) (jobs []string) {
	jobs = make([]string, 0, len(files))
	for i := range files {
		if strings.Contains(files[i], "onijob.") {
			jobs = append(jobs, files[i])
		}
	}
	return
}

func rmDeployedFiles(files []string) (undeployed []string) {
	// every onijob.*.sh that has a onilog.*.log are deployed jobs
	var nrOfFiles int
	for i := len(files) - 1; i > 0; i-- {
		if files[i] == "" || len(files[i]) < 15 {
			continue
		}

		var match bool
		for j := i - 1; j >= 0; j-- {
			if oniFilesRelate(files[i], files[j]) {
				files[j] = ""
				files[i] = ""
				match = true
				break
			}
		}
		if !match {
			nrOfFiles++
		}

		if nrOfFiles >= 40 {
			break
		}
	}

	undeployed = make([]string, 0, 40)
	for i := range files {
		if files[i] == "" {
			continue
		}

		undeployed = append(undeployed, files[i])
		if len(undeployed) == 40 {
			break
		}
	}
	return undeployed
}

func oniFilesRelate(a, b string) bool {
	const OnijobPrefix = "onijob."
	const OnilogPrefix = "onilog."

	if a == "" || b == "" {
		return false
	}

	if a[:len(OnijobPrefix)] != OnijobPrefix {
		return false
	} else if b[:len(OnilogPrefix)] != OnilogPrefix {
		return false
	}

	cardsA := a[len(OnijobPrefix) : len(a)-len(".sh")]
	cardsB := a[len(OnijobPrefix) : len(a)-len(".log")]
	if len(cardsA) != len(cardsB) {
		return false
	} else if cardsA != cardsB {
		return false
	}

	return true
}
