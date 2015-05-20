package severe

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
	"strings"
)

type Less struct {
	Focusable
	buffer [][]rune
	view   *Viewport
	w      int
	h      int
	maxw   int
}

func NewLess(w, h int) *Less {
	less := &Less{
		w:    w,
		h:    h,
		maxw: 0,
		view: &Viewport{
			h: 1,
			w: 1,
		},
	}
	less.view.bounds = func(_, _ int) (int, int) {
		return less.maxw - less.w + 1, len(less.buffer) - less.h + 1
	}
	return less
}

func (less *Less) Width() size.T {
	return size.Const(less.w)
}

func (less *Less) Height() size.T {
	return size.Const(less.h)
}

func (less *Less) Render(canvas wind.Canvas) {
	canvas.Clear()

	view := less.view
	ox, oy := view.Offset()
	w, h := less.w, less.h

	endY := min(oy+h, len(less.buffer))
	for y, row := range less.buffer[oy:endY] {
		endX := min(ox+w, len(row))
		if ox < len(row) && ox >= 0 {
			for x, c := range row[ox:endX] {
				canvas.Draw(x, y, c, 0, 0)
			}
		}
	}
}

func (less *Less) SetText(text string) {
	var buffer [][]rune
	less.maxw = 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > less.maxw {
			less.maxw = len(line)
		}
		buffer = append(buffer, []rune(line))
	}
	less.buffer = buffer
	less.view.CursorHome()
}

func (less *Less) ScrollUp() {
	less.view.CursorUp()
}

func (less *Less) ScrollDown() {
	less.view.CursorDown()
}

func (less *Less) ScrollLeft() {
	less.view.CursorLeft()
}

func (less *Less) ScrollRight() {
	less.view.CursorRight()
}

func (less *Less) PageUp() {
	view := less.view
	view.offY -= less.h
	if view.offY < 0 {
		view.offY = 0
	}
}

func (less *Less) PageDown() {
	view := less.view
	_, boundsY := view.pointBounds()
	view.offY += less.h
	if view.offY >= boundsY {
		view.offY = boundsY - 1
	}
}

func (less *Less) ScrollStartX() {
	less.view.CursorStartX()
}

func (less *Less) ScrollEndX() {
	less.view.CursorEndX()
}

func (less *Less) ScrollStartY() {
	less.view.CursorStartY()
}

func (less *Less) ScrollEndY() {
	less.view.CursorEndY()
}

func (less *Less) Control(flow *control.Flow) {
	flow.TermTransfer(control.Opts{}, func(flow *control.Flow, e term.Event) {
		switch e.Key {
		case term.KeyArrowUp:
			less.ScrollUp()
		case term.KeyArrowDown:
			less.ScrollDown()
		case term.KeyArrowLeft:
			less.ScrollLeft()
		case term.KeyArrowRight:
			less.ScrollRight()
		case term.KeyCtrlA:
			less.ScrollStartX()
		case term.KeyCtrlE:
			less.ScrollEndX()
		case term.KeyHome:
			less.ScrollStartY()
		case term.KeyEnd:
			less.ScrollEndY()
		}
	})
}
