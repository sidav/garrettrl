package main

type (
	pawn struct {
		ccell                                      *consoleCell
		hp, maxhp, x, y, nextTurnToAct int
		sightRangeCalm, sightRangeAlerted int
		name                                       string
		ai *aiData
	}
)

func (p *pawn) isDead() bool {
	return p.hp <= 0
}

func (p *pawn) spendTurnsForAction(turns int) {
	p.nextTurnToAct = CURRENT_TURN + turns
}

func (p *pawn) isTimeToAct() bool {
	return p.nextTurnToAct <= CURRENT_TURN
}

func (p *pawn) getCoords() (int, int) {
	return p.x, p.y
}

func (p *pawn) getHpPercent() int {
	return p.hp * 100 / p.maxhp
}
