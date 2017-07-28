package main

import (
	"github.com/Azure/acs-engine/cmd"
	log "github.com/Sirupsen/logrus"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatalln(err)
	}
}
