package main

import cw "github.com/sidav/golibrl/console"

type (
	pawn struct {
		code                              pawnCode
		hp, x, y, nextTurnToAct    int
		ai                                *aiData
		isRunning bool
	}
)

func initNewPawn(code pawnCode, x, y int, hasAi bool) *pawn {
	newPawn := &pawn{
		code:          code,
		hp:            0,
		x:             x,
		y:             y,
		nextTurnToAct: 0,
		ai:            nil,
		isRunning:     false,
	}
	newPawn.hp = newPawn.getStaticData().maxhp
	if hasAi {
		newPawn.ai = &aiData{}
	}
	return newPawn
}

func (p *pawn) isDead() bool {
	return p.hp <= 0
}

func (p *pawn) spendTurnsForAction(turns int) {
	p.nextTurnToAct = CURRENT_TURN + turns
}

func (p *pawn) isTimeToAct() bool {
	return p.nextTurnToAct <= CURRENT_TURN
}

func (p *pawn) isInLight() bool {
	return CURRENT_MAP.tiles[p.x][p.y].lightLevel > 0
}

func (p *pawn) getCoords() (int, int) {
	return p.x, p.y
}

func (p *pawn) getHpPercent() int {
	return p.hp * 100 / p.getStaticData().maxhp
}

func (p *pawn) createMovementNoise() *noise {
	intensity := p.getStaticData().walkingNoiseIntensity
	if p.isRunning {
		intensity = p.getStaticData().runningNoiseIntensity
	}
	nse := &noise{
		x:               p.x,
		y:               p.y,
		intensity:       intensity,
		visual:          consoleCell{
			appearance: p.getStaticData().ccell.appearance,
			color:      cw.BLUE,
			inverse:    false,
		},
		suspicious:      p.isRunning,
		showOnlyNotSeen: true,
	}
	return nse
}
