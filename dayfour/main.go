package main

import (
	"fmt"
	"math"
	"os"
)

type Card struct {
	id      int
	winning map[int]bool
	numbers []int
	winners int
	points  int
}

type LiteCard struct {
	id      int
	winners int
}

func (card LiteCard) String() string {
	return fmt.Sprintf("{Id: %d, Winners: %d}", card.id, card.winners)
}

func parseDigit(data []byte) int {
	val := 0
	for i, v := range data {
		val += int(rune(v)-rune('0')) * int(math.Pow(10, float64(len(data)-i-1)))
	}
	return val
}

func parseData(data []byte) []Card {
	cards := make([]Card, 0)
	digitStart := -1
	var card *Card
	winners := true
	for i, v := range data {
		if card == nil {
			card = &Card{}
			card.winning = make(map[int]bool)
		}
		if v >= '0' && v <= '9' {
			if digitStart == -1 {
				digitStart = i
			}
		}
		if v == ' ' {
			if digitStart > -1 {
				number := parseDigit(data[digitStart:i])
				digitStart = -1
				if winners {
					card.winning[number] = true
				} else {
					card.numbers = append(card.numbers, number)
				}
			}
		}
		if v == ':' {
			if digitStart > -1 {
				card.id = parseDigit(data[digitStart:i])
				digitStart = -1
			}
		}
		if v == '|' {
			winners = false
		}
		if v == '\n' {
			if digitStart > -1 {
				number := parseDigit(data[digitStart:i])
				digitStart = -1
				if winners {
					card.winning[number] = true
				} else {
					card.numbers = append(card.numbers, number)
				}
			}
			cards = append(cards, *card)
			card = nil
			winners = true
		}
	}
	return cards
}

func calculateWinners(cards []Card) int {
	sum := 0
	for i := 0; i < len(cards); i++ {
		c := &cards[i]
		for _, v := range c.numbers {
			_, ok := c.winning[v]
			if ok {
				c.winners++
			}
		}
		c.points = int(math.Pow(2, float64(c.winners-1)))
		sum += c.points
	}
	return sum
}

func processDeck(cards []LiteCard) int {
	total := 0
	queue := append(make([]LiteCard, 0, len(cards)), cards...)
	for len(queue) != 0 {
		card := queue[0]
		if card.winners > 0 {
			queue = append(queue, cards[card.id:card.id+card.winners]...)
		}
		total += 1
		queue = queue[1:]
	}
	return total
}

func main() {
	data, err := os.ReadFile("./data.txt")
	if err != nil {
		return
	}
	cards := parseData(data)
	sum := calculateWinners(cards)
	liteCards := make([]LiteCard, 0, len(cards))
	for _, v := range cards {
		liteCards = append(liteCards, LiteCard{id: v.id, winners: v.winners})
	}
	total := processDeck(liteCards)
	fmt.Println(sum, total)
}
