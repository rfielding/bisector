package main

import "bufio"

import "strings"
import "os"
import "fmt"
import "strconv"
import "libbisect"

func bisectWrite(bs libbisect.BisectState, mid int) {
	if bs.IsDone(mid) {
		fmt.Println("a", bs.Low, bs.LowOutcome, bs.High, bs.HighOutcome)
	} else {
		fmt.Println("s", mid, bs.Low, bs.LowOutcome, bs.High, bs.HighOutcome)
	}
}

func main() {
	bio := bufio.NewReader(os.Stdin)
	bs := libbisect.NewBisectState()
	running := true
	for running {
		line, _, err := bio.ReadLine()
		str := string(line)
		if err == nil {
			//Get a command
			tokens := strings.Split(str, " ")
			cmd := tokens[0]

			if cmd == "n" && len(tokens) == 5 {
				low, _ := strconv.Atoi(tokens[1])
				lowOutcome := libbisect.GetOutcome(tokens[2])
				high, _ := strconv.Atoi(tokens[3])
				highOutcome := libbisect.GetOutcome(tokens[4])
				bs.Begin(low, lowOutcome, high, highOutcome)
				bisectWrite(bs, bs.GetMid())
			} else if cmd == "p" || cmd == "f" {
				midOutcome := libbisect.GetOutcome(cmd)
				bs.Result(bs.GetMid(), midOutcome)
				bisectWrite(bs, bs.GetMid())
			} else if cmd == "q" {
				running = false
			} else {
				fmt.Println("unknown: ", cmd)
			}
		}
	}
}
