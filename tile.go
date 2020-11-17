package main

import cw "github.com/sidav/golibrl/console"

type tileCode uint8

const (
	TILE_UNDEFINED tileCode = iota
	TILE_WALL
	TILE_FLOOR
	TILE_DOOR
	TILE_WINDOW
)

type tileStaticData struct {
	blocksMovement, blocksVision bool
	appearance                   *consoleCell
}

var tileStaticTable = map[tileCode] tileStaticData {
	TILE_UNDEFINED: {
		blocksMovement: false,
		blocksVision:   false,
		appearance: &consoleCell{
			appearance: '?',
			color:      cw.MAGENTA,
			inverse:    true,
		},
	},
	TILE_WALL: {
		blocksMovement: false,
		blocksVision:   true,
		appearance: &consoleCell{
			appearance: ' ',
			color:      cw.DARK_RED,
			inverse:    true,
		},
	},
	TILE_DOOR: {
		blocksMovement: false,
		blocksVision:   true,
		appearance: &consoleCell{
			appearance: '+',
			color:      cw.DARK_MAGENTA,
			inverse:    false,
		},
	},
	TILE_FLOOR: {
		blocksMovement: true,
		blocksVision:   false,
		appearance: &consoleCell{
			appearance: '.',
			color:      cw.YELLOW,
			inverse:    false,
		},
	},
	TILE_WINDOW: {
		blocksMovement: false,
		blocksVision:   false,
		appearance: &consoleCell{
			appearance: '#',
			color:      cw.CYAN,
			inverse:    false,
		},
	},
}

type d_tile struct {
	code tileCode
	wasSeenByPlayer bool
	lightLevel int
	isOpened bool // only if tile is a door
}

func (t *d_tile) getAppearance() *consoleCell {
	if t.isOpened {
		return &consoleCell{
			appearance: '\\',
			color:      cw.DARK_MAGENTA,
			inverse:    false,
		}
	}
	return tileStaticTable[t.code].appearance
}

func (t *d_tile) isDoor() bool {
	return t.code == TILE_DOOR
}

func (t *d_tile) isPassable() bool {
	if t.isOpened {
		return true
	}
	return tileStaticTable[t.code].blocksMovement
}

func (t *d_tile) isOpaque() bool {
	if t.isOpened {
		return false
	}
	return tileStaticTable[t.code].blocksVision
}
