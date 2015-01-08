package libbisect

type BisectState struct {
	Low        int
	LowAction  string
	High       int
	HighAction string
}

func NewBisectState() BisectState {
	return BisectState{0, "p", 100, "f"}
}

func (this *BisectState) Begin(low int, lowAction string, high int, highAction string) {
	this.Low = low
	this.LowAction = lowAction
	this.High = high
	this.HighAction = highAction
}

func (this *BisectState) GetMid() int {
	return (this.Low + this.High) / 2
}

func (this *BisectState) Result(mid int, midAction string) {
	if midAction == this.LowAction {
		this.Low = mid
		this.LowAction = midAction
	} else if midAction == this.HighAction {
		this.High = mid
		this.HighAction = midAction
	}
}

func (this *BisectState) IsDone(mid int) bool {
	return this.Low == mid || mid == this.High
}
