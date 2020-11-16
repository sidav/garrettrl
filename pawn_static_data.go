package main

import "github.com/sidav/golibrl/console"

type pawnCode uint8
type responseSituation uint8

const (
	PAWN_GUARD pawnCode = iota
	PAWN_PLAYER
)

const (
	SITUATION_NOISE responseSituation = iota
	SITUATION_ENEMY_SIGHTED
	SITUATION_ENEMY_DISAPPEARED
	SITUATION_SEARCH_STOPPED
)

type pawnStaticData struct {
	ccell *consoleCell
	name  string
	maxhp int

	timeForWalking, timeForRunning                           int
	runningNoiseIntensity, walkingNoiseIntensity             int
	sightRangeCalm, sightRangeAlerted, sightRangeAlertedDark int

	responsesForSituations map[responseSituation][]string
}

func (p *pawn) getStaticData() *pawnStaticData {
	pds := pawnStaticTable[p.code]
	return &pds
}

func (p *pawnStaticData) getRandomResponseTo(situation responseSituation) string {
	resp := p.responsesForSituations[situation][rnd.Rand(len(p.responsesForSituations[situation]))]
	return resp
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
		responsesForSituations: map[responseSituation][]string{
			SITUATION_NOISE: {
				"What was that?",
				"Huh?",
				"Did you hear that?",
			},
			SITUATION_ENEMY_SIGHTED: {
				"There you are!",
				"Don't run, taffer!",
				"Haha! I see ya, thief!",
			},
			SITUATION_ENEMY_DISAPPEARED: {
				"Where did he go?",
				"I'll find thee, taffer.",
				"You think you can hide?",
			},
			SITUATION_SEARCH_STOPPED: {
				"Nothing.",
				"Taff it.",
				"I'll better return.",
			},
		},
	},
	PAWN_PLAYER: {
		ccell: &consoleCell{
			appearance: '@',
			color:      console.WHITE,
			inverse:    false,
		},
		sightRangeAlerted:     10,
		name:                  "Taffer",
		maxhp:                 5,
		timeForWalking:        10,
		timeForRunning:        6,
		runningNoiseIntensity: 5,
		walkingNoiseIntensity: 0,
	},
}
