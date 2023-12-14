package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Record int

const (
	RecordOk Record = iota
	RecordDamaged
	RecordUnkonw
)

type Row struct {
	Groups  []int
	Records []Record
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

func ParseInput(input []string) []Row {
	var res []Row
	for _, s := range input {
		newRow := Row{}

		splited := strings.Split(s, " ")
		recordsStr := splited[0]
		groupsStr := splited[1]

		records := ParseRecords(recordsStr)
		groups := ParseGroups(groupsStr)

		newRow.Records = records
		newRow.Groups = groups

		res = append(res, newRow)
	}
	return res
}

func ParseRecords(input string) []Record {
	var res []Record
	for _, c := range input {
		switch c {
		case '.':
			res = append(res, RecordOk)
		case '#':
			res = append(res, RecordDamaged)
		case '?':
			res = append(res, RecordUnkonw)
		}
	}
	return res
}

func ParseGroups(input string) []int {
	var res []int
	splitted := strings.Split(input, ",")
	for _, s := range splitted {
		val, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		res = append(res, val)
	}
	return res
}

// bruteforce
func Sol1(input []Row) int {
	res := 0
	for _, r := range input {
		res += CalcPossibleArrangementsByGroups(r)
	}
	return res
}

// cheated...
func Sol2(input []Row) int {
	res := 0
	for _, r := range input {
		gruopsClone := slices.Clone(r.Groups)
		recordsClone := slices.Clone(r.Records)
		for i := 1; i <= 4; i++ {
			r.Records = append(r.Records, RecordUnkonw)
			r.Records = append(r.Records, recordsClone...)
			r.Groups = append(r.Groups, gruopsClone...)
		}
		ResetMemo()
		res += CalcArrangementsDyn(r)
	}
	return res
}

func CalcPossibleArrangementsByGroups(r Row) int {
	mask := r.Records
	groups := r.Groups
	rowLen := len(mask)
	res := 0

	maxOffset := rowLen
	for _, v := range groups {
		maxOffset -= v
	}
	okRecordsRequired := len(groups) - 1
	maxOffset -= okRecordsRequired
	base := maxOffset + 1

	isAvailable := true
	curOffsets := make([]int, len(groups))
	for isAvailable {
		rec := ConstructRecsByGroupsAndOffsets(groups, curOffsets, rowLen)

		if IsValidForMask(mask, rec) {
			res++
		}
		isAvailable = IncrOffsetWithBase(curOffsets, base)
	}

	return res
}

func ConstructRecsByGroupsAndOffsets(groups []int, offsets []int, rowLen int) []Record {
	var rec []Record
	for i, v := range offsets {
		damaged := groups[i]
		for c := 0; c < v; c++ {
			rec = append(rec, RecordOk)
		}
		if i != 0 {
			rec = append(rec, RecordOk)
		}

		for c := 0; c < damaged; c++ {
			rec = append(rec, RecordDamaged)
		}
	}
	tail := rowLen - len(rec)
	for tail > 0 {
		rec = append(rec, RecordOk)
		tail--
	}
	return rec
}

func IncrOffsetWithBase(offsets []int, base int) bool {
	if len(offsets) <= 0 {
		return true
	}
	offsets[0]++
	if len(offsets) == 1 && offsets[0] >= base {
		return false
	}
	if Sum(offsets) >= base {
		offsets[0] = 0
		return IncrOffsetWithBase(offsets[1:], base)
	}
	return true
}

func IsValidForMask(mask []Record, r []Record) bool {
	isValid := true
	for i, v := range mask {
		if v == RecordUnkonw {
			continue
		}
		if v != r[i] {
			isValid = false
			break
		}
	}
	return isValid
}

func Sum(arr []int) int {
	res := 0
	for _, v := range arr {
		res += v
	}
	return res
}

func CalcArrangementsDyn(r Row) int {
	gLen := len(r.Groups)
	rLen := len(r.Records)
	rStr := RecordsToStr(r.Records)
	v, isCached := GetMemo(gLen, rStr)
	if isCached {
		return v
	}
	if gLen == 0 && rLen == 0 {
		return 1
	}
	if rLen == 0 {
		return 0
	}
	curRec := r.Records[0]
	curGroup := 0
	if gLen != 0 {
		curGroup = r.Groups[0]
	}
	switch curRec {
	case RecordOk:
		newRow := Row{
			Records: r.Records[1:],
			Groups:  r.Groups,
		}
		res := CalcArrangementsDyn(newRow)
		SaveMemo(gLen, rStr, res)
		return res
	case RecordUnkonw:
		newRecDam := slices.Clone(r.Records)
		newRecDam[0] = RecordDamaged
		newRowDam := Row{
			Records: newRecDam,
			Groups:  r.Groups,
		}
		newRecOk := slices.Clone(r.Records)
		newRecOk[0] = RecordOk
		newRowOk := Row{
			Records: newRecOk,
			Groups:  r.Groups,
		}
		res := CalcArrangementsDyn(newRowDam) + CalcArrangementsDyn(newRowOk)
		SaveMemo(gLen, rStr, res)
		return res
	case RecordDamaged:
		if CanMatch(r.Records, curGroup) {
			newIdx := curGroup + 1
			if curGroup == len(r.Records) {
				newIdx = curGroup
			}
			newRow := Row{
				Groups:  r.Groups[1:],
				Records: r.Records[newIdx:],
			}
			res := CalcArrangementsDyn(newRow)
			SaveMemo(gLen, rStr, res)
			return res
		} else {
			SaveMemo(gLen, rStr, 0)
			return 0
		}
	}
	return 0
}

func CanMatch(r []Record, group int) bool {
	if len(r) < group {
		return false
	}
	isEndValid := len(r) == group || r[group] == RecordOk || r[group] == RecordUnkonw
	isIntervalValid := true
	for i := 0; i < group; i++ {
		isIntervalValid = isIntervalValid && r[i] != RecordOk
	}

	return isEndValid && isIntervalValid
}

var cache = make(map[int]map[string]int)

func ResetMemo() {
	cache = make(map[int]map[string]int)
}

func SaveMemo(g int, r string, v int) {
	_, isSubMapExist := cache[g]
	if !isSubMapExist {
		cache[g] = make(map[string]int)
	}
	cache[g][r] = v
}

func GetMemo(g int, r string) (int, bool) {
	_, isSubMapExists := cache[g]
	if !isSubMapExists {
		return 0, false
	}

	v, b := cache[g][r]
	return v, b
}

func RecordsToStr(rs []Record) string {
	res := ""
	for _, r := range rs {
		switch r {
		case RecordDamaged:
			res += "#"
		case RecordOk:
			res += "."
		case RecordUnkonw:
			res += "?"
		}
	}
	return res
}
