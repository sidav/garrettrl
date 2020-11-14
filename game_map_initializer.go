package main

import cw "github.com/sidav/golibrl/console"

var testMap = []string{
	"########################",
	"#......#...............#",
	"#......#...............#",
	"#......#...............#",
	"#......#.......#######+#",
	"###+####.......#.......#",
	"#..............#.......#",
	"#..............#.......#",
	"########################",
}

func (dung *gameMap) initialize_level() { //crap of course
	dung.pawns = make([]*pawn, 0)
}

func (dung *gameMap) initTilesArrayForSize(sx, sy int) {
	dung.tiles = make([][]d_tile, sx)
	for i := range dung.tiles {
		dung.tiles[i] = make([]d_tile, sy)
	}
}

func (dung *gameMap) init_placeItemsAndEnemies() {

}

func (dung *gameMap) MakeMapFromGenerated(generated_map *[]string) {
	levelsizex = len(*generated_map)
	levelsizey = len((*generated_map)[0])
	dung.initTilesArrayForSize(levelsizex, levelsizey)

	for x := 0; x < levelsizex; x++ {
		for y := 0; y < levelsizey; y++ {
			currDungCell := &dung.tiles[x][y]
			currGenCell := (*generated_map)[x][y] //GetCell(x, y)
			switch currGenCell {
			case '#':
				currDungCell.code = TILE_WALL
			case '.':
				currDungCell.code = TILE_FLOOR
			case '+':
				currDungCell.code = TILE_DOOR
			default:
				currDungCell.code = TILE_UNDEFINED
			}
		}
	}
	// dung.furnitures = append(dung.furnitures, &furniture{code: FURNITURE_TORCH, x: 4, y: 5})
	dung.furnitures = append(dung.furnitures, &furniture{code: FURNITURE_TORCH, x: 1, y: 7})
	dung.pawns = append(dung.pawns, &pawn{
		ccell:         &consoleCell{
			appearance: 'G',
			color:      cw.RED,
			inverse:    false,
		},
		hp:            0,
		maxhp:         0,
		x:             7,
		y:             7,
		nextTurnToAct: 0,
		sightRange:    6,
		name:          "Guard",
		ai:            &aiData{},
	})
}


func (dung *gameMap) spawnPlayerAtRandomPosition() {
}

func (dung *gameMap) spawnPawnAtRandomPosition(name string, count int) {
}

func (dung *gameMap) spawnItemAtRandomPosition(name string, count int) {
}
