package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Hand struct {
	Cards1    [5]int
	Cards     [5]int
	Bid       int
	HandType  HandType
	HandType1 HandType
}

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfKind
	FullHouse
	FourOfKind
	FiveOfKind
)

var RuneLabelToInt = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var RuneLabelToInt1 = map[rune]int{
	'J': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 11,
	'K': 12,
	'A': 13,
}

func main() {
	input := ReadInput("./input.txt")
	hands := ParseInput(input)
	fmt.Println(Sol1(hands))
	fmt.Println(Sol2(hands))
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

func ParseInput(input []string) []Hand {
	var res []Hand

	for _, s := range input {
		hand := Hand{}
		splitted := strings.Split(s, " ")
		hand.Cards = ParseCards(splitted[0])
		hand.Cards1 = ParseCards1(splitted[0])
		bid, _ := strconv.Atoi(splitted[1])
		hand.Bid = bid

		res = append(res, hand)
	}

	return res
}

func ParseCards(input string) [5]int {
	var res [5]int
	runes := []rune(input)
	for i, r := range runes {
		res[i] = RuneLabelToInt[r]
	}
	return res
}

func ParseCards1(input string) [5]int {
	var res [5]int
	runes := []rune(input)
	for i, r := range runes {
		res[i] = RuneLabelToInt1[r]
	}
	return res
}

func Sol1(hands []Hand) int {
	res := 0
	for i, h := range hands {
		hands[i].HandType = CalcHandType(h)
	}

	slices.SortFunc(hands, HandCmp)

	for i, h := range hands {
		rank := i + 1
		res += rank * h.Bid
	}
	return res
}

func CalcHandType(h Hand) HandType {
	labelCount := CalcLabelCounts(h)

	if IsFiveOfKind(labelCount) {
		return FiveOfKind
	}
	if IsFourOfKind(labelCount) {
		return FourOfKind
	}
	if IsFullHouse(labelCount) {
		return FullHouse
	}
	if IsThreeOfKind(labelCount) {
		return ThreeOfKind
	}
	if IsTwoPair(labelCount) {
		return TwoPair
	}
	if IsPair(labelCount) {
		return OnePair
	}
	return HighCard
}

func HandCmp(x, y Hand) int {
	res := int(x.HandType - y.HandType)
	if res == 0 {
		res = CompareHandsByCardsValue(x, y)
	}

	return res
}

func CompareHandsByCardsValue(x, y Hand) int {
	for i := range x.Cards {
		diff := x.Cards[i] - y.Cards[i]
		if diff == 0 {
			continue
		}
		return diff
	}
	return 0
}

func Sol2(hands []Hand) int {
	res := 0
	for i, h := range hands {
		hands[i].HandType1 = CalcHandType1(h)
	}

	slices.SortFunc(hands, HandCmp1)

	for i, h := range hands {
		rank := i + 1
		res += rank * h.Bid
	}
	return res
}

func HandCmp1(x, y Hand) int {
	res := int(x.HandType1 - y.HandType1)
	if res == 0 {
		res = CompareHandsByCardsValue1(x, y)
	}
	return res
}

func CompareHandsByCardsValue1(x, y Hand) int {
	for i := range x.Cards1 {
		diff := x.Cards1[i] - y.Cards1[i]
		if diff == 0 {
			continue
		}
		return diff
	}
	return 0
}

func CalcHandType1(h Hand) HandType {
	labelCount := CalcLabelCounts1(h)
	res := GetHandTypeByLabelCounts(labelCount)
	jockerCount := labelCount[1]
	for i := jockerCount; i > 0; i-- {
		switch res {
		case FiveOfKind:
			res = FiveOfKind
		case FourOfKind:
			res = FiveOfKind
		case FullHouse:
			res = FourOfKind
		case ThreeOfKind:
			res = FourOfKind
		case TwoPair:
			res = FullHouse
		case OnePair:
			res = ThreeOfKind
		case HighCard:
			res = OnePair
		}
	}
	return res
}

func GetHandTypeByLabelCounts(labelCount map[int]int) HandType {
	if IsFiveOfKind(labelCount) {
		return FiveOfKind
	}
	if IsFourOfKind(labelCount) {
		return FourOfKind
	}
	if IsFullHouse(labelCount) {
		return FullHouse
	}
	if IsThreeOfKind(labelCount) {
		return ThreeOfKind
	}
	if IsTwoPair(labelCount) {
		return TwoPair
	}
	if IsPair(labelCount) {
		return OnePair
	}
	return HighCard
}

func CalcLabelCounts(h Hand) map[int]int {
	res := make(map[int]int)

	for _, v := range h.Cards {
		res[v] = res[v] + 1
	}
	return res
}

func CalcLabelCounts1(h Hand) map[int]int {
	res := make(map[int]int)

	for _, v := range h.Cards1 {
		res[v] = res[v] + 1
	}

	return res
}

func IsFiveOfKind(h map[int]int) bool {
	for i, v := range h {
		if v == 5 && i != 1 {
			return true
		}
	}
	return false
}

func IsFourOfKind(h map[int]int) bool {
	for i, v := range h {
		if v == 4 && i != 1 {
			return true
		}
	}
	return false
}

func IsFullHouse(h map[int]int) bool {
	haveThree := false
	havePair := false
	for i, v := range h {
		if i == 1 {
			continue
		}

		if v == 3 {
			haveThree = true
		}
		if v == 2 {
			havePair = true
		}
	}
	return havePair && haveThree
}

func IsThreeOfKind(h map[int]int) bool {
	for i, v := range h {
		if v == 3 && i != 1 {
			return true
		}
	}
	return false
}

func IsTwoPair(h map[int]int) bool {
	count := 0
	for i, v := range h {
		if i == 1 {
			continue
		}
		if v == 2 {
			count++
		}
	}
	return count == 2
}

func IsPair(h map[int]int) bool {
	for i, v := range h {
		if v == 2 && i != 1 {
			return true
		}
	}
	return false
}
