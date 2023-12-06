package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	Time, Distance int
}

func main() {
	input := ReadInput("./input.txt")
	races, race := ParseInput(input)
	fmt.Println(Sol1(races))
	fmt.Println(Sol2(race))
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

func ParseInput(input []string) ([]Race, Race) {
	resForSol1 := ParseForSol1(input)
	resForSol2 := ParseForSol2(input)
	return resForSol1, resForSol2
}

func ParseForSol1(input []string) []Race {
	var res []Race
	times := []int{}
	distances := []int{}
	timeStr := input[0]
	distanceStr := input[1]
	timeStrSplited := strings.Split(timeStr, " ")
	distanceStrSplited := strings.Split(distanceStr, " ")
	for _, s := range timeStrSplited[1:] {
		if s == "" || s == " " {
			continue
		}
		v, _ := strconv.Atoi(s)
		times = append(times, v)
	}

	for _, s := range distanceStrSplited[1:] {
		if s == "" || s == " " {
			continue
		}
		v, _ := strconv.Atoi(s)
		distances = append(distances, v)
	}

	for i := range times {
		race := Race{
			Time:     times[i],
			Distance: distances[i],
		}
		res = append(res, race)
	}

	return res
}

func ParseForSol2(input []string) Race {
	res := Race{}
	timeStr := input[0]
	distanceStr := input[1]
	timeStrSplited := strings.Split(timeStr, " ")
	distanceStrSplited := strings.Split(distanceStr, " ")
	timeResStr := ""
	distResStr := ""

	for _, s := range timeStrSplited[1:] {
		if s == "" || s == " " {
			continue
		}
		timeResStr += s
	}

	for _, s := range distanceStrSplited[1:] {
		if s == "" || s == " " {
			continue
		}
		distResStr += s
	}
	time, _ := strconv.Atoi(timeResStr)
	dist, _ := strconv.Atoi(distResStr)
	res.Time = time
	res.Distance = dist

	return res
}

func Sol1(input []Race) int {
	res := 1
	for _, race := range input {
		res *= FastCalcWinStrategiesCount(race)
	}
	return res
}

func CalcPossibleWinStrategiesCount(r Race) int {
	res := 0
	for i := 0; i <= r.Time; i++ {
		if Predict(r, i) {
			res++
		}
	}
	return res
}

func Predict(r Race, tw int) bool {
	vel := tw
	timeLeft := r.Time - tw
	dist := vel * timeLeft
	return dist > r.Distance
}

func FastCalcWinStrategiesCount(r Race) int {
	D := r.Time*r.Time - 4*r.Distance
	if D < 0 {
		return 0
	}
	if D == 0 {
		return 1
	}
	k1 := float64(r.Time) + (math.Sqrt(float64(D)))
	k1 = k1 / 2
	k2 := float64(r.Time) - (math.Sqrt(float64(D)))
	k2 = k2 / 2
	if k2 < 0 {
		k2 = 0
	}
	k1int := int(math.Floor(k1))
	k2int := int(math.Trunc(k2))
	res := k1int - k2int
	return res
}

func Sol2(input Race) int {
	res := FastCalcWinStrategiesCount(input)
	return res
}
