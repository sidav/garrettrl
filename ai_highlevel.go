package main

import "parcelcreationtool/parcel"

type aiState uint8

const (
	AI_ROAM aiState = iota
	AI_PATROLLING
	AI_ALERTED
	AI_SEARCHING
)

type aiData struct {
	currentState aiState
	targetPawn   *pawn
	dirx, diry   int // for roaming

	route                *parcel.Route // for patrol
	currentWaypointIndex int
}

func (p *pawn) ai_checkSituation() {
	switch p.ai.currentState {
	case AI_ROAM, AI_PATROLLING:
		p.ai_checkRoam()
	case AI_ALERTED:
		p.ai_checkAlerted()
	default:
		log.AppendMessage("No CHECK func for some ai state!")
	}
}

func (p *pawn) ai_act() {
	switch p.ai.currentState {
	case AI_ROAM:
		p.ai_actRoam()
	case AI_PATROLLING:
		p.ai_actPatrolling()
	case AI_ALERTED:
		p.ai_actAlerted()
	default:
		log.AppendMessage("No ACT func for some ai state!")
	}
}

func (p *pawn) ai_checkRoam() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		return
	}
	p.ai_resetStateToCalm()
}

func (p *pawn) ai_actRoam() {
	ai := p.ai
	for ai.dirx == 0 && ai.diry == 0 {
		ai.dirx, ai.diry = rnd.RandomUnitVectorInt()
	}
	if !p.ai_TryMoveOrOpenDoorOrAlert(ai.dirx, ai.diry) {
		ai.dirx = 0
		ai.diry = 0
	}
}

func (p *pawn) ai_actPatrolling() {
	ai := p.ai
	currWaypoint := ai.route.Waypoints[ai.currentWaypointIndex]
	px, py := p.getCoords()
	if px == currWaypoint.X && py == currWaypoint.Y {
		ai.currentWaypointIndex++
	}
	if ai.currentWaypointIndex >= len(ai.route.Waypoints) {
		ai.currentWaypointIndex = 0
	}
	currWaypoint = ai.route.Waypoints[ai.currentWaypointIndex]
	path := CURRENT_MAP.getPathFromTo(p.x, p.y, currWaypoint.X, currWaypoint.Y, false)
	dirx, diry := path.GetNextStepVector()
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}

func (p *pawn) ai_checkAlerted() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		return
	}
	p.ai_resetStateToCalm()
}

func (p *pawn) ai_actAlerted() {
	ai := p.ai
	path := CURRENT_MAP.getPathFromTo(p.x, p.y, ai.targetPawn.x, ai.targetPawn.y, false)
	dirx, diry := path.GetNextStepVector()
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}
