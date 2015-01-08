package libbisect

import "testing"

//Apply a test iteration and say what outcome is expected
func (bs *BisectState) resultIteration(t *testing.T, outcome Outcome, expected bool) {
	oldMid := bs.GetMid()
	bs.Result(oldMid, outcome)
	newMid := bs.GetMid()
	if (oldMid == newMid) == expected {
		t.Failed()
	}
}

func TestSimpleBisect(t *testing.T) {
	bs := NewBisectState()
	bs.Begin(0, Pass, 4, Fail)
	bs.resultIteration(t, Pass, false)
	bs.resultIteration(t, Pass, false)
	bs.Begin(8, Pass, 8, Fail)
	bs.resultIteration(t, Pass, true)
}
