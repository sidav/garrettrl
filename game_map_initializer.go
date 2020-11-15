package main

import (
	generator2 "parcelcreationtool/generator"
)

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
	dung.pathfindingCostMap = make([][]int, sx)
	for i := range dung.pathfindingCostMap {
		dung.pathfindingCostMap[i] = make([]int, sy)
	}
}

func (dung *gameMap) generateAndInitMap() {
	dung.initialize_level()
	generator := generator2.Generator{}
	generatedMap := generator.Generate("parcels", "templates", 0, 0, 9)
	generatedMapString := make([]string, 0)
	for i := range generatedMap.Terrain {
		currStr := ""
		for j := range generatedMap.Terrain[i] {
			currStr += string(generatedMap.Terrain[i][j])
		}
		generatedMapString = append(generatedMapString, currStr)
	}
	dung.applyRuneMap(&generatedMapString)
	dung.spawnPlayer(generatedMap)
	dung.spawnFurniture(generatedMap)
	dung.spawnEnemiesAtRoutes(generatedMap)
	dung.spawnRoamingEnemies(5)
}

func (dung *gameMap) applyRuneMap(generated_map *[]string) {
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
}

func (dung *gameMap) spawnPlayer(l *generator2.Level) {
	CURRENT_MAP.player = &pawn{
		code: PAWN_PLAYER,
		hp:             3,
		x:              1,
		y:              1,
		nextTurnToAct:  0,
	}
	// check if generated map has an entry point
	for _, i := range l.Items {
		if i.Name == "ENTRYPOINT" {
			CURRENT_MAP.player.x = i.X
			CURRENT_MAP.player.y = i.Y
		}
	}
}

func (dung *gameMap) spawnFurniture(l *generator2.Level) {
	// check if generated map has an entry point
	for _, i := range l.Items {
		if i.Name == "TORCH" {
			dung.furnitures = append(dung.furnitures, &furniture{code: FURNITURE_TORCH, x: i.X, y: i.Y})
		}
	}
}

func (dung *gameMap) spawnEnemiesAtRoutes(l *generator2.Level) {
	for r_index := range l.Routes {
		r := l.Routes[r_index]
		if len(r.Waypoints) > 0 {
			newEnemy := pawn{
				code: PAWN_GUARD,
				hp:                0,
				x:                 r.Waypoints[0].X,
				y:                 r.Waypoints[0].Y,
				nextTurnToAct:     0,
				ai:                &aiData{},
			}
			newEnemy.ai.route = &r
			newEnemy.ai.currentState = AI_PATROLLING
			dung.pawns = append(dung.pawns, &newEnemy)
		}
	}
}

func (dung *gameMap) spawnRoamingEnemies(count int) {
	var x, y int
	w, h := dung.getSize()
	for i := 0; i < count; i++ {
		for !dung.isTilePassableAndNotOccupied(x, y) {
			x, y = rnd.Rand(w), rnd.Rand(h)
		}
		dung.pawns = append(dung.pawns, &pawn{
			code: PAWN_GUARD,
			hp:                0,
			x:                 x,
			y:                 y,
			nextTurnToAct:     0,
			ai:                &aiData{},
		})
	}
}
