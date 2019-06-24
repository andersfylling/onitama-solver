package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var cmdDeploy = cli.Command{
	Name:  "show-undeployed",
	Usage: "shows up to 40 undeployed sbatch jobs",
	Action: func(c *cli.Context) error {
		var files []string
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			logrus.Error(err)
			return err
		}

		// print the filenames to the terminal, instead of running a sh task from here
		undeployed := rmDeployedFiles(files)
		for i := range undeployed[:40] {
			fmt.Println(undeployed[i])
		}

		return nil
	},
}

func rmDeployedFiles(files []string) (undeployed []string) {
	// every onijob.*.sh that has a onilog.*.log are deployed jobs
	var nrOfFiles int
	for i := len(files) - 1; i > 0; i-- {
		if files[i] == "" {
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
