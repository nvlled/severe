package severe

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	//sevtool "github.com/nvlled/severe/tool"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"runtime/debug"
	"testing"
)

const viewW = 3
const viewH = 3

func ignore(_ ...interface{}) {}

func newTestBuffer() [][]rune {
	return [][]rune{
		[]rune("1\n"),
		[]rune("134\n"),
		[]rune("1345\n"),
		[]rune("13456\n"),
		[]rune("13456789012345\n"),
		[]rune("13456\n"),
		[]rune("13\n"),
		[]rune("1345\n"),
		[]rune("\n"),
	}
}

func TestViewport(t *testing.T) {
	buffer := newTestBuffer()
	fn := func() [][]rune { return buffer }
	view := &Viewport{
		w:      viewW,
		h:      viewH,
		bounds: makeBufferBounds(bufferFunc(fn)),
	}
	checkCursor := func(x, y int) {
		cx, cy := view.Cursor()
		if cx != x || cy != y {
			t.Errorf("cursor expected (%d, %d), got (%d, %d)", x, y, cx, cy)
			debug.PrintStack()
		}
	}
	checkOffset := func(x, y int) {
		ox, oy := view.Offset()
		if ox != x || oy != y {
			t.Errorf("offset expected (%d, %d), got (%d, %d)", x, y, ox, oy)
			debug.PrintStack()
		}
	}
	checkBoundsX := func(x int) {
		bx, _ := view.pointBounds()
		if bx != x {
			t.Errorf("bounds expected (%d), got (%d)", x, bx)
			debug.PrintStack()
		}
	}
	reset := func() {
		view.cursX = 0
		view.cursY = 0
		view.offX = 0
		view.offY = 0
	}

	view.CursorRight()
	checkCursor(1, 0)
	checkOffset(0, 0)
	view.CursorDown()
	checkCursor(1, 1)
	checkOffset(0, 0)
	view.CursorRight()
	checkCursor(2, 1)
	checkOffset(0, 0)
	view.CursorRight()
	checkCursor(2, 1)
	checkOffset(1, 0)
	view.CursorUp()
	checkCursor(0, 0)
	checkOffset(1, 0)

	reset()

	view.CursorDown()
	view.CursorRight()
	view.CursorRight()
	checkCursor(2, 1)
	checkOffset(0, 0)
	view.CursorUp()
	checkBoundsX(1)
	checkCursor(1, 0)
	checkOffset(0, 0)

	reset()

	view.CursorDown()
	view.CursorRight()
	view.CursorRight()
	view.CursorRight()
	view.CursorUp()

	reset()

	view.CursorDown()
	view.CursorDown()
	view.CursorRight()
	view.CursorRight()
	view.CursorRight()
	view.CursorRight()
	view.CursorUp()
	fmt.Println(view)
	view.CursorUp()
	fmt.Println(view)
}

func TestListbox(t *testing.T) {
	defer antiFuck()
	term.Init()
	canvas := wind.NewTermCanvas()

	items := ItemSlice([]string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"eight",
		"nine",
		"ten",
		"eleven",
	})

	lbox := Listbox(20, 5, items)
	layer := wind.Vlayer(
		wind.Border('-', '|', lbox),
		wind.Text(`** Ctrl-c to exit`),
		wind.Text(`** Enter to select`),
	)

	drawLayer := func() {
		term.Clear(0, 0)
		layer.Render(canvas)
		term.Flush()
	}

	cancelled := false
	control.New(
		control.TermSource,
		control.Opts{
			EventEnded: func(_ interface{}) { drawLayer() },
			Interrupt: control.Interrupts(
				control.KeyInterrupt(term.KeyEnter),
				func(e interface{}, stop func()) {
					if e, ok := e.(term.Event); ok && e.Key == term.KeyCtrlC {
						cancelled = true
						stop()
					}
				},
			),
		},
		func(flow *control.Flow) {
			drawLayer()
			lbox.Control(flow)
		},
	)

	term.Close()

	if cancelled {
		println("listbox cancelled!")
	} else {
		_, item := lbox.SelectedItem()
		println("listbox selected:", item)
	}
}

func TestTextbox(t *testing.T) {
	defer antiFuck()
	term.Init()
	canvas := wind.NewTermCanvas()

	tbox := Textbox(50, 10)
	buffer := newTestBuffer()
	tbox.buffer = buffer
	layer := wind.Vlayer(
		wind.Border('-', '|', tbox),
		wind.Text("** Arrow keys to move cursor"),
		wind.Text("** Ctrl-c to exit"),
		wind.Text("** Enter key to insert newline"),
		wind.Text("** Backspace probably doesn't work, use delete key"),
	)
	drawLayer := func() {
		term.Clear(0, 0)
		layer.Render(canvas)
		term.Flush()
	}

	control.New(
		control.TermSource,
		control.Opts{
			EventEnded: func(_ interface{}) { drawLayer() },
			Interrupt:  control.KeyInterrupt(term.KeyCtrlC)},
		func(flow *control.Flow) {
			drawLayer()
			tbox.Control(flow)
		})

	term.Close()
}

func TestToolbar(t *testing.T) {
	defer antiFuck()
	term.Init()
	canvas := wind.NewTermCanvas()

	tbar := Toolbar(
		[]rune("abcdef"),
		[]rune("ghijkl"),
		[]rune("@#$$%^"),
	)
	layer := wind.Vlayer(
		wind.Border('-', '|', tbar),
		wind.Text(`
		** Arrow keys to move cursor
		** Enter to select
		** Ctrl-c to cancel`),
	)
	drawLayer := func() {
		term.Clear(0, 0)
		layer.Render(canvas)
		term.Flush()
	}

	selected := false
	control.New(
		control.TermSource,
		control.Opts{
			EventEnded: func(_ interface{}) { drawLayer() },
			Interrupt: control.Interrupts(
				control.KeyInterrupt(term.KeyCtrlC),
				control.TermInterrupt(func(e term.Event, stop func()) {
					if e.Key == term.KeyEnter {
						stop()
						selected = true
					}
				}),
			)},
		func(flow *control.Flow) {
			drawLayer()
			tbar.Control(flow)
		})

	term.Close()
	if selected {
		println("toolbar selected:", string(tbar.Selected()))
	} else {
		println("toolbar cancelled")
	}
}

func TestSevere1(t *testing.T) {
	defer antiFuck()
	term.Init()

	editor := Textbox(20, 10)
	editor.SetBuffer("nope")
	colorList := Listbox(10, 10, ItemSlice([]string{
		"default", "red", "blue", "yellow", "green",
	}))
	editBtn := Label("|edit text|")
	setColorBtn := Label("|set bgcolor|")

	layer := wind.Vlayer(
		wind.Hlayer(
			wind.Vlayer(
				wind.Text("editor"),
				wind.Border('.', '.', editor),
				editBtn,
			),
			wind.SizeW(3, wind.CharBlock(' ')),
			wind.Vlayer(
				wind.Text("color"),
				wind.Border('.', '.', colorList),
				setColorBtn,
			),
		),
		wind.Line('─'),
		wind.Text(`
		** Arrow keys to move focus
		** Enter to control focused component
		** Esc to stop component control
		** Ctrl-c to exit`),
	)

	color := uint16(term.ColorDefault)
	layer = wind.TapRender(layer, func(layer wind.Layer, canvas wind.Canvas) {
		canvas = wind.ChangeDefaultColor(color, 0, canvas)
		layer.Render(canvas)
	})

	canvas := wind.NewTermCanvas()
	drawLayer := func() {
		term.Clear(0, 0)
		layer.Render(canvas)
		term.Flush()
	}

	// TODO:
	// focuser := NewFocuser(extractGroup(layer))
	focuser := NewFocuser(XGroup(CompGroup(editBtn), CompGroup(setColorBtn)))

	drawLayer()

	control.TermStart(
		control.TermSource,
		control.Opts{
			EventEnded: func(_ interface{}) { drawLayer() },
			Interrupt:  control.KeyInterrupt(term.KeyCtrlC),
		},
		func(flow *control.Flow, e term.Event) {
			switch e.Key {
			case term.KeyArrowLeft:
				focuser.FocusLeft()
			case term.KeyArrowRight:
				focuser.FocusRight()
			case term.KeyEnter:
				component := focuser.Current()
				component.Unfocus()
				drawLayer()

				// *** TODO:
				//flow.New(
				//	control.Opts{Interrupt: control.KeyInterrupt(term.KeyEsc)},
				//	func(flow *control.Flow) {
				//		component.Control(flow)
				//      editBtn.handler = func() { }
				//	})

				if component == editBtn {
					flow.New(
						control.Opts{Interrupt: control.KeyInterrupt(term.KeyEsc)},
						func(flow *control.Flow) {
							editor.Control(flow)
						})
				} else if component == setColorBtn {
					flow.New(
						control.Opts{Interrupt: control.KeyInterrupt(term.KeyEnter)},
						func(flow *control.Flow) {
							colorList.Control(flow)
							_, colorName := colorList.SelectedItem()
							color = uint16(colorValue(colorName))
						})
				}
				component.Focus()
			}
		},
	)

	term.Close()
}

func colorValue(name string) term.Attribute {
	switch name {
	case "default":
		return term.ColorWhite
	case "red":
		return term.ColorRed
	case "blue":
		return term.ColorBlue
	case "yellow":
		return term.ColorYellow
	case "green":
		return term.ColorGreen
	}
	return term.ColorDefault
}

// editor         background                  `
// ------------   ---------                   `
// |          |   |red    |                   `
// |          |   |blue   |                   `
// |          |   |green  |                   `
// |          |   |yellow |                   `
// ------------   ---------                   `
//                                            `
// |edit text|    |set bgcolor|               `
//                                            `

func TestFocusComponents(t *testing.T) {
}

func antiFuck() {
	err := recover()
	if err != nil {
		term.Close()
		fmt.Println("****error:", err)
		debug.PrintStack()
	}
}
