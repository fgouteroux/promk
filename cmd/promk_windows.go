//go:build windows
// +build windows

package main

import (
	"log"

	"github.com/fgouteroux/promk/pkg/pusher"
	win "github.com/fgouteroux/promk/pkg/windows"
	"golang.org/x/sys/windows/svc"
)

func main() {
	p := pusher.Setup()

	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatal(err)
	}

	stopCh := make(chan bool)
	if !isInteractive {
		go func() {
			err = svc.Run("Puppet Agent Exporter", win.NewWindowsExporterService(stopCh))
			if err != nil {
				log.Fatalf("Failed to start service: %v", err)
			}
		}()
	}

	go func() {
		p.Run()
	}()

	for {
		if <-stopCh {
			log.Printf("Shutting down %s", "Puppet Agent Exporter")
			break
		}
	}

}
