package main

import "github.com/sidav/golibrl/fov/permissive_fov"

func (g *gameMap) getFieldOfVisionFor(p *pawn) *[][]bool {
	x, y := p.getCoords()
	return permissive_fov.GetFovMapFrom(x, y, 10, levelsizex, levelsizey, g.isTileOpaque)
}
