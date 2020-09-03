// +build ignore

package wtf

import (
	"errors"
	"gtihub.com/hariolate/tonneau/shared/models"
)

type Player struct {
	round *Round
	cards []models.Card

	announceUNO bool
}

func (p *Player) DrawCards() {
	cards := p.round.getDrawCards()
	p.cards = append(p.cards, cards...)
}

func (p *Player) onChooseColor() models.Color {
	// TODO
	panic("no impl")
}

func (p *Player) onTurn() {
	// TODO
	panic("no impl")
}

func (p *Player) Cards() []models.Card {
	return p.cards
}

func (p *Player) popCard(idx int) models.Card {
	card := p.cards[idx]
	p.cards = append(p.cards[:idx], p.cards[idx+1:]...)
	if len(p.cards) == 0 {
		if card.Action != models.ActionNone {
			p.cards = append(p.cards, p.round.getCard())
		}
	}
	return card
}

func (p *Player) PutCard(idx int) error {
	if !p.round.canPutCard(p.cards[idx]) {
		return errors.New("invalid card to put")
	}
	p.round.putCard(p.popCard(idx))
	return nil
}

func (p *Player) CanPutCard(idx int) bool {
	return p.round.canPutCard(p.cards[idx])
}

func (p *Player) AnnounceUNO() {
	p.announceUNO = true
}

func (p *Player) onOtherPlayerUNO(o *Player) {
	// TODO
	panic("not impl")
}

func (p *Player) IsUNO() bool {
	if len(p.cards) == 1 {
		return p.cards[0].Action == models.ActionNone
	}
	return false
}
