package main

import (
	"github.com/nsf/termbox-go"
)

type Status struct {
	mode    string
	cur_pos CursorsPos
	cmd_pos int
	command string
}

var status Status
var event_queue chan termbox.Event
var command_queue chan string

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Init key event channel
	event_queue = make(chan termbox.Event)
	command_queue = make(chan string)

	InitLayout(&status)
	go KeyEventPoll()
	go Command()

	for {
		ProcessCmd(&status)
		Draw()
	}
}
