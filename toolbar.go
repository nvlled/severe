package severe

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

// toolbar
// | @ # $ % ^ & |
// | A B C D E F |

type toolbar struct {
	cursX int
	cursY int
	icons [][]rune
}

func Toolbar(icons ...[]rune) *toolbar {
	return &toolbar{
		cursX: 0,
		cursY: 0,
		icons: icons,
	}
}

func (tb *toolbar) CursorDown() {
	if tb.cursY < len(tb.icons)-1 {
		tb.cursY++
	}
}

func (tb *toolbar) CursorUp() {
	if tb.cursY > 0 {
		tb.cursY--
	}
}

func (tb *toolbar) CursorRight() {
	row := tb.icons[tb.cursY]
	if tb.cursX < len(row)-1 {
		tb.cursX++
	}
}

func (tb *toolbar) CursorLeft() {
	if tb.cursX > 0 {
		tb.cursX--
	}
}

func (tb *toolbar) Width() size.T {
	if len(tb.icons) == 0 {
		return size.Const(0)
	}
	n := len(tb.icons[0])
	return size.Const(n*2 + 1)
}

func (tb *toolbar) Height() size.T {
	return size.Const(len(tb.icons))
}

func (tb *toolbar) Render(canvas wind.Canvas) {
	for y, row := range tb.icons {
		for x, c := range row {
			bg := term.ColorDefault
			if x == tb.cursX && y == tb.cursY {
				bg = term.ColorBlue
			}
			canvas.Draw(x*2, y, ' ', 0, 0)
			canvas.Draw(x*2+1, y, c, uint16(bg), 0)
		}
	}
}

func (tb *toolbar) Selected() rune {
	return tb.icons[tb.cursY][tb.cursX]
}

func (tb *toolbar) DefaultKeys() control.Keymap {
	return control.Keymap{
		term.KeyArrowDown:  func(_ *control.Flow) { tb.CursorDown() },
		term.KeyArrowUp:    func(_ *control.Flow) { tb.CursorUp() },
		term.KeyArrowLeft:  func(_ *control.Flow) { tb.CursorLeft() },
		term.KeyArrowRight: func(_ *control.Flow) { tb.CursorRight() },
		term.KeyEsc:        func(flow *control.Flow) {},
	}
}

func (tb *toolbar) Control(flow *control.Flow) {
	keymap := tb.DefaultKeys()
	opts := control.Opts{Interrupt: control.KeyInterrupt(term.KeyEsc)}
	flow.TermSwitch(opts, keymap)
}
