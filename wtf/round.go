// +build ignore

package wtf

import (
	"context"
	"gtihub.com/hariolate/tonneau/shared/models"
	"math/rand"
	"time"
)
import "github.com/google/uuid"

type Direction bool

const (
	Clockwise     Direction = true
	Anticlockwise           = false
)

type Round struct {
	discardCards []models.Card

	id uuid.UUID

	c context.Context

	draw int
	dir  Direction

	lastColor  models.Color
	lastNumber models.Number

	currentPlayer int
	players       []*Player
}

func NewRound(players []*Player) *Round {
	if len(players) < 2 || len(players) > 7 {
		panic("invalid player number")
	}

	round := &Round{
		id: uuid.New(),

		draw: 1,
		dir:  Clockwise,

		players:       players,
		currentPlayer: 0,

		discardCards: models.GenCards(),
		lastColor:    models.ColorAny,
		lastNumber:   models.NumberNone,
	}

	round.setup()

	return round
}

func (r *Round) setup() {
	for i := 0; i < len(r.players); i++ {
		for j := 0; j < 7; j++ {
			r.players[i].cards = append(r.players[i].cards, r.getCard())
		}
	}
}

func (r *Round) nextPlayer() {
	if r.dir == Clockwise {
		r.currentPlayer++
		if r.currentPlayer == len(r.players) {
			r.currentPlayer = 0
		}
	} else {
		r.currentPlayer--
		if r.currentPlayer == -1 {
			r.currentPlayer = len(r.players) - 1
		}
	}
}

func (r *Round) getCard() models.Card {
	card := r.discardCards[0]
	r.discardCards = r.discardCards[1:]
	return card
}

func (r *Round) putNumberCard(c models.Card) {
	if c.Action != models.ActionNone {
		panic("not a number card")
	}
	if r.lastColor != models.ColorAny || r.lastNumber != models.NumberNone {
		if c.Color != r.lastColor && c.Number != r.lastNumber {
			panic("not a valid number card")
		}
	}

	r.lastNumber = c.Number
	r.lastColor = c.Color
	r.putCardBack(c)
}

func (r *Round) putWildCard(c models.Card) {
	if c.Action != models.ActionWild && c.Action != models.ActionWildDrawFour {
		panic("not a wild card")
	}
	r.lastColor = r.players[r.currentPlayer].onChooseColor()
	r.lastNumber = models.NumberNone
	if c.Action == models.ActionWildDrawFour {
		if r.draw == 1 {
			r.draw = 4
		} else {
			r.draw += 4
		}
	}
	r.putCardBack(c)
}

func (r *Round) putActionCard(c models.Card) {
	if r.lastColor == models.ColorAny {
		r.lastColor = c.Color
		r.lastNumber = models.NumberNone
	} else if r.lastColor != c.Color {
		panic("invalid action card")
	}

	switch c.Action {
	case models.ActionReverse:
		r.dir = !r.dir
	case models.ActionSkip:
		r.nextPlayer()
	case models.ActionDrawTwo:
		if r.draw == 1 {
			r.draw = 2
		} else {
			r.draw += 2
		}
	default:
		panic("invalid action card")
	}

	r.putCardBack(c)
}

func (r *Round) shuffleDiscardCards() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(r.discardCards), func(i, j int) { r.discardCards[i], r.discardCards[j] = r.discardCards[j], r.discardCards[i] })
}

func (r *Round) getCurrentPlayer() *Player {
	return r.players[r.currentPlayer]
}

func (r *Round) putCardBack(c models.Card) {
	r.discardCards = append(r.discardCards, c)
	r.shuffleDiscardCards()
	r.nextPlayer()
	r.getCurrentPlayer().onTurn()
}

func (r *Round) putCard(c models.Card) {
	currentPlayer := r.getCurrentPlayer()

	if r.getCurrentPlayer().IsUNO() {
		for i := 0; i < len(r.players); i++ {
			if i == r.currentPlayer {
				continue
			}
			r.players[i].onOtherPlayerUNO(currentPlayer)
		}
	}

	switch c.Action {
	case models.ActionNone:
		r.putNumberCard(c)
	case models.ActionWild, models.ActionWildDrawFour:
		r.putWildCard(c)
	default:
		r.putActionCard(c)
	}
}

func (r *Round) canPutCard(c models.Card) bool {
	if c.Action == models.ActionNone {
		if r.lastNumber != models.NumberNone || r.lastColor != models.ColorAny {
			return c.Number == r.lastNumber || c.Color == r.lastColor
		}
		return true
	}
	if c.Action == models.ActionWild || c.Action == models.ActionWildDrawFour {
		return true
	}
	return c.Color == r.lastColor
}

func (r *Round) getDrawCards() []models.Card {
	var cards = make([]models.Card, r.draw)
	for i := 0; i < r.draw; i++ {
		cards[i] = r.getCard()
	}
	r.draw = 1
	return cards
}

func (r *Round) IsEnd() bool {
}
