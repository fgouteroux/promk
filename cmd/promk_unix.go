//go:build !windows
// +build !windows

package main

import (
	"github.com/fgouteroux/promk/pkg/pusher"
)

func main() {
	p := pusher.Setup()
	p.Run()
}
