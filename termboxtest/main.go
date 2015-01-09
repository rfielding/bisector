package main

import "github.com/nsf/termbox-go"

type UIState struct {
	V rune
	X int
	Y int
}

func newUIState() *UIState {
	return &UIState{'x', 20, 10}
}

func redraw(ui *UIState) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(ui.X, ui.Y, ui.V, termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func respond(k termbox.Key, ui *UIState) *UIState {
	switch k {
	case termbox.KeyArrowUp:
		ui.Y = ui.Y - 1
	case termbox.KeyArrowDown:
		ui.Y = ui.Y + 1
	case termbox.KeyArrowLeft:
		ui.X = ui.X - 1
	case termbox.KeyArrowRight:
		ui.X = ui.X + 1
	}
	return ui
}

func handlersPoll(ui *UIState) (bool, *UIState) {
	ev := termbox.PollEvent()
	t, k := ev.Type, ev.Key

	if t == termbox.EventKey && k == termbox.KeyEsc {
		return false, ui
	} else if t == termbox.EventResize {
		redraw(ui)
		return true, ui
	} else if isArrow(t, k) {
		ui = respond(k, ui)
		redraw(ui)
		return true, ui
	} else if isAlgebraChar(ev) {
		ui.V = ev.Ch
		redraw(ui)
		return true, ui
	}

	return true, ui
}

func setup() {
	termbox.SetInputMode(termbox.InputCurrent)
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Sync()
}

func isAlgebraChar(ev termbox.Event) bool {
	return isGrouping(ev) || isBinOp(ev) || isToken(ev)
}

func isGrouping(ev termbox.Event) bool {
	t, c, k := ev.Type, ev.Ch, ev.Key
	return (t == termbox.EventKey) && ((c == 0 && (k == termbox.KeySpace)) || ((c == 0) && (k == termbox.KeyEnter)) || ('(' == c) || (')' == c))
}

func isBinOp(ev termbox.Event) bool {
	t, c := ev.Type, ev.Ch
	return (t == termbox.EventKey) && (('/' == c) || ('=' == c) || ('*' == c) || ('+' == c) || ('-' == c) || ('^' == c) || ('.' == c) || ('_' == c) || ('>' == c) || ('<' == c))
}

func isToken(ev termbox.Event) bool {
	t, c := ev.Type, ev.Ch
	return (t == termbox.EventKey) && (('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9'))
}

func isArrow(t termbox.EventType, k termbox.Key) bool {
	return (t == termbox.EventKey) &&
		(k == termbox.KeyArrowUp ||
			k == termbox.KeyArrowDown ||
			k == termbox.KeyArrowLeft ||
			k == termbox.KeyArrowRight)
}

func handlers(ui *UIState) {
	more := true
	for more {
		more, ui = handlersPoll(ui)
	}
}

func main() {
	err := termbox.Init()
	if err == nil {
		ui := newUIState()
		redraw(ui)
		setup()
		handlers(ui)
	} else {
		panic(err)
	}
	termbox.Close()
}
