package main

import (
	"github.com/nsf/termbox-go"
	"os"
)

func Command() {
	cmd := <-command_queue
	if cmd == ":quit" {
		termbox.Close()
		os.Exit(0)
	}
}
