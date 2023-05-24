package main

import (
	"fmt"
	"os"

	"github.com/siddharth178/dd/internal"
	log "github.com/sirupsen/logrus"
)

const (
	CMDNAME = "detectdups"
)

func usageAndExit() {
	fmt.Printf(`usage: %s <dir path>
`, CMDNAME)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usageAndExit()
	}
	log.SetOutput(os.Stderr)

	dirPath := os.Args[1]
	d, err := os.Open(dirPath)
	if err != nil {
		log.Errorln("can't read", dirPath)
		usageAndExit()
	}
	stat, err := d.Stat()
	if err != nil {
		log.Errorln("can't stat", dirPath)
		usageAndExit()
	}
	if !stat.IsDir() {
		log.Errorln(dirPath, "is not a dir")
		usageAndExit()
	}
	d.Close()

	log.Infoln("processing dir:", dirPath)
	log.Infoln("running", CMDNAME)
	internal.LibMain(dirPath)
	log.Infoln(CMDNAME, "exiting")
}
