package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type CommandType int

const (
	Add CommandType = iota
	Delete
)

type Lens struct {
	Label      string
	FocalPower int
}

type Bukkit struct {
	Content []Lens
}

type HashMap struct {
	Boxes [256]Bukkit
}

type Command struct {
	Lens Lens
	Type CommandType
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

func ParseInput(input []string) []string {
	str := input[0]
	return strings.Split(str, ",")
}

func ParseAsCommand(s string) Command {
	res := Command{}
	isDelete := strings.Contains(s, "-")
	if isDelete {
		label, _ := strings.CutSuffix(s, "-")
		res.Type = Delete
		res.Lens.Label = label
		return res
	}
	splitted := strings.Split(s, "=")
	label := splitted[0]
	focalStr := splitted[1]
	focal, _ := strconv.Atoi(focalStr)

	res.Type = Add
	res.Lens.FocalPower = focal
	res.Lens.Label = label

	return res
}

func Sol1(input []string) int {
	res := 0
	for _, s := range input {
		res += HASH(s)
	}
	return res
}

func HASH(s string) int {
	curVal := 0
	for _, c := range s {
		ascii := int(c)
		curVal += ascii
		curVal *= 17
		curVal %= 256
	}
	return curVal
}

func HASHLens(l Lens) int {
	return HASH(l.Label)
}

func Sol2(input []string) int {
	res := 0
	hashMap := HashMap{}
	for _, s := range input {
		comm := ParseAsCommand(s)
		hashMap.Execute(comm)
	}

	for box, bukkit := range hashMap.Boxes {
		curBox := box + 1
		for slot, lens := range bukkit.Content {
			curSlot := slot + 1
			res += curSlot * curBox * lens.FocalPower
		}
	}

	return res
}

func (h *HashMap) Execute(c Command) {
	switch c.Type {
	case Add:
		h.Add(c.Lens)
	case Delete:
		h.Delete(c.Lens)
	}
}

func (h *HashMap) Add(l Lens) {
	idx := HASHLens(l)
	h.Boxes[idx] = h.Boxes[idx].Add(l)
}

func (h *HashMap) Delete(l Lens) {
	idx := HASHLens(l)
	h.Boxes[idx] = h.Boxes[idx].Delete(l)
}

func (b Bukkit) Add(l Lens) Bukkit {
	res := Bukkit{}
	lenses := slices.Clone(b.Content)
	idx := -1
	for i, lr := range lenses {
		if lr.Label == l.Label {
			idx = i
			break
		}
	}

	if idx == -1 {
		lenses = append(lenses, l)
	} else {
		lenses[idx] = l
	}

	res.Content = lenses
	return res
}

func (b Bukkit) Delete(l Lens) Bukkit {
	res := Bukkit{}
	lenses := slices.Clone(b.Content)
	idx := -1

	for i, lr := range lenses {
		if lr.Label == l.Label {
			idx = i
			break
		}
	}

	if idx != -1 {
		before := lenses[:idx]
		after := lenses[idx+1:]
		before = append(before, after...)
		lenses = before
	}
	res.Content = lenses
	return res
}
