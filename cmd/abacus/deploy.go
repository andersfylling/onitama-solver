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
	Usage: "shows up to 15 undeployed sbatch jobs",
	Action: func(c *cli.Context) error {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		filenames := make([]string, 0, len(files))
		for i := range files {
			filenames = append(filenames, files[i].Name())
		}

		// find logs
		logFiles := make([]string, 0, len(files))
		for i := range filenames {
			if strings.Contains(filenames[i], "onilog") {
				logFiles = append(logFiles, filenames[i])
			}
		}

		// job files
		jobFiles := make([]string, 0, len(files))
		for i := range filenames {
			if strings.Contains(filenames[i], "onijob") {
				jobFiles = append(jobFiles, filenames[i])
			}
		}

		// print the filenames to the terminal, instead of running a sh task from here
		undeployed := rmDeployedFiles(jobFiles, logFiles)
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

func rmDeployedFiles(jobFiles, logFiles []string) (undeployed []string) {
	// every onijob.*.sh that has a onilog.*.log are deployed jobs
	undeployed = make([]string, 0, 15)
	for i := len(jobFiles) - 1; i > 0; i-- {
		if jobFiles[i] == "" {
			continue
		}

		var match bool
		for j := range logFiles {
			if logFiles[j] == "" || !oniFilesRelate(jobFiles[i], logFiles[j]) {
				continue
			}
			logFiles[j] = ""
			jobFiles[i] = ""
			match = true
			break
		}
		if !match {
			undeployed = append(undeployed, jobFiles[i])
		}

		if len(undeployed) >= 15 {
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
	return cardsA == cardsB
}
