package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cell interface {
	Next(Beam) []Beam
}

type EmtpyCell struct {
}

func (c EmtpyCell) Next(b Beam) []Beam {
	res := Beam{
		Vel:    b.Vel,
		CurPos: b.CurPos.Add(b.Vel),
	}
	return []Beam{res}
}

type Mirror struct {
	Normal Vec2
}

func (c Mirror) Next(b Beam) []Beam {
	var res []Beam
	newVel := b.Vel.MirrorOver(c.Normal)
	bRes := Beam{
		Vel:    newVel,
		CurPos: b.CurPos.Add(newVel),
	}
	res = append(res, bRes)
	return res
}

type Splitter struct {
	Normal Vec2
}

func (c Splitter) Next(b Beam) []Beam {
	var res []Beam
	dotProduct := b.Vel.DotProduct(c.Normal)
	if dotProduct == 0 {
		bRes := Beam{
			Vel:    b.Vel,
			CurPos: b.CurPos.Add(b.Vel),
		}
		res = append(res, bRes)
		return res
	}

	vel1 := b.Vel.Rot90CW()
	vel2 := b.Vel.Rot90CCW()
	b1 := Beam{
		Vel:    vel1,
		CurPos: b.CurPos.Add(vel1),
	}
	b2 := Beam{
		Vel:    vel2,
		CurPos: b.CurPos.Add(vel2),
	}
	res = append(res, b1, b2)

	return res
}

type Vec2 struct {
	X, Y int
}

type Beam struct {
	CurPos, Vel Vec2
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	res := Vec2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
	return res
}

func (v1 Vec2) DotProduct(v2 Vec2) int {
	res := v1.X*v2.X + v1.Y*v2.Y
	return res
}

func (v Vec2) Rot90CW() Vec2 {
	res := Vec2{
		X: v.Y,
		Y: -v.X,
	}
	return res
}

func (v Vec2) Rot90CCW() Vec2 {
	res := Vec2{
		X: -v.Y,
		Y: v.X,
	}
	return res
}

func (v1 Vec2) MirrorOver(v2 Vec2) Vec2 {
	k := 2 * v1.DotProduct(v2) / v2.DotProduct(v2)
	negV1 := v1.Mult(-1)
	multV2 := v2.Mult(k)
	vRefl := negV1.Add(multV2)
	return vRefl
}

func (v Vec2) Mult(k int) Vec2 {
	res := Vec2{
		X: v.X * k,
		Y: v.Y * k,
	}
	return res
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

func ParseInput(input []string) [][]Cell {
	var res [][]Cell
	for _, s := range input {
		var curLine []Cell
		for _, c := range s {
			var curCell Cell
			switch c {
			case '.':
				curCell = EmtpyCell{}
			case '\\':
				fallthrough
			case '/':
				curCell = ParseMirror(c)
			case '-':
				fallthrough
			case '|':
				curCell = ParseSplitter(c)
			}
			curLine = append(curLine, curCell)
		}
		res = append(res, curLine)
	}
	return res
}

func ParseMirror(r rune) Mirror {
	res := Mirror{}
	normal := Vec2{}
	switch r {
	case '\\':
		normal = Vec2{X: 1, Y: 1}
	case '/':
		normal = Vec2{X: -1, Y: 1}
	}
	res.Normal = normal
	return res
}

func ParseSplitter(r rune) Splitter {
	res := Splitter{}
	normal := Vec2{}
	switch r {
	case '|':
		normal = Vec2{X: 1, Y: 0}
	case '-':
		normal = Vec2{X: 0, Y: 1}
	}
	res.Normal = normal
	return res
}

func Sol1(cells [][]Cell) int {
	zeroBeam := Beam{
		CurPos: Vec2{X: 0, Y: 0},
		Vel:    Vec2{X: 1, Y: 0},
	}
	res := CalcEnergyLevel(cells, zeroBeam)
	return res
}

func Sol2(cells [][]Cell) int {
	res := 0

	maxX := len(cells[0]) - 1
	maxY := len(cells) - 1

	entryVel := Vec2{X: 0, Y: 1}
	for i := range cells {
		entryPos := Vec2{X: i, Y: 0}
		entryBeam := Beam{Vel: entryVel, CurPos: entryPos}
		energy := CalcEnergyLevel(cells, entryBeam)
		if energy > res {
			res = energy
		}
	}

	entryVel = Vec2{X: 1, Y: 0}
	for i := range cells {
		entryPos := Vec2{X: 0, Y: i}
		entryBeam := Beam{Vel: entryVel, CurPos: entryPos}
		energy := CalcEnergyLevel(cells, entryBeam)
		if energy > res {
			res = energy
		}
	}

	entryVel = Vec2{X: 0, Y: -1}
	for i := range cells {
		entryPos := Vec2{X: i, Y: maxY}
		entryBeam := Beam{Vel: entryVel, CurPos: entryPos}
		energy := CalcEnergyLevel(cells, entryBeam)
		if energy > res {
			res = energy
		}
	}

	entryVel = Vec2{X: -1, Y: 0}
	for i := range cells {
		entryPos := Vec2{X: maxX, Y: i}
		entryBeam := Beam{Vel: entryVel, CurPos: entryPos}
		energy := CalcEnergyLevel(cells, entryBeam)
		if energy > res {
			res = energy
		}
	}

	return res
}

func CalcEnergyLevel(cells [][]Cell, enrtyBeam Beam) int {
	res := 0
	curState := []Beam{enrtyBeam}
	maxX := len(cells[0]) - 1
	maxY := len(cells) - 1

	history := make(map[Beam]bool)
	isGoing := true
	for isGoing {
		var newBeams []Beam

		isGoing = false
		for _, b := range curState {
			if history[b] {
				continue
			}
			x := b.CurPos.X
			y := b.CurPos.Y
			if x > maxX || x < 0 || y > maxY || y < 0 {
				continue
			}
			isGoing = true
			history[b] = true
			bs := cells[y][x].Next(b)
			newBeams = append(newBeams, bs...)
		}

		curState = newBeams
	}

	visited := make(map[Vec2]bool)
	for k := range history {
		visited[k.CurPos] = true
	}
	for _, v := range visited {
		if v {
			res++
		}
	}

	return res

}
