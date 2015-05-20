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
type Textbox struct {
	Focusable
	buffer [][]rune
	view   *Viewport
}

func NewTextbox(w, h int) *Textbox {
	tbox := &Textbox{
		buffer: nil,
		view: &Viewport{
			h: h,
			w: w,
		},
	}
	tbox.view.bounds = makeBufferBounds(tbox)
	return tbox
}

func Textfield(w int) *Textbox {
	return NewTextbox(w, 1)

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

func (tbox *Textbox) Buffer() [][]rune {
	return tbox.buffer
}

func (tbox *Textbox) SetBuffer(text string) {
	var buffer [][]rune
	for _, line := range strings.Split(text, "\n") {
		buffer = append(buffer, []rune(line+"\n"))
	}
	buffer = append(buffer, []rune("\n"))
	tbox.buffer = buffer
}

func (tbox *Textbox) Width() size.T {
	return size.Const(tbox.view.w)
}

func (tbox *Textbox) Height() size.T {
	return size.Const(tbox.view.h)
}

func (tbox *Textbox) Render(canvas wind.Canvas) {
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

func (tbox *Textbox) InsertChar(ch rune) {
	x, y := tbox.view.Point()
	line := tbox.buffer[y]
	rest := line[x:]
	line = line[:x]

	rest = append([]rune{ch}, rest...)
	line = append(line, rest...)
	tbox.buffer[y] = line
	tbox.view.CursorRight()
}

func (tbox *Textbox) InsertNewline() {
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

func (tbox *Textbox) DeleteBack() {
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

func (tbox *Textbox) CursorUp()    { tbox.view.CursorUp() }
func (tbox *Textbox) CursorDown()  { tbox.view.CursorDown() }
func (tbox *Textbox) CursorLeft()  { tbox.view.CursorLeft() }
func (tbox *Textbox) CursorRight() { tbox.view.CursorRight() }

func (tbox *Textbox) Control(flow *control.Flow) {
	flow.TermTransfer(control.Opts{}, func(_ *control.Flow, e term.Event) {
		if e.Ch != 0 {
			tbox.InsertChar(e.Ch)
		} else {
			switch e.Key {
			case term.KeyEnter:
				tbox.InsertNewline()
			case term.KeyDelete:
				tbox.DeleteBack()
			case term.KeySpace:
				tbox.InsertChar(' ')
			case term.KeyArrowDown:
				tbox.CursorDown()
			case term.KeyArrowRight:
				tbox.CursorRight()
			case term.KeyArrowLeft:
				tbox.CursorLeft()
			case term.KeyArrowUp:
				tbox.CursorUp()
			}
		}
	})
}
