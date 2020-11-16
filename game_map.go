package main

import (
	"github.com/sidav/golibrl/graphic_primitives"
)

type gameMap struct {
	player                     *pawn
	currentPlayerVisibilityMap [][]bool
	pathfindingCostMap         [][]int
	tiles                      [][]d_tile
	pawns                      []*pawn
	furnitures                 []*furniture
	noises                     []*noise
	//items       []*i_item
	//projectiles []*projectile
}

func (dung *gameMap) getSize() (int, int) {
	return len(dung.tiles), len(dung.tiles[0])
}

func (dung *gameMap) isPawnPresent(ix, iy int) bool {
	x, y := dung.player.x, dung.player.y
	if ix == x && iy == y {
		return true
	}
	for i := 0; i < len(dung.pawns); i++ {
		x, y = dung.pawns[i].x, dung.pawns[i].y
		if ix == x && iy == y {
			return true
		}
	}
	return false
}

func (dung *gameMap) getPawnAt(x, y int) *pawn {
	px, py := dung.player.x, dung.player.y
	if px == x && py == y {
		return dung.player
	}
	for i := 0; i < len(dung.pawns); i++ {
		px, py = dung.pawns[i].x, dung.pawns[i].y
		if px == x && py == y {
			return dung.pawns[i]
		}
	}
	return nil
}

func (d *gameMap) removePawn(p *pawn) {
	for i := 0; i < len(d.pawns); i++ {
		if p == d.pawns[i] {
			d.pawns = append(d.pawns[:i], d.pawns[i+1:]...) // ow it's fucking... magic!
		}
	}
}

func (d *gameMap) getNumberOfTilesOfTypeAround(ttype tileCode, x, y int) int {
	number := 0
	for i := x-1; i <= x+1; i++ {
		for j := y-1; j <= y+1; j++ {
			if areCoordinatesValid(i, j) && i != x && j != y && d.tiles[i][j].code == ttype {
				number++
			}
		}
	}
	return number
}

func (dung *gameMap) isTilePassable(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return false
	}
	for _, f := range CURRENT_MAP.furnitures {
		if !f.getStaticData().canBeSteppedOn && f.x == x && f.y == y {
			return false
		}
	}
	return dung.tiles[x][y].isPassable()
}

func (dung *gameMap) isTileOpaque(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return true
	}
	return dung.tiles[x][y].isOpaque()
}

func (dung *gameMap) isTileADoor(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return false
	}
	return dung.tiles[x][y].isDoor()
}

func (dung *gameMap) openDoor(x, y int) {
	if !areCoordinatesValid(x, y) {
		return
	}
	dung.tiles[x][y].isOpened = true
}

func (dung *gameMap) visibleLineExists(fx, fy, tx, ty int, ignoreStart bool) bool {
	line := graphic_primitives.GetLine(fx, fy, tx, ty)
	for i, l := range *line {
		if i == len(*line)-1 {
			break
		}
		if i == 0 && ignoreStart {
			continue
		}
		if !areCoordinatesValid(l.X, l.Y) || dung.isTileOpaque(l.X, l.Y) {
			return false
		}
	}
	return true
}

// true if action has been commited
func (dung *gameMap) movePawnOrOpenDoorByVector(p *pawn, mayOpenDoor bool, vx, vy int) bool {
	x, y := p.getCoords()
	x += vx
	y += vy
	if dung.isTilePassable(x, y) {
		p.x = x
		p.y = y
		if p.isRunning {
			p.spendTurnsForAction(p.getStaticData().timeForRunning)
		} else {
			p.spendTurnsForAction(p.getStaticData().timeForWalking)
		}
		dung.createNoise(p.createMovementNoise())
		return true
	}
	if dung.isTileADoor(x, y) && mayOpenDoor {
		dung.tiles[x][y].isOpened = true
		p.spendTurnsForAction(10)
		return true
	}
	return false
}

func (dung *gameMap) isTilePassableAndNotOccupied(x, y int) bool {
	return dung.isTilePassable(x, y) && !dung.isPawnPresent(x, y)
}

func (dung *gameMap) createNoise(n *noise) {
	n.turnCreatedAt = CURRENT_TURN
	n.duration = 9
	dung.noises = append(dung.noises, n)
}

func (dung *gameMap) cleanupNoises() {
	i := 0
	for _, n := range dung.noises {
		if n.turnCreatedAt+n.duration >= CURRENT_TURN {
			dung.noises[i] = n
			i++
		}
	}
	for j:=i; j<len(dung.noises); j++ {
		dung.noises[j] = nil
	}
	dung.noises = dung.noises[:i]
}
