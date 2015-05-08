package severe

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
	"strings"
)

type label struct {
	Focusable
	lines  []string
	width  int
	height int
}

func Label(text string) *label {
	lines := strings.Split(text, "\n")
	maxw := 0
	for _, line := range lines {
		w := len(line)
		if w > maxw {
			maxw = w
		}
	}
	return &label{
		width:  maxw,
		height: len(lines),
		lines:  lines,
	}
}

func (l *label) Width() size.T {
	return size.Const(l.width)
}

func (l *label) Height() size.T {
	return size.Const(l.height)
}

func (l *label) Render(canvas wind.Canvas) {
	bg := term.ColorDefault
	if l.IsFocused() {
		bg = term.ColorRed
	}
	for y, row := range l.lines {
		for x, c := range row {
			canvas.Draw(x, y, c, 0, uint16(bg))
		}
	}
}

func (_ *label) Control(flow *control.Flow) {
}
