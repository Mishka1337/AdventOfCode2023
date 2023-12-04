package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	WinningNumbers map[int]bool
	ActualNumbers  []int
	Count          int
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

func ParseInput(input []string) []Card {
	var res []Card
	for _, s := range input {
		card := Card{
			WinningNumbers: make(map[int]bool),
			Count:          1,
		}
		_, nums, _ := strings.Cut(s, ":")
		winningNumsStr, actualNumsStr, _ := strings.Cut(nums, "|")

		winningNums := strings.Split(winningNumsStr, " ")
		actualNums := strings.Split(actualNumsStr, " ")
		for _, wn := range winningNums {
			clear := strings.ReplaceAll(wn, " ", "")
			if clear == "" {
				continue
			}
			num, _ := strconv.Atoi(clear)
			card.WinningNumbers[num] = true
		}

		for _, an := range actualNums {
			clear := strings.ReplaceAll(an, " ", "")
			if clear == "" {
				continue
			}
			num, _ := strconv.Atoi(clear)
			card.ActualNumbers = append(card.ActualNumbers, num)
		}

		res = append(res, card)
	}

	return res
}

func Sol1(cards []Card) int {
	res := 0
	for _, card := range cards {
		hits := CalcMatches(card)
		if hits == 0 {
			continue
		}

		cardScore := 1 << (hits - 1)

		res += cardScore
	}
	return res
}

func CalcMatches(card Card) int {
	hits := 0
	for _, num := range card.ActualNumbers {
		if card.WinningNumbers[num] {
			hits++
		}
	}
	return hits
}

func Sol2(cards []Card) int {
	res := 0

	for i, card := range cards {
		matches := CalcMatches(card)
		for k := 0; k < card.Count; k++ {
			for j := 1; j <= matches; j++ {
				cards[i+j].Count++
			}
			res++
		}
	}

	return res
}
