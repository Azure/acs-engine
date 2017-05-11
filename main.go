package main

import (
	"github.com/Azure/acs-engine/cmd"
	log "github.com/Sirupsen/logrus"
)

const (
	ClientID = "76e0feec-6b7f-41f0-81a7-b1b944520261"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatalln(err)
	}
}
