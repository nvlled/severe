package severe

import (
	"fmt"
)

type Viewport struct {
	w, h   int
	offX   int
	offY   int
	cursX  int
	cursY  int
	bounds func(x, y int) (int, int)
}

func (view *Viewport) Offset() (int, int) {
	return view.offX, view.offY
}

func (view *Viewport) Point() (int, int) {
	return view.cursX + view.offX, view.cursY + view.offY
}

func (view *Viewport) Cursor() (int, int) {
	return view.cursX, view.cursY
}

func (view *Viewport) Size() (int, int) {
	return view.w, view.h
}

func (view *Viewport) CursorHome() {
	view.cursX = 0
	view.offX = 0
}

func (view *Viewport) CursorLeft() {
	if view.cursX > 0 {
		view.cursX--
	} else if view.offX > 0 {
		view.offX--
	}
}

func (view *Viewport) CursorUp() {
	if view.cursY > 0 {
		view.cursY--
	} else if view.offY > 0 {
		view.offY--
	}
	view.repositionCursor()
}

func (view *Viewport) CursorRight() {
	cx, cy := view.Cursor()
	boundsX, _ := view.bounds(cx+view.offX, cy+view.offY)
	if view.cursX >= boundsX-view.offX {
		return
	}
	if view.cursX < view.w-1 {
		view.cursX++
	} else if view.offX+view.w <= boundsX {
		view.offX++
	}
}

func (view *Viewport) CursorDown() {
	cx, cy := view.Cursor()
	_, boundY := view.bounds(cx+view.offX, cy+view.offY)
	if cy >= boundY-view.offY-1 {
		return
	}
	if view.cursY < view.h-1 {
		view.cursY++
	} else if view.offY+view.h-1 < boundY-1 {
		view.offY++
	}
	view.repositionCursor()
}

func (view *Viewport) pointBounds() (int, int) {
	x, y := view.Point()
	return view.bounds(x, y)
}

func (view *Viewport) setCursorX(x int) {
	view.cursX += x
	//view.moveCursorRight()
	view.repositionCursor()
}

func (view *Viewport) repositionCursor() {
	boundsX, _ := view.pointBounds()
	if view.offX+view.cursX >= boundsX {
		view.cursX = boundsX - view.offX
	}
	view.FocusCursor()
}

func (view *Viewport) FocusCursor() {
	boundsX, _ := view.pointBounds()
	if view.cursX > view.w {
		cursX := view.cursX + view.offX
		view.offX = boundsX - view.w
		view.cursX = cursX - view.offX
	}
	if view.cursX < 0 {
		view.offX = view.offX + view.cursX
		view.cursX = 0
	}
}

func (view *Viewport) String() string {
	boundsX, _ := view.pointBounds()
	return fmt.Sprintf("cursor(%d, %d); offset(%d, %d); boundsX(%d)",
		view.cursX, view.cursY, view.offX, view.offY, boundsX)
}
