package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Cell int

const (
	CellEmpty Cell = iota
	CellBall
	CellRock
)

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

func ParseInput(input []string) [][]Cell {
	var res [][]Cell
	for _, s := range input {
		var row []Cell
		for _, c := range s {
			switch c {
			case '.':
				row = append(row, CellEmpty)
			case 'O':
				row = append(row, CellBall)
			case '#':
				row = append(row, CellRock)
			}
		}
		res = append(res, row)
	}
	return res
}

func Sol1(input [][]Cell) int {
	res := 0
	copy := Copy(input)
	maxY := len(copy)

	MoveNorth(copy)

	for y, row := range copy {
		for _, cell := range row {
			if cell != CellBall {
				continue
			}
			res += maxY - y
		}
	}

	return res
}

func MoveNorth(Map [][]Cell) {
	maxY := len(Map)
	maxX := len(Map[0])
	for x := 0; x < maxX; x++ {
		curBarier := 0
		for y := 0; y < maxY; y++ {
			cell := Map[y][x]
			switch cell {
			case CellEmpty:
				continue
			case CellRock:
				curBarier = y + 1
			case CellBall:
				Map[y][x] = CellEmpty
				Map[curBarier][x] = CellBall
				curBarier++
			}
		}
	}
}

func MoveWest(Map [][]Cell) {
	maxY := len(Map)
	maxX := len(Map[0])
	for y := 0; y < maxY; y++ {
		curBarier := 0
		for x := 0; x < maxX; x++ {
			cell := Map[y][x]
			switch cell {
			case CellEmpty:
				continue
			case CellRock:
				curBarier = x + 1
			case CellBall:
				Map[y][x] = CellEmpty
				Map[y][curBarier] = CellBall
				curBarier++
			}
		}
	}
}

func MoveSouth(Map [][]Cell) {
	maxY := len(Map)
	maxX := len(Map[0])
	for x := 0; x < maxX; x++ {
		curBarier := maxY - 1
		for y := maxY - 1; y >= 0; y-- {
			cell := Map[y][x]
			switch cell {
			case CellEmpty:
				continue
			case CellRock:
				curBarier = y - 1
			case CellBall:
				Map[y][x] = CellEmpty
				Map[curBarier][x] = CellBall
				curBarier--
			}
		}
	}
}

func MoveEast(Map [][]Cell) {
	maxY := len(Map)
	maxX := len(Map[0])
	for y := 0; y < maxY; y++ {
		curBarier := maxX - 1
		for x := maxX - 1; x >= 0; x-- {
			cell := Map[y][x]
			switch cell {
			case CellEmpty:
				continue
			case CellRock:
				curBarier = x - 1
			case CellBall:
				Map[y][x] = CellEmpty
				Map[y][curBarier] = CellBall
				curBarier--
			}
		}
	}
}

func PrintMap(Map [][]Cell) {
	for _, row := range Map {
		var s string
		for _, cell := range row {
			switch cell {
			case CellEmpty:
				s += "."
			case CellBall:
				s += "O"
			case CellRock:
				s += "#"
			}
		}
		fmt.Println(s)
	}
}

func Copy(input [][]Cell) [][]Cell {
	var res [][]Cell
	for _, row := range input {
		copy := slices.Clone(row)
		res = append(res, copy)
	}

	return res
}

func Equal(xs, ys [][]Cell) bool {
	res := true
	if len(xs) != len(ys) {
		return false
	}
	for i := range xs {
		if !slices.Equal(xs[i], ys[i]) {
			return false
		}
	}
	return res
}

func Find(history [][][]Cell, xs [][]Cell) int {
	for i, prev := range history {
		if Equal(prev, xs) {
			return i
		}
	}

	return -1
}

func Sol2(input [][]Cell) int {
	res := 0
	cycleCount := 1000000000
	current := Copy(input)
	maxY := len(current)

	hist, idx := FindCycle(current)
	cycleCount -= idx
	cycleCount = cycleCount % len(hist)
	current = hist[cycleCount]

	for y, row := range current {
		for _, cell := range row {
			if cell != CellBall {
				continue
			}
			res += maxY - y
		}
	}

	return res
}

func Cycle(state [][]Cell) {
	MoveNorth(state)
	MoveWest(state)
	MoveSouth(state)
	MoveEast(state)
}

// return cycle itself and number of steps to reach entry of cycle
func FindCycle(startState [][]Cell) ([][][]Cell, int) {
	current := Copy(startState)
	var history [][][]Cell
	idx := 0
	for {
		idx = Find(history, current)
		if idx == -1 {
			toSave := Copy(current)
			history = append(history, toSave)
		} else {
			break
		}
		Cycle(current)
	}
	return history[idx:], idx
}
