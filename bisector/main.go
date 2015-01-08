package main

import "bufio"

import "strings"
import "os"
import "fmt"
import "strconv"
import "libbisect"

func bisectWrite(bs libbisect.BisectState, mid int) {
	if bs.IsDone(mid) {
		fmt.Println("a", bs.Low, bs.LowAction, bs.High, bs.HighAction)
	} else {
		fmt.Println("s", mid, bs.Low, bs.LowAction, bs.High, bs.HighAction)
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
				lowAction := tokens[2]
				high, _ := strconv.Atoi(tokens[3])
				highAction := tokens[4]
				bs.Begin(low, lowAction, high, highAction)
				bisectWrite(bs, bs.GetMid())
			} else if cmd == "p" || cmd == "f" {
				bs.Result(bs.GetMid(), cmd)
				bisectWrite(bs, bs.GetMid())
			} else if cmd == "q" {
				running = false
			} else {
				fmt.Println("unknown: ", cmd)
			}
		}
	}
}
