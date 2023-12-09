package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func ParseInput(input []string) [][]int {
	var res [][]int
	for _, s := range input {
		var seq []int
		splitted := strings.Split(s, " ")
		for _, numS := range splitted {
			num, _ := strconv.Atoi(numS)
			seq = append(seq, num)
		}
		res = append(res, seq)
	}
	return res
}

func Sol1(input [][]int) int {
	res := 0
	for _, seq := range input {
		subSeqs := CalcDiffSeqs(seq)
		res += PredictNew(subSeqs)
	}
	return res
}

func CalcDiffSeqs(seq []int) [][]int {
	return CalcDiffSeqsRec(seq, [][]int{seq})
}

func CalcDiffSeqsRec(seq []int, acc [][]int) [][]int {
	isAllZeros := true
	for _, x := range seq {
		if x != 0 {
			isAllZeros = false
			break
		}
	}
	if isAllZeros {
		return acc
	}
	var newSeq []int
	for i := 1; i < len(seq); i++ {
		x0 := seq[i-1]
		x1 := seq[i]
		newSeq = append(newSeq, x1-x0)
	}
	acc = append(acc, newSeq)
	return CalcDiffSeqsRec(newSeq, acc)
}

func PredictNew(subSeqs [][]int) int {
	res := 0
	for i := len(subSeqs) - 2; i >= 0; i-- {
		curSeq := subSeqs[i]
		res += curSeq[len(curSeq)-1]
	}
	return res
}

func Sol2(input [][]int) int {
	res := 0
	for _, seq := range input {
		subSeqs := CalcDiffSeqs(seq)
		res += PredictOld(subSeqs)
	}
	return res
}

func PredictOld(subSeqs [][]int) int {
	res := 0
	for i := len(subSeqs) - 2; i >= 0; i-- {
		curSeq := subSeqs[i]
		x := curSeq[0]
		res = x - res
	}
	return res
}
