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

// TODO: Handle timeout on tests

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

	lbox := NewListbox(0, 0, items)
	lbox.AutoSize = true

	layer := wind.Vlayer(
		wind.Border('-', '|', wind.SizeW(20, lbox).FreeHeight()),
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
				func(e interface{}, ir control.Irctrl) {
					if e, ok := e.(term.Event); ok && e.Key == term.KeyCtrlC {
						cancelled = true
						ir.Stop()
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

	tbox := NewTextbox(30, 10)
	tbox.AutoSize = false

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

	tbar := NewToolbar(
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
				control.TermInterrupt(func(e term.Event, ir control.Irctrl) {
					if e.Key == term.KeyEnter {
						ir.Stop()
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

func TestLess(t *testing.T) {
	less := NewLess(0, 0)
	less.AutoSize = true

	less.SetText(`*********************************************
********************* ************************************ *****
*  ****   ******* *********** * **** *** ******** *** **********
**     *********  **   ***** ***** *    ********* ** *** *******
*   ***************** *****  ******************** ****** *** ***
***  **** *  ** ************* ******** **   * ***** *******  ***
***** ** **  ** *** ***** ****  ** ****** ******* * ****** ** **
****  ** *** *** ** ** *** *** *** *** *** ****** **  ***** ***
***  **** *  ** *** *******   *****  * ** ******* **** *********
**  ** ***  ****  *************  *********    *********** * ****
*****   ***************** **************************************
*********`)

	layer := wind.Vlayer(
		wind.Border('.', '.', wind.Size(30, 10, less)),
		wind.Text(`
		** Arrow keys to scroll view
		** Ctrl-c to cancel
		`),
	)

	canvas := wind.NewTermCanvas()
	draw := func() {
		layer.Render(canvas)
		term.Flush()
	}

	term.Init()
	draw()
	control.New(
		control.TermSource,
		control.Opts{
			EventEnded: func(_ interface{}) { draw() },
			Interrupt:  control.KeyInterrupt(term.KeyCtrlC),
		},
		func(flow *control.Flow) {
			less.Control(flow)
		},
	)
	term.Close()
}

func TestSevere1(t *testing.T) {
	defer antiFuck()
	term.Init()

	editor := NewTextbox(20, 10)
	editor.SetBuffer("nope")
	colorList := NewListbox(10, 10, ItemSlice([]string{
		"default", "red", "blue", "yellow", "green",
	}))

	editBtn := Button("|edit text|")
	editBtn.Controller = func(flow *control.Flow) {
		flow.New(
			control.Opts{Interrupt: control.KeyInterrupt(term.KeyEsc)},
			func(flow *control.Flow) {
				editor.Control(flow)
			})
	}

	color := uint16(term.ColorDefault)
	setColorBtn := Button("|set bgcolor|")
	setColorBtn.Controller = func(flow *control.Flow) {
		selected := false
		flow.New(
			control.Opts{
				Interrupt: control.TermInterrupt(func(e term.Event, ir control.Irctrl) {
					if e.Key == term.KeyEnter {
						selected = true
						ir.Stop()
					}
				})},
			func(flow *control.Flow) {
				i, _ := colorList.SelectedItem()
				colorList.Control(flow)
				if selected {
					_, colorName := colorList.SelectedItem()
					color = uint16(colorValue(colorName))
				} else {
					colorList.SetIndex(i)
				}
			})
	}

	layer := wind.Vlayer(
		wind.Hlayer(
			wind.Vlayer(
				wind.Text("editor"),
				wind.Border('.', '.', editor),
				editBtn,
			),
			wind.LineV(' '),
			wind.Vlayer(
				wind.Text("color"),
				wind.Border('.', '.', colorList),
				setColorBtn,
			),
		),
		wind.LineH('─'),
		wind.Text(`
		** Arrow keys to move focus
		** Enter to control focused component
		** Esc to stop component control
		** Ctrl-c to exit`),
	)

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
				flow.New(
					control.Opts{Interrupt: control.KeyInterrupt(term.KeyEsc)},
					func(flow *control.Flow) {
						component.Control(flow)
					})

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
