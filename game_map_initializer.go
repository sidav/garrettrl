package main

import (
	"fmt"
	generator2 "parcelcreationtool/generator"
)

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
	dung.spawnFurnitureFromGenerated(generatedMap)
	dung.addRandomFurniture()
	dung.spawnEnemiesAtRoutes(generatedMap)
	dung.spawnRoamingEnemies(5)
	dung.distributeLootBetweenCabinets(1000)
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
	CURRENT_MAP.player = initNewPawn(PAWN_PLAYER, 1, 1, false)
	CURRENT_MAP.player.inv = &inventory{}
	CURRENT_MAP.player.inv.init()
	// check if generated map has an entry point
	// and select one at random
	entrypoints := make([][2]int, 0)
	for _, i := range l.Items {
		if i.Name == "ENTRYPOINT" {
			entrypoints = append(entrypoints, [2]int{i.X, i.Y})
		}
	}
	if len(entrypoints) > 0 {
		randEntryIndex := rnd.Rand(len(entrypoints))
		CURRENT_MAP.player.x = entrypoints[randEntryIndex][0]
		CURRENT_MAP.player.y = entrypoints[randEntryIndex][1]
	}
}

func (dung *gameMap) spawnFurnitureFromGenerated(l *generator2.Level) {
	// check if generated map has an entry point
	for _, i := range l.Items {
		if i.Name == "TORCH" {
			newF := furniture{code: FURNITURE_TORCH, x: i.X, y: i.Y, isLit: true}
			dung.furnitures = append(dung.furnitures, &newF)
		}
		if i.Name == "TABLE" {
			dung.furnitures = append(dung.furnitures, &furniture{code: FURNITURE_TABLE, x: i.X, y: i.Y})
		}
		if i.Name == "CABINET" {
			dung.furnitures = append(dung.furnitures, &furniture{code: FURNITURE_CABINET, x: i.X, y: i.Y})
		}
		if i.Name == "BUSH" {
			dung.furnitures = append(dung.furnitures, &furniture{code: FURNITURE_BUSH, x: i.X, y: i.Y})
		}
	}
}

func (dung *gameMap) addRandomFurniture() {
	w, h := CURRENT_MAP.getSize()
	// tables
	const TABLES = 0
	suitableTableCoords := make([][2]int, 0)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// placement rule
			if dung.isTilePassableAndNotOccupied(x, y) &&
				dung.getNumberOfTilesOfTypeAround(TILE_WALL, x, y) <= 3 &&
				dung.getNumberOfTilesOfTypeAround(TILE_FLOOR, x, y) > 2 {
				suitableTableCoords = append(suitableTableCoords, [2]int{x, y})
			}
		}
	}
	log.AppendMessage(fmt.Sprintf("Found %d suitable coords.", len(suitableTableCoords)))
	if len(suitableTableCoords) == 0 {
		log.AppendMessage("NO TABLE COORDS FOUND")
		return
	}
	currTableNum := 0
	currTableCoordIndex := rnd.Rand(len(suitableTableCoords))
	needChangeIndex := true
	for currTableNum < len(suitableTableCoords) && currTableNum < TABLES {
		for needChangeIndex {
			needChangeIndex = false
			currTableCoordIndex = rnd.Rand(len(suitableTableCoords))
			for _, f := range dung.furnitures {
				if f.x == suitableTableCoords[currTableCoordIndex][0] &&
					f.y == suitableTableCoords[currTableCoordIndex][1] {
					needChangeIndex = true
					break
				}
			}
		}
		dung.furnitures = append(dung.furnitures, &furniture{
			code: FURNITURE_TABLE,
			x:    suitableTableCoords[currTableCoordIndex][0],
			y:    suitableTableCoords[currTableCoordIndex][1],
			inv:  nil,
		})
		currTableNum++
		needChangeIndex = true
	}
}

func (dung *gameMap) spawnEnemiesAtRoutes(l *generator2.Level) {
	for r_index := range l.Routes {
		r := l.Routes[r_index]
		if len(r.Waypoints) > 0 {
			newEnemy := initNewPawn(PAWN_GUARD, r.Waypoints[0].X, r.Waypoints[0].Y, true)
			newEnemy.ai.route = &r
			newEnemy.ai.currentState = AI_PATROLLING
			dung.pawns = append(dung.pawns, newEnemy)
		}
	}
}

func (dung *gameMap) spawnRoamingEnemies(count int) {
	x := -1
	y := -1
	w, h := dung.getSize()
	for i := 0; i < count; i++ {
		for !dung.isTilePassableAndNotOccupied(x, y) {
			x, y = rnd.Rand(w), rnd.Rand(h)
		}
		newEnemy := initNewPawn(PAWN_GUARD, x, y, true)
		dung.pawns = append(dung.pawns, newEnemy)
	}
}

func (d *gameMap) distributeLootBetweenCabinets(minimumGoldAmount int) {
	totalCabinetsOnMap := 0
	for _, f := range d.furnitures {
		if f.code == FURNITURE_CABINET {
			totalCabinetsOnMap++
		}
	}
	avgGoldPerCabinet := minimumGoldAmount/totalCabinetsOnMap
	minGoldPerCabinet := avgGoldPerCabinet - 25
	if minGoldPerCabinet < 0 {
		minGoldPerCabinet = 0
	}
	maxGoldPerCabinet := avgGoldPerCabinet + 75
	for _, f := range d.furnitures {
		if f.code == FURNITURE_CABINET {
			f.inv = &inventory{}
			f.inv.init()
			f.inv.gold = rnd.RandInRange(minGoldPerCabinet, maxGoldPerCabinet)
			if rnd.OneChanceFrom(2) {
				f.inv.arrows[0].amount = 1
			}
			if rnd.OneChanceFrom(5) {
				f.inv.arrows[1].amount = 1
			}
			if rnd.OneChanceFrom(10) {
				f.inv.arrows[2].amount = 1
			}
			if rnd.OneChanceFrom(4) {
				f.inv.arrows[3].amount = 1
			}
		}
	}
}
