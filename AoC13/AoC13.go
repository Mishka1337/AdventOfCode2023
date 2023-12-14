package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Pattern [][]int

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

func ParseInput(input []string) []Pattern {
	var res []Pattern
	splitted := SplitByEmptyString(input)
	for _, group := range splitted {
		pattern := ParsePattern(group)
		res = append(res, pattern)
	}

	return res
}

func ParsePattern(input []string) Pattern {
	var res [][]int

	for _, s := range input {
		var row []int
		for _, c := range s {
			switch c {
			case '.':
				row = append(row, 0)
			case '#':
				row = append(row, 1)
			}
		}
		res = append(res, row)
	}

	return res
}

func SplitByEmptyString(input []string) [][]string {
	var res [][]string
	var curGroup []string

	for _, s := range input {
		if s == "" {
			res = append(res, curGroup)
			curGroup = []string{}
			continue
		}
		curGroup = append(curGroup, s)
	}
	res = append(res, curGroup)

	return res
}

func Sol1(input []Pattern) int {
	res := 0
	for _, p := range input {
		res += FindHorReflectionLine(p) * 100
		t := Transpose(p)
		res += FindHorReflectionLine(t)
	}
	return res
}

func Sol2(input []Pattern) int {
	res := 0
	for _, p := range input {
		res += FindHorSludgedReflectionLine(p) * 100
		t := Transpose(p)
		res += FindHorSludgedReflectionLine(t)
	}
	return res
}

func FindHorSludgedReflectionLine(pattern Pattern) int {
	for i := 1; i < len(pattern); i++ {
		res := CheckHorSludgedReflection(pattern, i)
		if res {
			return i
		}
	}
	return 0
}

func CheckHorSludgedReflection(pattern Pattern, linePos int) bool {
	iterMax := Min(linePos, len(pattern)-linePos)
	res := 0
	for i := 0; i < iterMax; i++ {
		row1 := pattern[linePos-i-1]
		row2 := pattern[linePos+i]
		res += Diff(row1, row2)
		if res > 1 {
			return false
		}
	}
	return res == 1
}

func FindHorReflectionLine(pattern Pattern) int {
	for i := 1; i < len(pattern); i++ {
		res := CheckHorReflection(pattern, i)
		if res {
			return i
		}
	}
	return 0
}

func CheckHorReflection(pattern Pattern, linePos int) bool {
	iterMax := Min(linePos, len(pattern)-linePos)
	for i := 0; i < iterMax; i++ {
		row1 := pattern[linePos-i-1]
		row2 := pattern[linePos+i]
		if !slices.Equal(row1, row2) {
			return false
		}
	}
	return true
}

func Transpose(pattern Pattern) Pattern {
	res := make([][]int, len(pattern[0]))
	for i := range res {
		res[i] = make([]int, len(pattern))
	}

	for i, r := range pattern {
		for j := range r {
			res[j][i] = pattern[i][j]
		}
	}
	return res
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// number of diffent cells
func Diff(x, y []int) int {
	res := 0
	for i := range x {
		if x[i] != y[i] {
			res++
		}
	}
	return res
}
