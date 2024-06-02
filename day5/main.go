package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TokenType int32

const (
	STRING TokenType = iota
	NUMBER TokenType = iota
	MAP    TokenType = iota
)

type Token struct {
	tokenType TokenType
	value     string
}

func isString(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func getString(data []byte, i *int) string {
	start := *i
	for ; isString(data[*i]); *i++ {
	}
	return string(data[start:*i])
}

func getNumber(data []byte, i *int) string {
	start := *i
	for ; isDigit(data[*i]); *i++ {
	}
	return string(data[start:*i])
}

type FarmerMap struct {
	from     string
	to       string
	mappings []Mapping
}

type Mapping struct {
	min   int
	max   int
	delta int
}

func (m FarmerMap) String() string {
	s := fmt.Sprintf("%s to %s\n", m.from, m.to)
	sb := strings.Builder{}
	sb.WriteString(s)
	for _, v := range m.mappings {
		sb.WriteString(fmt.Sprintf("Diff: %d; Range: %d - %d\n", v.delta, v.min, v.max))
	}
	return sb.String()
}

type EntityRange struct {
	min int
	max int
}

func (m Mapping) inRange(i int) bool {
	return i >= m.min && i <= m.max
}

func getNumbers(tokens []Token, i *int) []int {
	numbers := make([]int, 0)
	for ; *i < len(tokens) && tokens[*i].tokenType == NUMBER; *i++ {
		token := tokens[*i]
		i, e := strconv.Atoi(token.value)
		if e != nil {
			fmt.Println("SOMETHING WENT WRONG")
		}
		numbers = append(numbers, i)
	}
	return numbers
}

func getSeedRanges(seeds []int) []EntityRange {
	ranges := make([]EntityRange, 0)
	for i := 0; i < len(seeds); i += 2 {
		ranges = append(ranges, EntityRange{min: seeds[i], max: seeds[i] + seeds[i+1] - 1})
	}
	return ranges
}

func main() {
	data, err := os.ReadFile("./example.txt")
	if err != nil {
		return
	}
	tokens := make([]Token, 0)
	i := 0
	for ; i < len(data); i++ {
		v := data[i]
		if isString(v) {
			s := getString(data, &i)
			if s == "map" {
				tokens = append(tokens, Token{tokenType: MAP, value: s})
			} else {
				tokens = append(tokens, Token{tokenType: STRING, value: s})
			}
		}
		if isDigit(v) {
			s := getNumber(data, &i)
			tokens = append(tokens, Token{tokenType: NUMBER, value: s})
		}
	}
	i = 0
	initialseeds := make([]EntityRange, 0)
	maps := make([]FarmerMap, 0)
	for ; i < len(tokens); i++ {
		token := tokens[i]
		if token.tokenType == STRING && token.value == "seeds" {
			i++ // next token
			seeds := getNumbers(tokens, &i)
			initialseeds = getSeedRanges(seeds)
			fmt.Println(initialseeds)
			token = tokens[i] // the value of `i` has changed
		}
		if token.tokenType == MAP {
			to := tokens[i-1].value
			from := tokens[i-3].value
			i++ // next token
			mapNumbers := getNumbers(tokens, &i)
			mappings := make([]Mapping, 0)
			for i := 0; i < len(mapNumbers); i += 3 {
				min := mapNumbers[i+1]
				max := mapNumbers[i+1] + mapNumbers[i+2] - 1
				delta := mapNumbers[i] - mapNumbers[i+1]
				mappings = append(mappings, Mapping{min: min, max: max, delta: delta})
			}
			farmerMap := FarmerMap{from: from, to: to, mappings: mappings}
			fmt.Println(farmerMap)
			maps = append(maps, farmerMap)
		}
	}
	min := -1
	for _, s := range initialseeds {
		fmt.Println(s)
		sts := processRange(maps[0], []EntityRange{s})
		fmt.Println(sts)
		stf := processRange(maps[1], sts)
		fmt.Println(stf)
		ftw := processRange(maps[2], stf)
		fmt.Println(ftw)
		wtl := processRange(maps[3], ftw)
		fmt.Println(wtl)
		ltt := processRange(maps[4], wtl)
		fmt.Println(ltt)
		tth := processRange(maps[5], ltt)
		fmt.Println(tth)
		htl := processRange(maps[6], tth)
		fmt.Println(htl)
		for _, v := range htl {
			if min == -1 {
				min = v.min
			}
			if v.min < min {
				min = v.min
			}
			if v.max < min {
				min = v.max
			}
		}
	}
	println(min)
}

func processRange(fmap FarmerMap, ranges []EntityRange) []EntityRange {
	entityRanges := make([]EntityRange, 0)
	for i := 0; i < len(ranges); i++ {
		r := &ranges[i]
		minIndex := -1
		maxIndex := -1
		for i, v := range fmap.mappings {
			if v.inRange(r.min) {
				minIndex = i
			}
			if v.inRange(r.max) {
				maxIndex = i
			}
		}
		if minIndex == maxIndex {
			if minIndex == -1 {
				entityRanges = append(entityRanges, *r)
			} else {
				r.min = r.min + fmap.mappings[minIndex].delta
				r.max = r.max + fmap.mappings[minIndex].delta
				entityRanges = append(entityRanges, *r)
			}
		} else {
			if minIndex > -1 && maxIndex > -1 {
				minr := EntityRange{min: r.min, max: fmap.mappings[minIndex].max}
				maxr := EntityRange{min: fmap.mappings[maxIndex].min, max: r.max}
				entityRanges = append(entityRanges, processRange(fmap, []EntityRange{minr})...)
				entityRanges = append(entityRanges, processRange(fmap, []EntityRange{maxr})...)
			}
			if minIndex > -1 && maxIndex == -1 {
				minr := EntityRange{min: r.min, max: fmap.mappings[minIndex].max}
				maxr := EntityRange{min: fmap.mappings[minIndex].max + 1, max: r.max}
				entityRanges = append(entityRanges, processRange(fmap, []EntityRange{minr})...)
				entityRanges = append(entityRanges, processRange(fmap, []EntityRange{maxr})...)
			}
			if maxIndex > -1 && minIndex == -1 {
				minr := EntityRange{min: r.min, max: fmap.mappings[maxIndex].min - 1}
				maxr := EntityRange{min: fmap.mappings[maxIndex].min, max: r.max}
				entityRanges = append(entityRanges, processRange(fmap, []EntityRange{minr})...)
				entityRanges = append(entityRanges, processRange(fmap, []EntityRange{maxr})...)
			}
		}
	}
	return entityRanges
}
