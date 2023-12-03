package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Sample struct {
	Red, Green, Blue int
}

type Game struct {
	Id      int
	Samples []Sample
}

var referenceSample Sample = Sample{
	Red:   12,
	Green: 13,
	Blue:  14,
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

func ParseInput(input []string) []Game {
	var res []Game
	for _, s := range input {
		var game Game
		splited := strings.Split(s, ":")
		gameString := splited[0]
		samplesString := splited[1]
		splited = strings.Split(gameString, " ")
		id, _ := strconv.Atoi(splited[1])
		game.Id = id
		splited = strings.Split(samplesString, ";")
		for _, sampleString := range splited {
			sample := ParseSample(sampleString)
			game.Samples = append(game.Samples, sample)
		}
		res = append(res, game)
	}

	return res
}

func ParseSample(input string) Sample {
	var res Sample
	splited := strings.Split(input, ",")
	for _, s := range splited {
		splitedColorCount := strings.Split(s, " ")
		FillSample(&res, splitedColorCount[1], splitedColorCount[2])
	}

	return res
}

func FillSample(sample *Sample, count string, color string) {
	if sample == nil {
		return
	}
	countInt, _ := strconv.Atoi(count)
	switch color {
	case "red":
		sample.Red = countInt
	case "green":
		sample.Green = countInt
	case "blue":
		sample.Blue = countInt
	}
}

func Sol1(games []Game) int {
	res := 0
	for _, game := range games {
		isValid := CheckGame(game)
		if isValid {
			res += game.Id
		}
	}
	return res
}

func CheckGame(game Game) bool {
	res := true
	for _, sample := range game.Samples {
		if !res {
			break
		}
		res = CheckSample(sample)
	}
	return res
}

func CheckSample(sample Sample) bool {
	res := true
	res = res && (sample.Red <= referenceSample.Red)
	res = res && (sample.Green <= referenceSample.Green)
	res = res && (sample.Blue <= referenceSample.Blue)

	return res
}

func Sol2(games []Game) int {
	res := 0
	for _, game := range games {
		res += CalcPower(game)
	}
	return res
}

func CalcPower(game Game) int {
	maxRed, maxGreen, maxBlue := 0, 0, 0
	for _, sample := range game.Samples {
		if maxRed < sample.Red {
			maxRed = sample.Red
		}

		if maxGreen < sample.Green {
			maxGreen = sample.Green
		}

		if maxBlue < sample.Blue {
			maxBlue = sample.Blue
		}
	}

	return maxRed * maxGreen * maxBlue
}
