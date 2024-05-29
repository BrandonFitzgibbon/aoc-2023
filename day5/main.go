package main

import (
	"fmt"
	"os"
	"strconv"
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
	mappings []int
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
	seeds := make([]int, 0)
	maps := make([]FarmerMap, 0)
	for ; i < len(tokens); i++ {
		token := tokens[i]
		if token.tokenType == STRING && token.value == "seeds" {
			i++ // next token
			seeds = getNumbers(tokens, &i)
			fmt.Println(seeds)
			token = tokens[i] // the value of `i` has changed
		}
		if token.tokenType == MAP {
			to := tokens[i-1].value
			from := tokens[i-3].value
			i++ // next token
			mappings := getNumbers(tokens, &i)
			farmerMap := FarmerMap{from: from, to: to, mappings: mappings}
			fmt.Println(farmerMap)
			maps = append(maps, farmerMap)
		}
	}
	for _, v := range seeds {
		for _, m := range maps {
			for i, k := range m.mappings {

			}
		}
	}
}
