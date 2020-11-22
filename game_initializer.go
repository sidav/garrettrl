package main

import (
	"encoding/json"
	"fmt"
	cw "github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/console_menu"
	"io/ioutil"
	generator2 "parcelcreationtool/generator"
)

type missionInitializer struct {
}

func (m *missionInitializer) initializeMission(missionNumber int) { //crap of course
	CURRENT_MAP = gameMap{}
	CURRENT_MAP.pawns = make([]*pawn, 0)
	filesDir := fmt.Sprintf("missions/mission%d/", missionNumber)
	m.generateAndInitMap(filesDir)
}

func (m *missionInitializer) generateAndInitMap(filesPath string) {
	generator := generator2.Generator{}
	generatedMap := generator.Generate(filesPath + "parcels", filesPath+ "templates", 0, 0, 9)
	generatedMapString := make([]string, 0)
	for i := range generatedMap.Terrain {
		currStr := ""
		for j := range generatedMap.Terrain[i] {
			currStr += string(generatedMap.Terrain[i][j])
		}
		generatedMapString = append(generatedMapString, currStr)
	}

	// mission unmarshalling
	mis := &Mission{}
	jsn, err := ioutil.ReadFile(filesPath+"mission.json")
	if err == nil {
		json.Unmarshal(jsn, mis)
	} else {
		panic(err)
	}
	renderer.putTextInRect(mis.BriefingText, 0, 0, 0)
	cw.ReadKey()
	difficulty := console_menu.ShowSingleChoiceMenu("Select difficulty:", "", []string{"Easy", "Medium", "Hard"})

	m.applyRuneMap(&generatedMapString)
	m.spawnPlayer(generatedMap)
	m.spawnFurnitureFromGenerated(generatedMap)
	m.addRandomFurniture()
	m.spawnEnemiesAtRoutes(generatedMap)
	m.spawnRoamingEnemies(mis.AdditionalGuardsNumber[difficulty])
	m.distributeLootBetweenCabinets(mis.TotalLoot[difficulty])
}

func (m *missionInitializer) applyRuneMap(generated_map *[]string) {
	levelsizex = len(*generated_map)
	levelsizey = len((*generated_map)[0])
	CURRENT_MAP.initTilesArrayForSize(levelsizex, levelsizey)

	for x := 0; x < levelsizex; x++ {
		for y := 0; y < levelsizey; y++ {
			currCURRENT_MAPCell := &CURRENT_MAP.tiles[x][y]
			currGenCell := (*generated_map)[x][y] //GetCell(x, y)
			switch currGenCell {
			case '#':
				currCURRENT_MAPCell.code = TILE_WALL
			case '.':
				currCURRENT_MAPCell.code = TILE_FLOOR
			case '+':
				currCURRENT_MAPCell.code = TILE_DOOR
			case '\'':
				currCURRENT_MAPCell.code = TILE_WINDOW
			default:
				currCURRENT_MAPCell.code = TILE_UNDEFINED
			}
		}
	}
}

func (m *missionInitializer) spawnPlayer(l *generator2.Level) {
	CURRENT_MAP.player = initNewPawn(PAWN_PLAYER, 1, 1, false)
	CURRENT_MAP.player.inv = &inventory{}
	CURRENT_MAP.player.inv.init()
	CURRENT_MAP.player.inv.arrows[0].amount = 2
	CURRENT_MAP.player.inv.arrows[1].amount = 1
	// check if generated map has an entry point
	// and select one at random
	entrypoints := make([][2]int, 0)
	for _, i := range l.Items {
		if i.Name == "ENTRYPOINT" {
			entrypoints = append(entrypoints, [2]int{i.X, i.Y})
		}
	}
	if len(entrypoints) > 0 {
		randEntryIndex := rnd.Rand(101) % len(entrypoints) // TODO: find why that hack is even needed.
		log.AppendMessage(fmt.Sprintf("Used %dth entry from %d entrypoints found.", randEntryIndex+1, len(entrypoints)))
		CURRENT_MAP.player.x = entrypoints[randEntryIndex][0]
		CURRENT_MAP.player.y = entrypoints[randEntryIndex][1]
	}
}

func (m *missionInitializer) spawnFurnitureFromGenerated(l *generator2.Level) {
	for _, i := range l.Items {
		switch i.Name {
		case "ENTRYPOINT":
			continue // do nothing
		case "TORCH":
			newF := furniture{code: FURNITURE_TORCH, x: i.X, y: i.Y, isLit: true}
			CURRENT_MAP.furnitures = append(CURRENT_MAP.furnitures, &newF)
		case "TABLE":
			CURRENT_MAP.furnitures = append(CURRENT_MAP.furnitures, &furniture{code: FURNITURE_TABLE, x: i.X, y: i.Y})
		case "CABINET":
			CURRENT_MAP.furnitures = append(CURRENT_MAP.furnitures, &furniture{code: FURNITURE_CABINET, x: i.X, y: i.Y})
		case "BUSH":
			CURRENT_MAP.furnitures = append(CURRENT_MAP.furnitures, &furniture{code: FURNITURE_BUSH, x: i.X, y: i.Y})
		default:
			CURRENT_MAP.furnitures = append(CURRENT_MAP.furnitures, &furniture{code: FURNITURE_UNDEFINED, x: i.X, y: i.Y})
		}
	}
}

func (m *missionInitializer) addRandomFurniture() {
	w, h := CURRENT_MAP.getSize()
	// tables
	const TABLES = 0
	suitableTableCoords := make([][2]int, 0)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// placement rule
			if CURRENT_MAP.isTilePassableAndNotOccupied(x, y) &&
				CURRENT_MAP.getNumberOfTilesOfTypeAround(TILE_WALL, x, y) <= 3 &&
				CURRENT_MAP.getNumberOfTilesOfTypeAround(TILE_FLOOR, x, y) > 4 {
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
			for _, f := range CURRENT_MAP.furnitures {
				if f.x == suitableTableCoords[currTableCoordIndex][0] &&
					f.y == suitableTableCoords[currTableCoordIndex][1] {
					needChangeIndex = true
					break
				}
			}
		}
		CURRENT_MAP.furnitures = append(CURRENT_MAP.furnitures, &furniture{
			code: FURNITURE_TABLE,
			x:    suitableTableCoords[currTableCoordIndex][0],
			y:    suitableTableCoords[currTableCoordIndex][1],
			inv:  nil,
		})
		currTableNum++
		needChangeIndex = true
	}
}

func (m *missionInitializer) spawnEnemiesAtRoutes(l *generator2.Level) {
	for r_index := range l.Routes {
		r := l.Routes[r_index]
		if len(r.Waypoints) > 0 {
			newEnemy := initNewPawn(PAWN_GUARD, r.Waypoints[0].X, r.Waypoints[0].Y, true)
			newEnemy.ai.route = &r
			newEnemy.ai.currentState = AI_PATROLLING
			CURRENT_MAP.pawns = append(CURRENT_MAP.pawns, newEnemy)
		}
	}
}

func (m *missionInitializer) spawnRoamingEnemies(roamingEnemiesCount int) {
	x := -1
	y := -1
	w, h := CURRENT_MAP.getSize()
	for i := 0; i < roamingEnemiesCount; i++ {
		for !CURRENT_MAP.isTilePassableAndNotOccupied(x, y) {
			x, y = rnd.Rand(w), rnd.Rand(h)
		}
		newEnemy := initNewPawn(PAWN_GUARD, x, y, true)
		CURRENT_MAP.pawns = append(CURRENT_MAP.pawns, newEnemy)
	}
}

func (m *missionInitializer) distributeLootBetweenCabinets(totalDesiredLootAmount int) {
	totalCabinetsOnMap := 0
	for _, f := range CURRENT_MAP.furnitures {
		if f.code == FURNITURE_CABINET {
			totalCabinetsOnMap++
		}
	}
	avgGoldPerCabinet := totalDesiredLootAmount / totalCabinetsOnMap
	minGoldPerCabinet := avgGoldPerCabinet - 25
	if minGoldPerCabinet < 0 {
		minGoldPerCabinet = 0
	}
	maxGoldPerCabinet := avgGoldPerCabinet + 75
	for _, f := range CURRENT_MAP.furnitures {
		if f.code == FURNITURE_CABINET {
			f.inv = &inventory{}
			f.inv.init()
			f.inv.gold = rnd.RandInRange(minGoldPerCabinet, maxGoldPerCabinet)
			// water
			if rnd.OneChanceFrom(3) {
				f.inv.arrows[0].amount = 1
			}
			// noise
			if rnd.OneChanceFrom(10) {
				f.inv.arrows[1].amount = 1
			}
			// gas
			if rnd.OneChanceFrom(10) {
				f.inv.arrows[2].amount = 1
			}
			// explosive
			if rnd.OneChanceFrom(15) {
				f.inv.arrows[3].amount = 1
			}
		}
	}
}