package libbisect

type Outcome string

const (
	Pass         Outcome = "p"
	Fail         Outcome = "f"
	Undetermined Outcome = "u"
)

type BisectState struct {
	Low         int
	LowOutcome  Outcome
	High        int
	HighOutcome Outcome
}

func GetOutcome(v string) Outcome {
	if v == "p" {
		return Pass
	}
	if v == "f" {
		return Fail
	}
	return Undetermined
}

func NewBisectState() BisectState {
	return BisectState{0, Pass, 100, Fail}
}

func (this *BisectState) Begin(low int, lowOutcome Outcome, high int, highOutcome Outcome) {
	this.Low = low
	this.LowOutcome = lowOutcome
	this.High = high
	this.HighOutcome = highOutcome
}

func (this *BisectState) GetMid() int {
	return (this.Low + this.High) / 2
}

func (this *BisectState) Result(mid int, midOutcome Outcome) {
	if midOutcome == this.LowOutcome {
		this.Low = mid
		this.LowOutcome = midOutcome
	} else if midOutcome == this.HighOutcome {
		this.High = mid
		this.HighOutcome = midOutcome
	}
}

func (this *BisectState) IsDone(mid int) bool {
	return this.Low == mid || mid == this.High
}
