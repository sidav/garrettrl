package main

type aiState uint8

const (
	AI_ROAM aiState = iota
	AI_ALERTED
	AI_ENGAGING
)

type aiData struct {
	currentState aiState
	targetPawn *pawn
	dirx, diry int // for roaming
}

func (p *pawn) ai_checkSituation() {
	switch p.ai.currentState {
	case AI_ROAM:
		p.ai_checkRoam()
	case AI_ALERTED:
	}
}

func (p *pawn) ai_act() {
	switch p.ai.currentState {
	case AI_ROAM:
		p.ai_actRoam()
	case AI_ALERTED:
	}
}

func (p *pawn) ai_checkRoam() {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		if CURRENT_MAP.tiles[px][py].lightLevel > 0 {
			p.ai.targetPawn = CURRENT_MAP.player
			p.ai.currentState = AI_ALERTED
			return
		}
	}
	p.ai.currentState = AI_ROAM
}

func (p *pawn) ai_actRoam() {
	ai := p.ai
	for ai.dirx == 0 && ai.diry == 0 {
		ai.dirx, ai.diry = rnd.RandomUnitVectorInt()
	}
	newx, newy := p.x + ai.dirx, p.y + ai.diry
	if CURRENT_MAP.isTilePassableAndNotOccupied(newx, newy) {
		CURRENT_MAP.movePawnOrOpenDoorByVector(p, true, ai.dirx, ai.diry)
	} else {
		ai.dirx = 0
		ai.diry = 0
	}
}

func (p *pawn) ai_checkAlerted() {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		if CURRENT_MAP.tiles[px][py].lightLevel > 0 {
			p.ai.targetPawn = CURRENT_MAP.player
			p.ai.currentState = AI_ALERTED
			return
		}
	}
	p.ai.currentState = AI_ROAM
}
