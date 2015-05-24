package severe

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
)

// adding items in the list

// cursX is ignored
// │─────────────│
// │1. one       │
// │-------------│
// │2. two       │
// │3. three     │
// │4. four      │
// │-------------│
// │5. five      │
// │6. six       │
// │─────────────│

// TODO: rename to ListItems
type Items interface {
	GetItems() []string
}

type ItemSlice []string

func (items ItemSlice) GetItems() []string {
	return []string(items)
}

type ItemsFn func() []string

func (fn ItemsFn) GetItems() []string { return fn() }

// TODO: rename to List
type Listbox struct {
	Focusable
	Sizable
	focused bool
	TabSym  string

	Items Items
	view  *Viewport
	//OnSelect func(index int, item string)
}

func NewListbox(w, h int, items Items) *Listbox {
	return &Listbox{
		Sizable: Sizable{w: w, h: h},
		Items:   items,
		view: &Viewport{
			bounds: func(_, _ int) (int, int) {
				// note: lbox.items and items may happen
				// to be different if lbox.items is re-assigned
				return 0, len(items.GetItems())
			},
		},
	}
}

func (lbox *Listbox) Render(canvas wind.Canvas) {
	lbox.SetSize(canvas.Dimension())
	lbox.view.SetSize(lbox.Size())

	_, cursY := lbox.view.Cursor()
	_, offY := lbox.view.Offset()
	_, h := lbox.view.Size()
	items := lbox.Items.GetItems()

	endY := min(offY+h, len(items))
	for y, item := range items[offY:endY] {
		var bgColor uint16 = 0
		if lbox.IsFocused() {
			bgColor = uint16(term.ColorRed)
		}
		if y == cursY {
			bgColor = uint16(term.ColorBlue)
		}

		x := 0
		var c rune

		for x, c = range item {
			canvas.Draw(x, y, c, 0, bgColor)
		}

		for x = x + 1; x < canvas.Width(); x++ {
			canvas.Draw(x, y, ' ', 0, bgColor)
		}
	}
}

func (lbox *Listbox) SelectDown() {
	lbox.view.CursorDown()
}

func (lbox *Listbox) SelectUp() {
	lbox.view.CursorUp()
}

func (lbox *Listbox) SelectedItem() (int, string) {
	_, i := lbox.view.Point()
	items := lbox.Items.GetItems()
	if i < 0 || i >= len(items) {
		return -1, ""
	}
	return i, items[i]
}

func (lbox *Listbox) SetIndex(i int) {
	lbox.view.SetCursorY(i)
}

func (lbox *Listbox) DefaultKeys() control.Keymap {
	return control.Keymap{
		term.KeyArrowUp:   func(_ *control.Flow) { lbox.SelectUp() },
		term.KeyArrowDown: func(_ *control.Flow) { lbox.SelectDown() },
		//term.KeyEnter: func() {
		//	if lbox.OnSelect != nil {
		//		items := lbox.items.GetItems()
		//		_, i := lbox.view.cursor()
		//		if i > 0 && i < len(items) {
		//			lbox.OnSelect(i, items[i])
		//		}
		//	}
		//},
	}
}

func (lbox *Listbox) Control(flow *control.Flow) {
	//flow.Switch(lbox.DefaultKeys())
	keymap := lbox.DefaultKeys()
	flow.TermSwitch(control.Opts{}, keymap)
	//flow.TermControl(func(flow *control.Flow, e term.Event) {
	//	if fn, ok := keymap[e.Key]; ok {
	//		fn(flow)
	//	}
	//})
}
