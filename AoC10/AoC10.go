package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type State struct {
	MainLoop *Node
	PipeMap  [][]rune
}

type Node struct {
	Next         *Node
	Prev         *Node
	IsUp         bool
	IsHorizontal bool
	IsVertical   bool
	IsStart      bool
	X, Y         int
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

func ParseInput(input []string) State {
	res := State{}
	startX := 0
	startY := 0
	for i, s := range input {
		for j, c := range s {
			if c == 'S' {
				startX = j
				startY = i
				break
			}
		}
	}
	pipeMap := [][]rune{}
	for _, s := range input {
		pipeMap = append(pipeMap, []rune(s))
	}
	res.PipeMap = pipeMap

	nextX, nextY := GetStartDirection(startX, startY, pipeMap)
	curX, curY := startX, startY
	startNode := Node{}
	startNode.IsStart = true
	startNode.X = startX
	startNode.Y = startY

	//hardcoded, but not guarantied
	startNode.IsHorizontal = true
	startNode.IsVertical = true
	curNode := &startNode

	for pipeMap[nextY][nextX] != 'S' {
		nextNode := Node{}
		nextNode.Prev = curNode
		curNode.Next = &nextNode
		curNode = &nextNode
		curNode.X = nextX
		curNode.Y = nextY

		tX, tY := nextX, nextY
		nextX, nextY = GetNextPoint(curX, curY, tX, tY, pipeMap)
		curX, curY = tX, tY
		curNode.IsHorizontal = IsHorizontal(curX, curY, pipeMap)
		curNode.IsVertical = IsVertical(curX, curY, pipeMap)
	}
	curNode.Next = &startNode
	startNode.Prev = curNode
	res.MainLoop = &startNode

	return res
}

func IsVertical(x, y int, pipeMap [][]rune) bool {
	pipe := pipeMap[y][x]
	return pipe != '-'
}

func IsHorizontal(x, y int, pipeMap [][]rune) bool {
	pipe := pipeMap[y][x]
	return pipe != '|'
}

func GetStartDirection(startX, startY int, pipeMap [][]rune) (int, int) {
	resX := startX
	resY := startY

	width := len(pipeMap[0])
	depth := len(pipeMap)

	//check from left
	curX := startX - 1
	curY := startY
	validLeft := []rune{
		'L', 'F', '-',
	}
	if curX > 0 && curX < width && curY > 0 && curY < depth {
		if slices.Contains(validLeft, pipeMap[curY][curX]) {
			return curX, curY
		}

	}

	//check from top
	curX = startX
	curY = startY - 1
	validTop := []rune{
		'F', '7', '|',
	}
	if curX > 0 && curX < width && curY > 0 && curY < depth {
		if slices.Contains(validTop, pipeMap[curY][curX]) {
			return curX, curY
		}
	}

	//check from right
	curX = startX + 1
	curY = startY
	validRight := []rune{
		'7', 'J', '-',
	}
	if curX > 0 && curX < width && curY > 0 && curY < depth {
		if slices.Contains(validRight, pipeMap[curY][curX]) {
			return curX, curY
		}
	}

	//check from bottom
	curX = startX
	curY = startY + 1
	validBottom := []rune{
		'L', 'J', '|',
	}
	if curX > 0 && curX < width && curY > 0 && curY < depth {
		if slices.Contains(validBottom, pipeMap[curY][curX]) {
			return curX, curY
		}
	}

	return resX, resY
}

func GetNextPoint(curX, curY, nextX, nextY int, pipeMap [][]rune) (int, int) {
	point := pipeMap[nextY][nextX]
	deltaX := nextX - curX
	deltaY := nextY - curY
	switch point {
	case '|':
		{
			if deltaY > 0 {
				return nextX, nextY + 1
			} else {
				return nextX, nextY - 1
			}
		}
	case '-':
		{
			if deltaX > 0 {
				return nextX + 1, nextY
			} else {
				return nextX - 1, nextY
			}
		}
	case 'L':
		{
			if deltaX != 0 {
				return nextX, nextY - 1
			} else {
				return nextX + 1, nextY
			}
		}
	case 'J':
		{
			if deltaX != 0 {
				return nextX, nextY - 1
			} else {
				return nextX - 1, nextY
			}
		}
	case '7':
		{
			if deltaX != 0 {
				return nextX, nextY + 1
			} else {
				return nextX - 1, nextY
			}
		}
	case 'F':
		{
			if deltaX != 0 {
				return nextX, nextY + 1
			} else {
				return nextX + 1, nextY
			}
		}
	}

	return nextX, nextY
}

func Sol1(state State) int {
	res := 1
	curNode := state.MainLoop.Next
	for !curNode.IsStart {
		curNode = curNode.Next
		res++
	}
	res++
	return res / 2
}

func Sol2(state State) int {
	res := 0

	verticalPoints := GetVertMap(state.MainLoop)
	allPoints := GetLoopMap(state.MainLoop)
	TraverseOrientations(state.MainLoop)
	pipeMap := state.PipeMap
	for i, s := range pipeMap {
		isOutSide := true
		isIntersectedVert := false
		isUp := true
		for j := range s {
			vertNode, isVertical := verticalPoints[i][j]
			if isVertical {
				if !isIntersectedVert {
					isUp = vertNode.IsUp
					isOutSide = false
				} else {
					if vertNode.IsUp != isUp {
						isOutSide = !isOutSide
						isUp = vertNode.IsUp
					}
				}
				isIntersectedVert = true

				continue
			}
			if _, exist := allPoints[i][j]; exist {
				continue
			}
			if !isOutSide {
				res++
			}

		}
	}
	return res
}

func GetLoopMap(startNode *Node) map[int]map[int]*Node {
	res := make(map[int]map[int]*Node)

	curNode := startNode.Next
	if _, exist := res[startNode.Y]; !exist {
		res[startNode.Y] = make(map[int]*Node)
	}
	res[startNode.Y][startNode.X] = startNode

	for !curNode.IsStart {
		if _, exist := res[curNode.Y]; !exist {
			res[curNode.Y] = make(map[int]*Node)
		}
		res[curNode.Y][curNode.X] = curNode

		curNode = curNode.Next
	}

	return res
}

func GetVertMap(startNode *Node) map[int]map[int]*Node {
	res := make(map[int]map[int]*Node)

	curNode := startNode.Next
	if _, exist := res[startNode.Y]; !exist {
		res[startNode.Y] = make(map[int]*Node)
	}
	if startNode.IsVertical {
		res[startNode.Y][startNode.X] = startNode
	}

	for !curNode.IsStart {
		if _, exist := res[curNode.Y]; !exist {
			res[curNode.Y] = make(map[int]*Node)
		}
		if curNode.IsVertical {
			res[curNode.Y][curNode.X] = curNode
		}

		curNode = curNode.Next
	}

	return res
}

func TraverseOrientations(startNode *Node) {
	CalcOrientation(startNode)
	curNode := startNode.Next
	for !curNode.IsStart {
		CalcOrientation(curNode)
		curNode = curNode.Next
	}
}

func CalcOrientation(node *Node) {
	if !node.IsVertical {
		return
	}
	if node.X == node.Next.X && node.Y > node.Next.Y {
		node.IsUp = true
	}

	if node.X == node.Prev.X && node.Y < node.Prev.Y {
		node.IsUp = true
	}
}
