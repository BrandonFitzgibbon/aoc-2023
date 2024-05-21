package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("data.txt")
	if err != nil {
		return
	}
	digitMap := map[string]rune{
		"one":   '1',
		"two":   '2',
		"three": '3',
		"four":  '4',
		"five":  '5',
		"six":   '6',
		"seven": '7',
		"eight": '8',
		"nine":  '9',
	}
	magicNumbers := make([]int32, 0)
	runer := make([]rune, 0)
	var firstDigit, lastDigit, zeroRune, sum rune
	firstDigit = 0
	lastDigit = 0
	zeroRune = '0'
	line := 1
	for _, v := range data {
		if v == 10 {
			number := ((firstDigit - zeroRune) * 10) + (lastDigit - zeroRune)
			fmt.Printf("%d %d\n", line, number)
			line++
			magicNumbers = append(magicNumbers, number)
			firstDigit = 0
			lastDigit = 0
			runer = runer[:0]
			continue
		}
		if v >= 48 && v <= 57 {
			if firstDigit == 0 {
				firstDigit = rune(v)
			}
			lastDigit = rune(v)
			runer = runer[:0]
			continue
		}
		runer = append(runer, rune(v))
		if len(runer) > 2 {
			for k, v := range digitMap {
				if strings.Contains(string(runer), k) {
					if firstDigit == 0 {
						firstDigit = v
					}
					lastDigit = v
					runer = runer[len(runer)-1:]
					break
				}
			}
		}
	}
	sum = 0
	for _, v := range magicNumbers {
		sum += v
	}
	fmt.Println(sum)
	return
}
