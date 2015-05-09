package severe

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
	"strings"
)

type button struct {
	Focusable
	lines  []string
	width  int
	height int
}

func Button(text string) *button {
	lines := strings.Split(text, "\n")
	maxw := 0
	for _, line := range lines {
		w := len(line)
		if w > maxw {
			maxw = w
		}
	}
	return &button{
		width:  maxw,
		height: len(lines),
		lines:  lines,
	}
}

func (btn *button) Width() size.T {
	return size.Const(btn.width)
}

func (btn *button) Height() size.T {
	return size.Const(btn.height)
}

func (btn *button) Render(canvas wind.Canvas) {
	bg := term.ColorDefault
	if btn.IsFocused() {
		bg = term.ColorRed
	}
	for y, row := range btn.lines {
		for x, c := range row {
			canvas.Draw(x, y, c, 0, uint16(bg))
		}
	}
}

func (btn *button) Control(flow *control.Flow) {
}
