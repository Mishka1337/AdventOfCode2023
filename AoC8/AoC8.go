package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Step int

const (
	StepLeft Step = iota
	StepRight
)

type Node struct {
	Val          string
	Left         string
	Right        string
	PotentialEnd bool
}

type Cycle struct {
	ReachDistance int
	Length        int
}

type Map struct {
	Nodes map[string]Node
	Steps []Step
}

func main() {
	input := ReadInput("./input.txt")
	parsed := ParseInput(input)
	fmt.Println(Sol1(parsed))
	fmt.Println(Sol2(parsed))
}

func ReadInput(inputFilename string) []string {
	result := make([]string, 0)

	file, err := os.Open(inputFilename)
	if err != nil {
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}

func ParseInput(input []string) Map {
	res := Map{}
	res.Nodes = make(map[string]Node)

	stepsString := input[0]
	for _, c := range stepsString {
		switch c {
		case 'R':
			res.Steps = append(res.Steps, StepRight)
		case 'L':
			res.Steps = append(res.Steps, StepLeft)
		}
	}
	nodeStrings := input[2:]
	for _, s := range nodeStrings {
		node := Node{}

		splitedOnEqual := strings.Split(s, "=")
		nodeVal := splitedOnEqual[0]
		nodeVal = strings.ReplaceAll(nodeVal, " ", "")
		node.Val = nodeVal

		leafsString := splitedOnEqual[1]
		leafsString = strings.ReplaceAll(leafsString, " ", "")
		leafsString = strings.ReplaceAll(leafsString, "(", "")
		leafsString = strings.ReplaceAll(leafsString, ")", "")
		splitedOnComa := strings.Split(leafsString, ",")

		leftVal := splitedOnComa[0]
		rightVal := splitedOnComa[1]
		node.Left = leftVal
		node.Right = rightVal

		if strings.HasSuffix(nodeVal, "Z") {
			node.PotentialEnd = true
		}

		res.Nodes[node.Val] = node
	}

	return res
}

func Sol1(m Map) int {
	res := 0
	stepsLen := len(m.Steps)
	curStep := 0
	curNodeVal := "AAA"

	for {
		if curNodeVal == "ZZZ" {
			break
		}
		curStepType := m.Steps[curStep]
		curNode := m.Nodes[curNodeVal]
		switch curStepType {
		case StepLeft:
			curNodeVal = curNode.Left
		case StepRight:
			curNodeVal = curNode.Right
		}

		res++
		curStep = (curStep + 1) % stepsLen
	}
	return res
}

func Sol2(m Map) int {
	var curNodes []Node
	var cycleLens []int

	for _, v := range m.Nodes {
		if strings.HasSuffix(v.Val, "A") {
			curNodes = append(curNodes, v)
		}
	}

	for _, v := range curNodes {
		cycle := DetectCycle(v, m)
		cycleLens = append(cycleLens, cycle.Length)
	}

	lcm := LCM(cycleLens[0], cycleLens[1], cycleLens[2:]...)

	return lcm

}

//for research purposes
//turns out, that cyclelength and reachdestination is equal

func DetectCycle(n Node, m Map) Cycle {
	res := Cycle{}
	reachDistance := 0
	cycleLen := 0
	curNode := n
	curStep := 0
	stepsLen := len(m.Steps)
	for {
		if curNode.PotentialEnd {
			break
		}
		curStepType := m.Steps[curStep]
		switch curStepType {
		case StepLeft:
			curNode = m.Nodes[curNode.Left]
		case StepRight:
			curNode = m.Nodes[curNode.Right]
		}

		curStep = (curStep + 1) % stepsLen
		reachDistance++
	}

	onFirstNode := true

	for {
		if !onFirstNode && curNode.PotentialEnd {
			break
		}
		onFirstNode = false

		curStepType := m.Steps[curStep]
		switch curStepType {
		case StepLeft:
			curNode = m.Nodes[curNode.Left]
		case StepRight:
			curNode = m.Nodes[curNode.Right]
		}

		curStep = (curStep + 1) % stepsLen
		cycleLen++
	}

	res.ReachDistance = reachDistance
	res.Length = cycleLen

	return res
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
