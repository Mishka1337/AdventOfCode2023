package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Almanac struct {
	Seeds      []int64
	Maps       []Map
	SeedRanges []Range
}

type Map struct {
	Ranges []MapRange
}

type MapRange struct {
	Offset      int64
	SourceStart int64
	SourceEnd   int64
}

type Range struct {
	Start int64
	End   int64
}

type Set struct {
	Ranges []Range
}

func main() {
	input := ReadInput("./input.txt")
	maps := ParseInput(input)
	fmt.Println(Sol1(maps))
	fmt.Println(Sol2(maps))
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

func ParseInput(input []string) Almanac {
	var maps []Map

	seeds := ParseSeeds(input[0])
	seedRanges := ParseSeedRanges(input[0])

	var curRanges []string
	for _, s := range input[3:] {
		if strings.HasSuffix(s, "map:") {
			continue
		}

		if s == "" {
			curMap := ParseRanges(curRanges)
			maps = append(maps, curMap)
			curRanges = []string{}
			continue
		}
		curRanges = append(curRanges, s)
	}
	lastMap := ParseRanges(curRanges)
	maps = append(maps, lastMap)

	res := Almanac{}
	res.Seeds = seeds
	res.Maps = maps
	res.SeedRanges = seedRanges
	return res
}

func ParseRanges(input []string) Map {
	curMap := Map{}
	for _, s := range input {
		splitted := strings.Split(s, " ")

		destStart, _ := strconv.ParseInt(splitted[0], 0, 64)
		sourceStart, _ := strconv.ParseInt(splitted[1], 0, 64)
		length, _ := strconv.ParseInt(splitted[2], 0, 64)
		curRange := MapRange{
			Offset:      destStart - sourceStart,
			SourceStart: sourceStart,
			SourceEnd:   sourceStart + length - 1,
		}
		curMap.Ranges = append(curMap.Ranges, curRange)
	}
	return curMap
}

func ParseSeeds(input string) []int64 {
	var res []int64
	splitted := strings.Split(input, " ")
	for _, s := range splitted[1:] {
		parsed, _ := strconv.ParseInt(s, 0, 64)
		res = append(res, parsed)
	}
	return res
}

func ParseSeedRanges(input string) []Range {
	var res []Range
	splitted := strings.Split(input, " ")
	for i := 1; i < len(splitted)-1; i += 2 {
		s := splitted[i]
		start, _ := strconv.ParseInt(s, 0, 64)
		s = splitted[i+1]
		length, _ := strconv.ParseInt(s, 0, 64)
		r := Range{
			Start: start,
			End:   start + length - 1,
		}
		res = append(res, r)
	}
	return res
}

func Sol1(input Almanac) int64 {
	res := int64(-1)
	for _, seed := range input.Seeds {
		loc := ApplyMaps(seed, input.Maps)
		if res == -1 {
			res = loc
			continue
		}
		if res > loc {
			res = loc
		}
	}
	return res
}

func ApplyMaps(val int64, maps []Map) int64 {
	res := val
	for _, m := range maps {
		res = m.Apply(res)
	}

	return res
}

func (r MapRange) Apply(k int64) (int64, bool) {
	var res int64
	if k < r.SourceStart || k > r.SourceEnd {
		return res, false
	}
	res = k + r.Offset
	return res, true
}

func (m Map) Apply(k int64) int64 {
	for _, r := range m.Ranges {
		val, success := r.Apply(k)
		if success {
			return val
		}
	}
	return k
}

func Sol2(input Almanac) int64 {
	min := int64(math.MaxInt64)
	cur := input.SeedRanges
	for _, m := range input.Maps {
		cur = m.ApplyRanges(cur)
	}

	for _, r := range cur {
		if r.Start < min {
			min = r.Start
		}
	}

	return min
}

func (m Map) ApplyRanges(ranges []Range) []Range {
	var res []Range
	var rem []Range
	rem = append(rem, ranges...)
	for _, mr := range m.Ranges {
		nr, newRem := mr.ApplyRanges(rem)
		res = append(res, nr...)
		rem = newRem
	}

	res = append(res, rem...)

	return res
}

func (r MapRange) ApplyRanges(ks []Range) ([]Range, []Range) {
	res := []Range{}
	rem := []Range{}
	for _, k := range ks {
		newR, newRem, success := r.ApplyRange(k)
		rem = append(rem, newRem...)
		if success {
			res = append(res, newR)
		}

	}

	return res, rem
}

func (r MapRange) ApplyRange(k Range) (Range, []Range, bool) {
	res := Range{}
	rem := k.Substract(r)

	if k.End < r.SourceStart || k.Start > r.SourceEnd {
		return res, rem, false
	}
	res.Start = Max(r.SourceStart, k.Start) + r.Offset
	res.End = Min(r.SourceEnd, k.End) + r.Offset

	return res, rem, true
}

func (r Range) Substract(mr MapRange) []Range {
	if mr.SourceEnd < r.Start || r.End < mr.SourceStart {
		return []Range{r}
	}

	if r.Start < mr.SourceStart && mr.SourceStart < r.End && r.End < mr.SourceEnd {
		rn := Range{
			Start: r.Start,
			End:   mr.SourceStart - 1,
		}
		return []Range{rn}
	}

	if mr.SourceStart < r.Start && r.Start < mr.SourceEnd && mr.SourceEnd < r.End {
		rn := Range{
			Start: mr.SourceEnd + 1,
			End:   r.End,
		}
		return []Range{rn}
	}

	if r.Start < mr.SourceStart && mr.SourceStart < mr.SourceEnd && mr.SourceEnd < r.End {
		r1 := Range{
			Start: r.Start,
			End:   mr.SourceStart - 1,
		}
		r2 := Range{
			Start: mr.SourceEnd + 1,
			End:   r.End,
		}

		return []Range{r1, r2}
	}

	return []Range{}
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
