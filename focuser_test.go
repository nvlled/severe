package severe

import (
	"fmt"
	"testing"
)

type testComp struct {
	*nilComp
	name string
}

func (t *testComp) String() string {
	return "comp{" + t.name + "}"
}

// ---------------------------------------
// |                a                    |
// ---------------------------------------
// |   b     |       c        |     d    |
// ---------------------------------------
// |                 |         |         |
// |       e         |         |   g     |
// |                 |   f     |         |
// |-----------------|         |---------|
// |                 |         |         |
// |       h         |         |   i     |
// |                 |         |         |
// ---------------------------------------
// |                  j                  |
// ---------------------------------------
func TestFocusing(t *testing.T) {
	a := &testComp{name: "a"}
	b := &testComp{name: "b"}
	c := &testComp{name: "c"}
	d := &testComp{name: "d"}
	e := &testComp{name: "e"}
	f := &testComp{name: "f"}
	g := &testComp{name: "g"}
	h := &testComp{name: "h"}
	i := &testComp{name: "i"}
	j := &testComp{name: "j"}

	cg := CompGroup
	group := YGroup(
		cg(a),
		XGroup(cg(b), cg(c), cg(d)),
		XGroup(
			YGroup(cg(e), cg(h)),
			cg(f),
			YGroup(cg(g), cg(i)),
		),
		cg(j),
	)

	focuser := NewFocuser(group)

	type entry struct {
		dir      string
		expected string
	}

	tests := []entry{
		entry{"down", "b"},
		entry{"up", "a"},
		entry{"down", "b"},
		entry{"right", "c"},
		entry{"up", "a"},
		entry{"down", "c"},
		entry{"right", "d"},
		entry{"right", "d"},
		entry{"down", "e"},
		entry{"up", "d"},
		entry{"down", "e"},
		entry{"down", "h"},
		entry{"right", "f"},
		entry{"right", "g"},
		entry{"down", "i"},
		entry{"left", "f"},
		entry{"left", "h"},
		entry{"down", "j"},
		entry{"left", "j"},
		entry{"right", "j"},
		entry{"down", "j"},
		entry{"up", "h"},
		entry{"right", "f"},
		entry{"down", "j"},
		entry{"up", "f"},
		entry{"up", "d"},
		entry{"up", "a"},
	}

	for _, e := range tests {
		switch e.dir {
		case "up":
			focuser.FocusUp()
		case "down":
			focuser.FocusDown()
		case "left":
			focuser.FocusLeft()
		case "right":
			focuser.FocusRight()
		}
		comp := focuser.Current().(*testComp)
		if comp.name != e.expected {
			t.Errorf("expected: %v, got %v", e.expected, comp.name)
		}
		fmt.Printf("%s -> %s\n", e.dir, comp.name)
	}

}
