package main

import (
	"github.com/nsf/termbox-go"
)

type CursorsPos struct {
	x int
	y int
}

const CMD_LINE = 1
const BAR_POS = 2
const FOOTER_HIGH = 4

func clearc(status *Status) {
	_, h := termbox.Size()
	if status.cmd_pos < 0 {
		return
	}

	termbox.SetCell(status.cmd_pos, h-CMD_LINE, ' ', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(status.cmd_pos, h-CMD_LINE)

	if status.cmd_pos > 0 {
		status.cmd_pos -= 1
		return
	}
}

func putc(status *Status, c rune) {
	w, h := termbox.Size()
	if status.cmd_pos >= w {
		return
	}

	termbox.SetCell(status.cmd_pos, h-CMD_LINE, rune(c), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(status.cmd_pos+1, h-CMD_LINE)
	status.cmd_pos++
}

func LineClear(status *Status) {
	w, h := termbox.Size()
	for i := 0; i < w; i++ {
		termbox.SetCell(i, h-CMD_LINE, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCursor(0, h-CMD_LINE)
	status.cmd_pos = 0

	termbox.Flush()
}

func LinePrint(status *Status, message string) {
	for _, v := range message {
		putc(status, v)
	}
	termbox.Flush()
}

func LinePutc(status *Status, c rune) {
	putc(status, c)
	termbox.Flush()
}

func InitLayout(status *Status) {
	message := "Asterix Emulatore Seriale!"

	w, h := termbox.Size()
	status.cur_pos = CursorsPos{w / 2, h / 2}
	status.cmd_pos = 0

	status.mode = "IDLE"
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	status.cur_pos.x -= len(message) / 2
	for _, v := range message {
		termbox.SetCell(status.cur_pos.x, status.cur_pos.y, rune(v), termbox.ColorDefault, termbox.ColorDefault)
		status.cur_pos.x++
	}

	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	termbox.HideCursor()

	Draw()
}

func Reset(status *Status) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
}

func Draw() {
	w, h := termbox.Size()
	for x := 0; x < w; x++ {
		termbox.SetCell(x, h-BAR_POS, ' ', termbox.ColorDefault, termbox.ColorBlue)
	}

	termbox.Flush()
}

func idle_mode_up(status *Status) int {
	if status.cur_pos.y <= 0 {
		return 0
	}
	status.cur_pos.y -= 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}
func idle_mode_down(status *Status) int {
	_, h := termbox.Size()

	if status.cur_pos.y > h-FOOTER_HIGH {
		return 0
	}
	status.cur_pos.y += 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}

func idle_mode_right(status *Status) int {
	w, _ := termbox.Size()
	if status.cur_pos.x >= w {
		return 0
	}
	status.cur_pos.x += 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}

func idle_mode_left(status *Status) int {
	if status.cur_pos.x <= 0 {
		return 0
	}
	status.cur_pos.x -= 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}

// Command function

func command_mode(status *Status) int {
	status.mode = "CMD"
	LineClear(status)
	LinePutc(status, ':')
	return 0
}

func command_mode_line(status *Status, s rune) int {
	LinePutc(status, s)
	status.command = status.command + string(s)
	return 0
}

func cmd_mode_up(status *Status) int {
	return 0
}

func cmd_mode_down(status *Status) int {
	return 0
}

func cmd_mode_right(status *Status) int {
	w, h := termbox.Size()
	if status.cmd_pos >= w {
		return 0
	}
	status.cmd_pos += 1
	termbox.SetCursor(status.cmd_pos, h-CMD_LINE)
	return 0
}

func cmd_mode_left(status *Status) int {
	_, h := termbox.Size()
	if status.cmd_pos <= 1 {
		return 0
	}
	status.cmd_pos -= 1
	termbox.SetCursor(status.cmd_pos, h-CMD_LINE)
	return 0
}

func cmd_mode_canc(status *Status) int {
	clearc(status)
	return 0
}

func cmd_mode_exec(status *Status) int {
	LineClear(status)
	LinePrint(status, ">>> ")
	LinePrint(status, status.command)
	command_queue <- status.command
	status.command = ""
	return 0
}

type callback func(status *Status) int

var idle_mode_map = map[rune]callback{
	rune(termbox.KeyArrowUp):    idle_mode_up,
	rune(termbox.KeyArrowDown):  idle_mode_down,
	rune(termbox.KeyArrowRight): idle_mode_right,
	rune(termbox.KeyArrowLeft):  idle_mode_left,
	rune('k'):                   idle_mode_up,
	rune('j'):                   idle_mode_down,
	rune('l'):                   idle_mode_right,
	rune('h'):                   idle_mode_left,
	rune(':'):                   command_mode,
}

var cmd_mode_map = map[rune]callback{
	rune(termbox.KeyArrowUp):    cmd_mode_up,
	rune(termbox.KeyArrowDown):  cmd_mode_down,
	rune(termbox.KeyArrowRight): cmd_mode_right,
	rune(termbox.KeyArrowLeft):  cmd_mode_left,
	rune(termbox.KeyBackspace):  cmd_mode_canc,
	rune(termbox.KeyBackspace2): cmd_mode_canc,
	rune(termbox.KeyEnter):      cmd_mode_exec,
}

func KeyEventPoll() {
	for {
		event_queue <- termbox.PollEvent()
	}
}

func ProcessCmd(status *Status) {
	ev := <-event_queue
	if ev.Type != termbox.EventKey {
		return
	}

	if ev.Key == termbox.KeyEsc {
		status.mode = "IDLE"
		Reset(status)
	}

	if status.mode == "IDLE" {
		foo_key, ok_key := idle_mode_map[rune(ev.Key)]
		foo, ok := idle_mode_map[rune(ev.Ch)]

		if ok_key {
			foo_key(status)
		} else if ok {
			foo(status)
		}
	}

	if status.mode == "CMD" {
		foo_key, ok_key := cmd_mode_map[rune(ev.Key)]
		foo, ok := cmd_mode_map[rune(ev.Ch)]

		if ok_key {
			foo_key(status)
		} else if ok {
			foo(status)
		} else {
			command_mode_line(status, ev.Ch)
		}
	}
}
