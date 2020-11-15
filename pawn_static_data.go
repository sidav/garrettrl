package main

import "github.com/sidav/golibrl/console"

type pawnCode uint8

const (
	PAWN_GUARD pawnCode = iota
	PAWN_PLAYER
)

type pawnStaticData struct {
	ccell *consoleCell
	name  string
	maxhp int

	timeForWalking, timeForRunning                           int
	runningNoiseIntensity, walkingNoiseIntensity             int
	sightRangeCalm, sightRangeAlerted, sightRangeAlertedDark int
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
		name:                  "Guard",
		maxhp:                 3,
		timeForWalking:        10,
		timeForRunning:        8,
		runningNoiseIntensity: 10,
		walkingNoiseIntensity: 7,

		sightRangeAlerted:     9,
		sightRangeAlertedDark: 3,
		sightRangeCalm:        6,
	},
	PAWN_PLAYER: {
		ccell: &consoleCell{
			appearance: '@',
			color:      console.WHITE,
			inverse:    false,
		},
		sightRangeAlerted:     10,
		name:                  "Taffer",
		maxhp:                 3,
		timeForWalking:        10,
		timeForRunning:        6,
		runningNoiseIntensity: 5,
		walkingNoiseIntensity: 0,
	},
}
