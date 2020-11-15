package main

import "github.com/sidav/golibrl/astar"

func (d *gameMap) recalculatePathfindingCostMap(considerPawns bool) {
	w, h := d.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if d.isTilePassable(x, y) || d.isTileADoor(x, y) {
				d.pathfindingCostMap[x][y] = 0
			} else {
				d.pathfindingCostMap[x][y] = -1
			}
		}
	}
	if considerPawns {
		for _, p := range CURRENT_MAP.pawns {
			d.pathfindingCostMap[p.x][p.y] = -1
		}
	}
	// consider furniture
	for _, f := range CURRENT_MAP.furnitures {
		if !f.getStaticData().canBeSteppedOn {
			d.pathfindingCostMap[f.x][f.y] = -1
		}
	}
}

func (d *gameMap) getPathFromTo(fx, fy, tx, ty int, considerPawns bool) *astar.Cell {
	d.recalculatePathfindingCostMap(considerPawns)
	path := astar.FindPath(&d.pathfindingCostMap, fx, fy, tx, ty, true, 100, true, true)
	return path
}
