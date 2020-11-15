package main

import "github.com/sidav/golibrl/console"

type pawnCode uint8

const (
	PAWN_GUARD pawnCode = iota
	PAWN_PLAYER
)

type pawnStaticData struct {
	ccell                             *consoleCell
	sightRangeCalm, sightRangeAlerted int
	name                              string
	maxhp                             int
}

func (p *pawn) getStaticData() *pawnStaticData {
	pds := pawnStaticTable[p.code]
	return &pds
}

var pawnStaticTable = map[pawnCode]pawnStaticData{
	PAWN_GUARD: {
		ccell: &consoleCell{
			appearance: 'G',
			color:      console.RED,
			inverse:    false,
		},
		sightRangeCalm:    6,
		sightRangeAlerted: 9,
		name:              "Guard",
	},
	PAWN_PLAYER: {
		ccell: &consoleCell{
			appearance: '@',
			color:      console.WHITE,
			inverse:    false,
		},
		sightRangeCalm:    6,
		sightRangeAlerted: 9,
		name:              "Taffer",
	},
}
