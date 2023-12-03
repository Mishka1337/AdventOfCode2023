package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type PartNumber struct {
	X, Y, Width, Val int
}

type Gear struct {
	X, Y int
}

type EngineSchema struct {
	PartsMap     [][]bool
	PartsNumbers []PartNumber
	Gears        []Gear
}

func main() {
	input := ReadInput("./input.txt")
	games := ParseInput(input)
	fmt.Println(Sol1(games))
	fmt.Println(Sol2(games))
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

func ParseInput(input []string) EngineSchema {
	var res EngineSchema

	for i, s := range input {
		var parts []bool

		runes := []rune(s)
		j := 0

		for j < len(runes) {
			curElem := runes[j]

			if curElem == '.' {
				parts = append(parts, false)
				j++
				continue
			}

			if curElem == '*' {
				res.Gears = append(res.Gears, Gear{X: j, Y: i})
				parts = append(parts, true)
				j++
				continue
			}

			if curElem < '0' || curElem > '9' {
				parts = append(parts, true)
				j++
				continue
			}

			var number PartNumber
			number.Y = i
			number.X = j
			width := 0
			var numberRunes []rune
			var curNumberRune = runes[number.X+width]

			for curNumberRune >= '0' && curNumberRune <= '9' {
				numberRunes = append(numberRunes, curNumberRune)
				parts = append(parts, false)
				width++
				j++
				if number.X+width < len(runes) {
					curNumberRune = runes[number.X+width]
				} else {
					break
				}
			}

			numberVal, _ := strconv.Atoi(string(numberRunes))

			number.Width = width
			number.Val = numberVal
			res.PartsNumbers = append(res.PartsNumbers, number)
		}
		res.PartsMap = append(res.PartsMap, parts)
	}

	return res
}

func Sol1(input EngineSchema) int {
	res := 0
	for _, number := range input.PartsNumbers {
		isNeighboringPart := false

		for i := 0; i < number.Width; i++ {
			x := number.X + i
			y := number.Y
			isNeighboringPart = isNeighboringPart || IsPartNearToPoint(input.PartsMap, x, y)
		}

		if isNeighboringPart {
			res += number.Val
		}
	}
	return res
}

func IsPartNearToPoint(partMap [][]bool, x int, y int) bool {
	res := false
	maxY := len(partMap) - 1
	maxX := len(partMap[0]) - 1
	minY := 0
	minX := 0
	for i := y - 1; i <= y+1; i++ {
		if i < minY || i > maxY {
			continue
		}
		for j := x - 1; j <= x+1; j++ {
			if j < minX || j > maxX {
				continue
			}

			res = res || partMap[i][j]
		}
	}
	return res
}

func Sol2(input EngineSchema) int {
	res := 0

	for _, gear := range input.Gears {
		res += CalcGearRatio(gear, input.PartsNumbers)
	}

	return res
}

func CalcGearRatio(gear Gear, numbers []PartNumber) int {
	var localNumbers []PartNumber

	for _, number := range numbers {
		minX := number.X - 1
		minY := number.Y - 1
		maxX := number.X + number.Width
		maxY := number.Y + 1

		if gear.X >= minX && gear.X <= maxX && gear.Y >= minY && gear.Y <= maxY {
			localNumbers = append(localNumbers, number)
		}
	}

	if len(localNumbers) != 2 {
		return 0
	}
	res := 1
	for _, number := range localNumbers {
		res *= number.Val
	}
	return res
}
