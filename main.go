package main

import (
	"fmt"
	"github.com/onedss/onegbs/app"
	"github.com/onedss/onegbs/buildtime"
	"github.com/onedss/onegbs/utils"
	"log"
)

var (
	gitCommitCode string
	buildDateTime string
)

func main() {
	log.SetPrefix("[OneGBS] ")
	log.SetFlags(log.LstdFlags)
	if utils.Debug {
		log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	}
	buildtime.BuildVersion = fmt.Sprintf("%s.%s", buildtime.BuildVersion, gitCommitCode)
	if buildDateTime != "" {
		buildtime.BuildTimeStr = buildDateTime
	} else {
		buildtime.BuildTimeStr = buildtime.BuildTime.Format(utils.DateTimeLayout)
	}
	utils.Info("git commit code :", gitCommitCode)
	utils.Info("build date :", buildtime.BuildTimeStr)

	app.StartApp()
}
