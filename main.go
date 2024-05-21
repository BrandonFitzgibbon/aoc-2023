package main

import (
	"fmt"
	"math"
	"os"
)

const (
	BLUE  = 14
	RED   = 12
	GREEN = 13
)

type Reveal struct {
	blue  int
	red   int
	green int
}

func (reveal Reveal) String() string {
	return fmt.Sprintf("\nBlue: %d Red: %d Green: %d", reveal.blue, reveal.red, reveal.green)
}

func (reveal Reveal) Power() int {
	return reveal.blue * reveal.red * reveal.green
}

type Game struct {
	id      int
	reveals []Reveal
	power   int
}

func (game Game) String() string {
	return fmt.Sprintf("\nID: %d Reveals: %v", game.id, game.reveals)
}

func parseDigit(data []rune) int {
	val := 0
	for i, v := range data {
		val += int(v-rune('0')) * int(math.Pow(10, float64(len(data)-i-1)))
	}
	return val
}

func parseReveal(data []byte, minReveal *Reveal) (Reveal, bool) {
	var blue, red, green int
	digits := make([]rune, 0)
	validReveal := true
	for i := 0; i < len(data); i++ {
		v := data[i]
		if v >= '0' && v <= '9' {
			digits = append(digits, rune(v))
		}
		if v == ' ' && len(digits) > 0 {
			val := parseDigit(digits)
			digits = digits[:0]
			i++
			switch data[i] {
			case 'b':
				blue = val
				if minReveal.blue < blue {
					minReveal.blue = blue
				}
				i += 3
			case 'r':
				red = val
				if minReveal.red < red {
					minReveal.red = red
				}
				i += 2
			case 'g':
				green = val
				if minReveal.green < green {
					minReveal.green = green
				}
				i += 4
			}
		}
	}
	if blue > BLUE || red > RED || green > GREEN {
		validReveal = true // toggling true for part 2
	}
	return Reveal{blue: blue, red: red, green: green}, validReveal
}

func parseGame(data []byte) (Game, bool) {
	var id int
	cursor := 0
	validGame := true
	if data[0] == 'G' && data[1] == 'a' && data[2] == 'm' && data[3] == 'e' && data[4] == ' ' {
		digits := make([]rune, 0, 3)
		cursor = 5
		for ; data[cursor] != ':'; cursor++ {
			if data[cursor] >= '0' && data[cursor] <= '9' {
				digits = append(digits, rune(data[cursor]))
			}
		}
		id = parseDigit(digits)
		game := Game{id: id, reveals: make([]Reveal, 0)}
		minReveal := Reveal{}
		if data[cursor] == ':' && data[cursor+1] == ' ' {
			cursor += 2
			lastBlock := cursor
			for ; cursor < len(data); cursor++ {
				if data[cursor] == ';' || cursor == len(data)-1 {
					reveal, validReveal := parseReveal(data[lastBlock:cursor], &minReveal)
					if !validReveal {
						validGame = false
					}
					game.reveals = append(game.reveals, reveal)
					cursor += 2
					lastBlock = cursor
				}
			}
		}
		game.power = minReveal.Power()
		return game, validGame
	}
	return Game{id: id}, false
}

func parseGames(data []byte) int {
	lastLine := 0
	sum := 0
	for i := 0; i < len(data); i++ {
		v := data[i]
		if v == 10 {
			game, validGame := parseGame(data[lastLine:i])
			if validGame {
				sum += game.power
			}
			lastLine = i + 1
		}
	}
	return sum
}

func main() {
	data, err := os.ReadFile("./data.txt")
	if err != nil {
		return
	}
	sum := parseGames(data)
	fmt.Println(sum)
}
