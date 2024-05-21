package main

import (
	"fmt"
	"math"
	"os"
)

func parseDigit(data []byte) (int, int) {
	val := 0
	for i, v := range data {
		val += int(rune(v)-rune('0')) * int(math.Pow(10, float64(len(data)-i-1)))
	}
	return val, len(data)
}

func main() {
	data, err := os.ReadFile("./data.txt")
	if err != nil {
		return
	}
	plane := make([][]byte, 0)
	currentLine := 0
	for i, v := range data {
		if v == '\n' {
			plane = append(plane, data[currentLine:i])
			currentLine = i + 1
		}
	}
	sum := 0
	total := 0
	digitStart := -1
	for i := 0; i < len(plane); i++ {
		row := plane[i]
		for j := 0; j < len(row); j++ {
			v := row[j]
			if v >= '0' && v <= '9' && j+1 != len(row) {
				if digitStart == -1 {
					digitStart = j
				}
			} else {
				if digitStart > -1 {
					val, length := parseDigit(row[digitStart:j])
					var startingRow, endingRow, startingColumn, endingColumn int
					if i-1 > 0 {
						startingRow = i - 1
					} else {
						startingRow = 0
					}
					if i+1 < len(plane) {
						endingRow = i + 1
					} else {
						endingRow = i
					}
					if digitStart-1 >= 0 {
						startingColumn = digitStart - 1
					} else {
						startingColumn = digitStart
					}
					if digitStart+length < len(row) {
						endingColumn = digitStart + length
					} else {
						endingColumn = len(row) - 1
					}
					part := false
					symbol := '.'
					for k := startingRow; k <= endingRow; k++ {
						checkRow := plane[k]
						for l := startingColumn; l <= endingColumn; l++ {
							checkItem := checkRow[l]
							if checkItem >= 32 && checkItem <= 47 && checkItem != '.' {
								part = true
								symbol = rune(checkItem)
								break
							}
						}
						if part {
							break
						}
					}
					fmt.Println(startingRow+1, endingRow+1, startingColumn+1, endingColumn+1, i+1, j+1, val, part, (string)(symbol))
					if part {
						sum += val
					}
					total += val
					digitStart = -1
				}
			}
		}
	}
	fmt.Println(sum, total)
}
