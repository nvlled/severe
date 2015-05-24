package severe

import (
	term "github.com/nsf/termbox-go"
	//"github.com/nvlled/severe/tool"
	"github.com/nvlled/control"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

type KeyMap map[term.Key]func()

type GroupType int

const (
	GtypeX = iota
	GtypeY
	GtypeComp
	GtypeNil
)

type Group interface {
	Group()
	Parent() Group
	Next() Group
	Prev() Group
	Children() []Group

	SetParent(g Group)
	SetNext(g Group)
	SetPrev(g Group)
	Gtype() GroupType
}

type ComponentGroup interface {
	Group
	Component() Component
}

type nilGroup struct {
	parent Group
	next   Group
	prev   Group
	elem   Component
}

// good luck inventing nil types for every interface!12311!
var NilGroup = &nilGroup{}

func (_ *nilGroup) Group()            {}
func (_ *nilGroup) Parent() Group     { return NilGroup }
func (_ *nilGroup) Next() Group       { return NilGroup }
func (_ *nilGroup) Prev() Group       { return NilGroup }
func (_ *nilGroup) Children() []Group { return nil }

func (_ *nilGroup) SetParent(parent Group) {}
func (_ *nilGroup) SetNext(next Group)     {}
func (_ *nilGroup) SetPrev(prev Group)     {}
func (_ *nilGroup) Gtype() GroupType       { return GtypeNil }

type gcomp struct {
	parent    Group
	next      Group
	prev      Group
	component Component
	gtype     GroupType
}

func (c *gcomp) Group()               {}
func (c *gcomp) Parent() Group        { return c.parent }
func (c *gcomp) Next() Group          { return c.next }
func (c *gcomp) Prev() Group          { return c.prev }
func (c *gcomp) Children() []Group    { return nil }
func (c *gcomp) Component() Component { return c.component }

func (c *gcomp) SetParent(parent Group) { c.parent = parent }
func (c *gcomp) SetNext(next Group)     { c.next = next }
func (c *gcomp) SetPrev(prev Group)     { c.prev = prev }
func (_ *gcomp) Gtype() GroupType       { return GtypeComp }

type group struct {
	parent   Group
	next     Group
	prev     Group
	children []Group
	gtype    GroupType
}

func (g *group) Group()            {}
func (g *group) Parent() Group     { return g.parent }
func (g *group) Next() Group       { return g.next }
func (g *group) Prev() Group       { return g.prev }
func (g *group) Children() []Group { return g.children }

func (g *group) SetParent(parent Group) { g.parent = parent }
func (g *group) SetNext(next Group)     { g.next = next }
func (g *group) SetPrev(prev Group)     { g.prev = prev }
func (g *group) Gtype() GroupType       { return g.gtype }

func createGroup(gtype GroupType, elems []Group) Group {
	g := &group{
		gtype:    gtype,
		children: elems,
		parent:   NilGroup,
		next:     NilGroup,
		prev:     NilGroup,
	}

	var lastchild Group = NilGroup
	for _, elem := range elems {
		elem.SetParent(g)
		if lastchild != NilGroup {
			lastchild.SetNext(elem)
		}
		elem.SetPrev(lastchild)
		lastchild = elem
	}
	return g
}

func XGroup(elems ...Group) Group {
	return createGroup(GtypeX, elems)
}

func YGroup(elems ...Group) Group {
	return createGroup(GtypeY, elems)
}

func CompGroup(comp Component) Group {
	return &gcomp{
		component: comp,
		gtype:     GtypeComp,
		parent:    NilGroup,
		next:      NilGroup,
		prev:      NilGroup,
	}
}

type Focusable struct {
	focused bool
}

func (f *Focusable) Focus()          { f.focused = true }
func (f *Focusable) Unfocus()        { f.focused = false }
func (f *Focusable) IsFocused() bool { return f.focused }

type focusIndex map[*group]Group

type Focuser struct {
	current   Group
	lastFocus focusIndex
}

func NewFocuser(g Group) *Focuser {
	focuser := &Focuser{
		current:   g,
		lastFocus: make(focusIndex),
	}
	focuser.focusFirstComp()
	return focuser
}

func (foc *Focuser) focusFirstComp() {
	group := foc.current
	if group.Gtype() != GtypeComp {
		group = foc.searchFirstComponent(foc.current)
	}
	if group, ok := group.(ComponentGroup); ok {
		foc.current = group
		group.Component().Focus()
	}
}

func (foc *Focuser) setLastFocus() {
	current := foc.current
	for {
		parent := current.Parent()
		if parent == NilGroup {
			break
		}
		if g, ok := parent.(*group); ok {
			foc.lastFocus[g] = current
		}
		current = parent
	}
}

func (foc *Focuser) Current() Component {
	if cgroup, ok := foc.current.(ComponentGroup); ok {
		return cgroup.Component()
	}
	return NilComponent
}

func (foc *Focuser) setCurrent(group Group) Component {
	if group, ok := group.(ComponentGroup); ok {
		if current := foc.Current(); current != nil {
			current.Unfocus()
		}
		foc.current = group
		foc.setLastFocus()
		comp := group.Component()
		comp.Focus()
		return comp
	}
	return NilComponent
}

func (foc *Focuser) FocusUp() Component {
	group := foc.searchPrev(GtypeY, foc.current)
	return foc.setCurrent(group)
}

func (foc *Focuser) FocusDown() Component {
	group := foc.searchNext(GtypeY, foc.current)
	return foc.setCurrent(group)
}

func (foc *Focuser) FocusLeft() Component {
	group := foc.searchPrev(GtypeX, foc.current)
	return foc.setCurrent(group)
}

func (foc *Focuser) FocusRight() Component {
	group := foc.searchNext(GtypeX, foc.current)
	return foc.setCurrent(group)
}

func (foc *Focuser) searchPrev(gtype GroupType, group Group) Group {
	if group == NilGroup {
		return group
	}

	parentType := group.Parent().Gtype()
	if parentType == gtype && group.Prev() != NilGroup {
		return foc.searchLastComponent(group.Prev())
	}

	return foc.searchPrev(gtype, group.Parent())
}

func (foc *Focuser) searchNext(gtype GroupType, group Group) Group {
	if group == NilGroup {
		return group
	}

	parentType := group.Parent().Gtype()
	if parentType == gtype && group.Next() != NilGroup {
		return foc.searchFirstComponent(group.Next())
	}

	return foc.searchNext(gtype, group.Parent())
}

func (foc *Focuser) searchComponent(g Group, indexOf func([]Group) int) Group {
	for {
		gt := g.Gtype()
		if gt == GtypeComp || gt == GtypeNil {
			break
		}

		if subg, ok := g.(*group); ok {
			if g_, ok := foc.lastFocus[subg]; ok {
				g = g_
				continue
			}
		}

		children := g.Children()
		n := len(children)
		if n == 0 {
			g = NilGroup
			break
		}
		i := indexOf(children)
		g = children[i]
	}
	return g
}

func (foc *Focuser) searchFirstComponent(group Group) Group {
	return foc.searchComponent(group, func(_ []Group) int {
		return 0
	})
}

func (foc *Focuser) searchLastComponent(group Group) Group {
	return foc.searchComponent(group, func(children []Group) int {
		return len(children) - 1
	})
}

type Component interface {
	//Group
	wind.Layer
	//Control(tun *tool.Tun)
	Control(flow *control.Flow)

	//Name() string
	//SetName(name string) Component

	Focus()
	Unfocus()
	IsFocused() bool
	//IsFocusable() bool

	//Label(label string)
	//Unlabel()
}

type nilComp struct {
}

func (_ *nilComp) Width() size.T           { return size.Const(0) }
func (_ *nilComp) Height() size.T          { return size.Const(0) }
func (_ *nilComp) Render(_ wind.Canvas)    {}
func (_ *nilComp) Control(_ *control.Flow) {}

func (_ *nilComp) Focus()   {}
func (_ *nilComp) Unfocus() {}

func (_ *nilComp) IsFocused() bool {
	return false
}

var NilComponent = new(nilComp)

type Sizable struct {
	w, h     int
	AutoSize bool
}

func (s *Sizable) Width() size.T {
	return size.Const(s.w)
}

func (s *Sizable) Height() size.T {
	return size.Const(s.h)
}

func (s *Sizable) Size() (int, int) {
	return s.w, s.h
}

func (s *Sizable) SetSize(w, h int) {
	if s.AutoSize {
		s.w, s.h = w, h
	}
}
