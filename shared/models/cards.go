package models

import (
	"fmt"
	"gtihub.com/hariolate/tonneau/shared"
	"math/rand"
	"time"
)

type Color int

const (
	ColorAny Color = iota - 1
	ColorBlue
	ColorGreen
	ColorRed
	ColorYellow
)

type Number int

const NumberNone Number = -1

type Action int

const (
	ActionNone Action = iota - 1
	ActionReverse
	ActionSkip
	ActionDrawTwo
	ActionWild
	ActionWildDrawFour
)

type Card struct {
	Color  Color  `json:"color"`
	Number Number `json:"number"`
	Action Action `json:"action"`
}

func (c Card) Score() uint {
	switch c.Action {
	case ActionDrawTwo, ActionReverse, ActionSkip:
		return 20
	case ActionWild, ActionWildDrawFour:
		return 50
	}
	return uint(c.Number)
}

var WildCard = Card{
	ColorAny, NumberNone, ActionWild,
}

var WildDrawFourCard = Card{
	ColorAny, NumberNone, ActionWildDrawFour,
}

func MakeDrawTwoCard(c Color) Card {
	return Card{
		c, NumberNone, ActionDrawTwo,
	}
}

func MakeSkipCard(c Color) Card {
	return Card{
		c, NumberNone, ActionSkip,
	}
}

func MakeReverseCard(c Color) Card {
	return Card{
		c, NumberNone, ActionReverse,
	}
}

func MakeNumberCard(c Color, n Number) Card {
	if n < 0 {
		panic("invalid number card")
	}
	return Card{c, n, ActionNone}
}

const cardRedisValueTemplate = "%d:%d:%d" // number:color:action

func (c Card) redisValue() string {
	return fmt.Sprintf(cardRedisValueTemplate, c.Number, c.Color, c.Action)
}

func makeCardFromRedisValue(s string) Card {
	var c Card
	_, err := fmt.Scanf(cardRedisValueTemplate, &c.Number, &c.Color, &c.Action)
	shared.NoError(err)
	return c
}

func GenCards() []Card {
	var cards []Card

	for color := ColorBlue; color <= ColorYellow; color++ {
		for number := Number(0); number <= Number(9); number++ {
			card := MakeNumberCard(color, number)
			cards = append(cards, card)
		}
		drawTwoCard := MakeDrawTwoCard(color)
		skipCard := MakeSkipCard(color)
		reverseCard := MakeReverseCard(color)
		cards = append(cards, drawTwoCard, skipCard, reverseCard)
	}

	for i := 0; i < 3; i++ {
		cards = append(cards, WildCard, WildDrawFourCard)
	}

	if len(cards) != 54 {
		panic("invalid card set")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return cards
}

var InvalidCard = Card{-1, -1, -1}
