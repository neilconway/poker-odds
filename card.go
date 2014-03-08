/*
 * Copyright 2011 Colin Patrick McCabe
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 2.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
)

const (
	PARSE_STATE_EAT_VAL = iota
	PARSE_STATE_EAT_VAL_SAW_1
	PARSE_STATE_EAT_SUIT
)

const (
	DIAMONDS = iota
	CLUBS
	HEARTS
	SPADES
)

type Card struct {
	val  int
	suit int
}

func (p Card) Copy() *Card {
	return &Card{p.val, p.suit}
}

func rankToStr(r int) string {
	switch r {
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	case 14:
		return "A"
	default:
		return fmt.Sprintf("%d", r)
	}
}

func suitToStr(s int) string {
	switch s {
	case CLUBS:
		return "♣C"
	case DIAMONDS:
		return "♦D"
	case HEARTS:
		return "♥H"
	case SPADES:
		return "♠S"
	default:
		panic(fmt.Sprintf("invalid suit %d", s))
	}
}

func (c *Card) String() string {
	return fmt.Sprintf("%s%s", rankToStr(c.val), suitToStr(c.suit))
}

/* It's important that the cards compare in this order. It makes detecting
 * straights easier because cards of a similar value (as opposed to suit) are
 * adjacent. Don't change this sort order without updating hand.go
 */
func (p *Card) Compare(rhs *Card) int {
	if p.suit < rhs.suit {
		return -1
	}
	if p.suit > rhs.suit {
		return 1
	}
	if p.val < rhs.val {
		return -1
	}
	if p.val > rhs.val {
		return 1
	}
	return 0
}

func StrToCard(str string, cnt *int) (myCard *Card) {
	myCard = new(Card)
	var parseState = PARSE_STATE_EAT_VAL
	for *cnt < len(str) {
		var c = str[*cnt]
		*cnt++
		switch {
		case parseState == PARSE_STATE_EAT_VAL:
			switch {
			case c == ' ' || c == '\t':
				continue
			case c == '1':
				parseState = PARSE_STATE_EAT_VAL_SAW_1
			case c >= '2' && c <= '9':
				myCard.val = (int)(c - '0')
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'J':
				myCard.val = JACK_VAL
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'Q':
				myCard.val = QUEEN_VAL
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'K':
				myCard.val = KING_VAL
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'A':
				myCard.val = ACE_VAL
				parseState = PARSE_STATE_EAT_SUIT
			default:
				return nil
			}
		case parseState == PARSE_STATE_EAT_VAL_SAW_1:
			switch {
			case c == '0':
				myCard.val = 10
				parseState = PARSE_STATE_EAT_SUIT
			default:
				return nil
			}
		case parseState == PARSE_STATE_EAT_SUIT:
			switch {
			case c == 'C':
				myCard.suit = CLUBS
			case c == 'D':
				myCard.suit = DIAMONDS
			case c == 'H':
				myCard.suit = HEARTS
			case c == 'S':
				myCard.suit = SPADES
			default:
				return nil
			}
			return myCard
		}
	}
	*cnt = -1
	return nil
}

type CardSlice []*Card

func (arr CardSlice) Copy() CardSlice {
	ret := make(CardSlice, len(arr))
	for i := range arr {
		ret[i] = arr[i].Copy()
	}
	return ret
}

/* In poker, the 'kicker' breaks ties between hands of the same type.
 * An example:
 * Hand1: 6D QC QS KS KC      Kicker: 6D
 * Hand2: 9H QH QD KC KS      Kicker: 9H (wins)
 *
 * Hand1: 6D 8C 9S JS JC      Kicker: 6D 8C 9S
 * Hand2: 2H 3H 10D JD JH     Kicker: 2H 3H 10D (wins)
 *
 * Kickers are compared in lexicographical order, starting with the highest
 * valued card. Suit is irrelevant; only card value matters.
 *
 * This function can return 0 even if the two CardSlices are different.
 */
func (arr CardSlice) CompareKicker(rhs CardSlice) int {
	var a int // needs to be signed
	var b int
	a = len(arr) - 1
	b = len(rhs) - 1
	for {
		if a < 0 {
			if b < 0 {
				return 0
			} else {
				return -1
			}
		} else if b < 0 {
			return 1
		}
		if arr[a].val < arr[b].val {
			return -1
		} else if arr[a].val > arr[b].val {
			return 1
		}
		// ignore suit!
		a--
		b--
	}
	return 0
}

func (arr CardSlice) Identical(rhs CardSlice) bool {
	if len(arr) != len(rhs) {
		return false
	}
	for i := range arr {
		if (arr[i].val != arr[i].val) || (arr[i].suit != arr[i].suit) {
			return false
		}
	}
	return true
}

func (arr CardSlice) Len() int {
	return len(arr)
}

func (arr CardSlice) Less(i, j int) bool {
	if arr[i].val < arr[j].val {
		return true
	}
	if arr[i].val > arr[j].val {
		return false
	}
	if arr[i].suit < arr[j].suit {
		return true
	}
	return false
}

func (arr CardSlice) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr CardSlice) String() string {
	ret := ""
	sep := ""
	for i := range arr {
		ret += fmt.Sprintf("%s%s", sep, arr[i].String())
		sep = ", "
	}
	return ret
}

// Could do this smarter if we knew that we were sorted...
func (arr CardSlice) HasDuplicates() *Card {
	for i := range arr {
		for j := range arr {
			if i == j {
				continue
			}
			if arr[i].Compare(arr[j]) == 0 {
				return arr[i]
			}
		}
	}
	return nil
}

func StrToCards(str string) (ret CardSlice, cnt int) {
	for cnt = 0; cnt != -1; {
		var c = StrToCard(str, &cnt)
		if c == nil {
			return
		}
		ret = append(ret, c)
	}
	return ret, cnt
}
