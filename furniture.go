package main

import cw "github.com/sidav/golibrl/console"

type furniture struct {
	code  furnitureCode
	isLit bool
	x, y  int
	inv   *inventory
}

func (f *furniture) canBeLooted() bool {
	return f.inv != nil
}

func (f *furniture) getCurrentLightLevel() int {
	if f.isLit {
		return f.getStaticData().lightStrength
	} else {
		return 0
	}
}

func (f *furniture) getStaticData() *furnitureStaticData {
	fsd := furnitureStaticTable[f.code]
	return &fsd
}

type furnitureCode uint8

const (
	FURNITURE_UNDEFINED furnitureCode = iota
	FURNITURE_TORCH
	FURNITURE_CABINET
	FURNITURE_TABLE
	FURNITURE_BUSH
)

type furnitureStaticData struct {
	lightStrength int
	appearance    *consoleCell

	isExtinguishable bool // for torches
	canBeSteppedOn   bool // ONLY AS NON-COVER MOVE!
	canBeUsedAsCover bool
}

var furnitureStaticTable = map[furnitureCode]furnitureStaticData{
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
			color:      cw.DARK_GRAY,
			inverse:    false,
		},
		canBeSteppedOn: false,
		isExtinguishable: true,
	},
	FURNITURE_CABINET: {
		lightStrength: 0,
		appearance: &consoleCell{
			appearance: '&',
			color:      cw.DARK_GRAY,
			inverse:    false,
		},
		canBeSteppedOn: false,
	},
	FURNITURE_TABLE: {
		lightStrength: 0,
		appearance: &consoleCell{
			appearance: '=',
			color:      cw.WHITE,
			inverse:    false,
		},
		canBeSteppedOn:   false,
		canBeUsedAsCover: true,
	},
	FURNITURE_BUSH: {
		lightStrength: 0,
		appearance: &consoleCell{
			appearance: '*',
			color:      cw.DARK_GREEN,
			inverse:    false,
		},
		canBeSteppedOn:   true,
		canBeUsedAsCover: true,
	},
}
