package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var spelledDigits = [10]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func main() {
	input := ReadInput("./input.txt")
	fmt.Println(Sol1(input))
	fmt.Println(Sol2(input))
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

func Sol1(input []string) int {
	res := 0
	for _, s := range input {
		first, last := FindFirstAndLastDigit(s)
		first *= 10
		res += (first + last)
	}
	return res
}

func FindFirstAndLastDigit(s string) (int, int) {
	first, last := 0, 0
	sArr := []rune(s)
	for i := 0; i < len(sArr); i++ {
		c := sArr[i]
		if c >= '0' && c <= '9' {
			first = int(c - '0')
			break
		}
	}

	for i := len(sArr) - 1; i >= 0; i-- {
		c := sArr[i]
		if c >= '0' && c <= '9' {
			last = int(c - '0')
			break
		}
	}

	return first, last
}

func Sol2(input []string) int {
	res := 0
	for _, s := range input {
		first, last := FindFirstAndLastSpelledDigits(s)
		first *= 10
		res += first + last
	}

	return res
}

func FindFirstAndLastSpelledDigits(input string) (int, int) {
	first, last := 0, 0
	minIndex := len(input)
	maxIndex := 0
	indexes := make(map[int]int)

	for i, s := range spelledDigits {
		index := strings.Index(input, s)
		indexes[i] = index
	}

	for k, v := range indexes {
		if v == -1 {
			continue
		}
		if v < minIndex {
			minIndex = v
			first = k
		}
	}

	sArr := []rune(input)
	for i := 0; i < minIndex; i++ {
		c := sArr[i]
		if c >= '0' && c <= '9' {
			first = int(c - '0')
			break
		}
	}
	indexes = make(map[int]int)

	for i, s := range spelledDigits {
		index := strings.LastIndex(input, s)
		indexes[i] = index
	}

	for k, v := range indexes {
		if v == -1 {
			continue
		}
		if v > maxIndex {
			maxIndex = v
			last = k
		}
	}

	for i := len(sArr) - 1; i >= maxIndex; i-- {
		c := sArr[i]
		if c >= '0' && c <= '9' {
			last = int(c - '0')
			break
		}
	}

	return first, last
}
