package main

import (
	"bufio"
	"fmt"
	"os"
)

type Vec3 struct {
	X, Y int
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

func ParseInput(input []string) [][]bool {
	var res [][]bool

	for _, s := range input {
		var row []bool
		for _, c := range s {
			if c == '#' {
				row = append(row, true)
			} else {
				row = append(row, false)
			}

		}
		res = append(res, row)
	}
	return res
}

// naive
func Sol1(starMap [][]bool) int {
	expanded := ExpandMap(starMap)
	stars := ExtractStars(expanded)
	res := 0
	for i := 0; i < len(stars)-1; i++ {
		for j := i + 1; j < len(stars); j++ {
			deltaX := Abs(stars[i].X - stars[j].X)
			deltaY := Abs(stars[i].Y - stars[j].Y)
			res += deltaX + deltaY
		}
	}
	return res
}

func Sol2(starMap [][]bool) int {
	res := 0
	stars := ExtractStars(starMap)
	emptyCols := ExtractEmptyCols(starMap)
	emptyRows := ExtractEmptyRows(starMap)
	for i := 0; i < len(stars)-1; i++ {
		for j := i + 1; j < len(stars); j++ {
			deltaX := Abs(stars[i].X - stars[j].X)
			deltaY := Abs(stars[i].Y - stars[j].Y)
			emptyRowsCount := CountEmpty(emptyRows, stars[i].Y, stars[j].Y)
			emptyColsCount := CountEmpty(emptyCols, stars[i].X, stars[j].X)

			deltaX -= emptyColsCount
			deltaY -= emptyRowsCount

			deltaX += emptyColsCount * 1_000_000
			deltaY += emptyRowsCount * 1_000_000

			res += deltaX + deltaY
		}
	}
	return res
}

func ExpandMap(starMap [][]bool) [][]bool {
	res := ExpandByRows(starMap)
	res = ExpandByCols(res)
	return res
}

func ExpandByRows(starMap [][]bool) [][]bool {
	var res [][]bool
	for _, oldRow := range starMap {
		var newRow []bool
		newRow = append(newRow, oldRow...)
		res = append(res, newRow)
		if IsRowEmpty(newRow) {
			res = append(res, newRow)
		}
	}
	return res
}

func ExpandByCols(starMap [][]bool) [][]bool {
	colsCount := len(starMap[0])
	var res [][]bool
	for i := 0; i < colsCount; i++ {
		res = AppendCol(res, starMap, i)
		if IsColEmpty(starMap, i) {
			res = AppendCol(res, starMap, i)
		}
	}
	return res
}

func IsRowEmpty(row []bool) bool {
	empty := true
	for _, c := range row {
		if c {
			empty = false
			break
		}
	}
	return empty
}

func AppendCol(acc, starMap [][]bool, col int) [][]bool {
	var res [][]bool
	if len(acc) == 0 {
		acc = make([][]bool, len(starMap))
	}
	res = acc
	for i := range starMap {
		res[i] = append(res[i], starMap[i][col])
	}
	return res
}

func IsColEmpty(starMap [][]bool, col int) bool {
	empty := true
	for i := range starMap {
		c := starMap[i][col]
		if c {
			empty = false
			break
		}
	}
	return empty
}

func ExtractStars(starMap [][]bool) []Vec3 {
	var res []Vec3
	for i, row := range starMap {
		for j := range row {
			if !starMap[i][j] {
				continue
			}
			vec := Vec3{
				X: j,
				Y: i,
			}
			res = append(res, vec)
		}
	}
	return res
}

func ExtractEmptyRows(starMap [][]bool) []int {
	var res []int

	for i, row := range starMap {
		if IsRowEmpty(row) {
			res = append(res, i)
		}
	}
	return res
}

func ExtractEmptyCols(starMap [][]bool) []int {
	var res []int
	colCount := len(starMap[0])
	for i := 0; i < colCount; i++ {
		if IsColEmpty(starMap, i) {
			res = append(res, i)
		}
	}
	return res
}

func CountEmpty(empty []int, a, b int) int {
	res := 0
	min := Min(a, b)
	max := Max(a, b)
	for _, v := range empty {
		if v > min && v < max {
			res++
		}
	}

	return res
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
