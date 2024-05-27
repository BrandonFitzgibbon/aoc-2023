package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	WIDTH = 141.0
)

type Container struct {
	number int
	nlen   int
	data   []byte
	part   bool
	line   int
	col    int
}

func (container Container) String() string {
	sb := strings.Builder{}
	span := container.nlen
	if container.col == 1 || container.col+container.nlen == WIDTH {
		span += 1
	} else {
		span += 2
	}
	sb.WriteString(fmt.Sprintf("\nNumber %d\nLength %d\nLine %d\nCol %d\nPart %t\nData:\n", container.number, container.nlen, container.line, container.col, container.part))
	for i, v := range container.data {
		sb.WriteByte(v)
		if (i+1)%span == 0 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func parseDigit(data []byte) int {
	val := 0
	for i, v := range data {
		val += int(rune(v)-rune('0')) * int(math.Pow(10, float64(len(data)-i-1)))
	}
	return val
}

func getDataSlice(data []byte, row int, col int, len int) ([]byte, bool) {
	part := false
	slice := make([]byte, 0)
	for i := -1; i < 2; i++ {
		if row+i < 1 {
			continue
		}
		if row+i > WIDTH-1 {
			continue
		}
		for j := -1; j < len+1; j++ {
			if col+j < 1 {
				continue
			}
			if col+j > WIDTH-1 {
				continue
			}
			ind := ((row - 1 + i) * WIDTH) + (col - 1 + j)
			ele := data[ind]
			if ele != '.' && !(ele >= '0' && ele <= '9') {
				part = true
			}
			slice = append(slice, ele)
		}
	}
	return slice, part
}

func main() {
	data, err := os.ReadFile("./data.txt")
	if err != nil {
		return
	}
	sb := strings.Builder{}
	containers := make([]Container, 0)
	digitStart := -1
	for i, v := range data {
		if v >= '0' && v <= '9' {
			if digitStart == -1 {
				digitStart = i
			}
		} else {
			if digitStart > -1 {
				number := parseDigit(data[digitStart:i])
				len := i - digitStart
				row := int(math.Ceil(float64(digitStart+1) / WIDTH))
				col := (WIDTH) - ((row * WIDTH) - digitStart) + 1
				container := Container{number: number, nlen: len, line: row, col: col}
				container.data, container.part = getDataSlice(data, row, col, len)
				containers = append(containers, container)
				digitStart = -1
			}
		}
	}
	sum := 0
	for _, v := range containers {
		sb.WriteString(fmt.Sprint(v))
		if v.part {
			sum += v.number
		}
	}
	fmt.Println(sum)
	os.WriteFile("./output.txt", []byte(sb.String()), 0777)
}
