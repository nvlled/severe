package severe

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
	"strings"
)

type bufferer interface {
	Buffer() [][]rune
}

type bufferFunc func() [][]rune

func (fn bufferFunc) Buffer() [][]rune { return fn() }

// *** there must be at least one newline in the buffer
type textbox struct {
	Focusable
	wrap   bool
	buffer [][]rune
	view   *Viewport
}

func Textbox(w, h int) *textbox {
	tbox := &textbox{
		buffer: nil,
		view: &Viewport{
			h: h,
			w: w,
		},
	}
	tbox.view.bounds = makeBufferBounds(tbox)
	return tbox
}

func Textfield(w int) *textbox {
	return Textbox(w, 1)

}

func makeBufferBounds(buf bufferer) func(int, int) (int, int) {
	return func(x, y int) (int, int) {
		buffer := buf.Buffer()
		if y >= len(buffer) {
			return 0, 0
		}
		row := buffer[y]
		return len(row) - 1, len(buffer) - 1
	}
}

func (tbox *textbox) Buffer() [][]rune {
	return tbox.buffer
}

func (tbox *textbox) SetBuffer(text string) {
	var buffer [][]rune
	for _, line := range strings.Split(text, "\n") {
		buffer = append(buffer, []rune(line+"\n"))
	}
	buffer = append(buffer, []rune("\n"))
	tbox.buffer = buffer
}

func (tbox *textbox) Width() size.T {
	return size.Const(tbox.view.w)
}

func (tbox *textbox) Height() size.T {
	return size.Const(tbox.view.h)
}

func (tbox *textbox) Render(canvas wind.Canvas) {
	if len(tbox.buffer) == 0 {
		return
	}
	view := tbox.view
	ox, oy := view.Offset()
	w, h := view.Size()

	bg := term.ColorDefault
	if tbox.IsFocused() {
		bg = term.ColorRed
	}

	endY := min(oy+h, len(tbox.buffer))
	for y, row := range tbox.buffer[oy:endY] {
		endX := min(ox+w, len(row))
		if ox < len(row) {
			for x, c := range row[ox:endX] {
				canvas.Draw(x, y, c, 0, uint16(bg))
			}
		}
	}
	cx, cy := view.Cursor()
	canvas.Draw(cx, cy, tbox.buffer[oy+cy][ox+cx], 0, uint16(term.ColorBlue))
}

func (tbox *textbox) insertChar(ch rune) {
	x, y := tbox.view.Point()
	line := tbox.buffer[y]
	rest := line[x:]
	line = line[:x]

	rest = append([]rune{ch}, rest...)
	line = append(line, rest...)
	tbox.buffer[y] = line
	tbox.view.CursorRight()
}

func (tbox *textbox) breakline() {
	// *** Assumes line has a line terminator in it
	//     that is to say, ∀line, len(line) >= 1 and ('\n' ∈ line)
	x, y := tbox.view.Point()
	line := tbox.buffer[y]
	newline := copyLine(line[x:])
	line = line[:x]
	tbox.buffer[y] = append(line, '\n')

	if y < len(tbox.buffer)-1 {
		nextlines := copyLines(tbox.buffer[y+1:])
		buffer := tbox.buffer[:y+1]
		buffer = append(buffer, newline)
		tbox.buffer = append(buffer, nextlines...)
	} else {
		tbox.buffer = append(tbox.buffer, newline)
	}

	tbox.view.CursorDown()
	tbox.view.CursorHome()
}

func (tbox *textbox) backspace() {
	x, y := tbox.view.Point()
	line := tbox.buffer[y]
	if x > 0 {
		left := line[0 : x-1]
		right := line[x:]
		tbox.buffer[y] = append(left, right...)
		tbox.view.CursorLeft()
	} else if y > 0 {
		prevline := tbox.buffer[y-1]
		prevline = prevline[:len(prevline)-1] // remove trailing terminator
		line := append(prevline, line...)
		tbox.buffer[y-1] = line

		buffer := tbox.buffer[:y]
		if y+1 < len(tbox.buffer) {
			buffer = append(buffer, tbox.buffer[y+1:]...)
		}
		tbox.buffer = buffer
		tbox.view.CursorUp()

		tbox.view.SetCursorX(len(prevline))
	}
}

// TODO: decouple control
func (tbox *textbox) Control(flow *control.Flow) {
	flow.TermTransfer(control.Opts{}, func(_ *control.Flow, e term.Event) {
		if e.Ch != 0 {
			tbox.insertChar(e.Ch)
		} else {
			switch e.Key {
			case term.KeyEnter:
				tbox.breakline()
			case term.KeyDelete:
				tbox.backspace()
			case term.KeySpace:
				tbox.insertChar(' ')
			case term.KeyArrowDown:
				tbox.view.CursorDown()
			case term.KeyArrowRight:
				tbox.view.CursorRight()
			case term.KeyArrowLeft:
				tbox.view.CursorLeft()
			case term.KeyArrowUp:
				tbox.view.CursorUp()
				//case term.KeyEsc:
				//return true
			}
		}
	})
}
