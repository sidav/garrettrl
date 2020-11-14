package main

import cw "github.com/sidav/golibrl/console"

type furniture struct {
	code furnitureCode
	x, y int
}

func (f *furniture) getStaticData() *furnitureStaticData {
	fsd := furnitureStaticTable[f.code]
	return &fsd
}

type furnitureCode uint8

const (
	FURNITURE_UNDEFINED furnitureCode = iota
	FURNITURE_TORCH
)

type furnitureStaticData struct {
	lightStrength int
	appearance *consoleCell
}

var furnitureStaticTable = map[furnitureCode] furnitureStaticData {
	FURNITURE_UNDEFINED: {
		lightStrength: 0,
		appearance: &consoleCell{
			appearance: '?',
			color:      cw.MAGENTA,
			inverse:    true,
		},
	},
	FURNITURE_TORCH: {
		lightStrength: 5,
		appearance: &consoleCell{
			appearance: '|',
			color:      cw.YELLOW,
			inverse:    true,
		},
	},
}
