package main

import "parcelcreationtool/parcel"

type aiState uint8

const (
	AI_ROAM aiState = iota
	AI_PATROLLING
	AI_SEARCHING
	AI_ALERTED
)

type aiData struct {
	currentState            aiState
	currentStateTimeoutTurn int
	targetPawn              *pawn
	dirx, diry              int // for roaming

	route                *parcel.Route // for patrol
	currentWaypointIndex int

	searchx, searchy int // for search
}

func (p *pawn) ai_checkSituation() {
	switch p.ai.currentState {
	case AI_ROAM, AI_PATROLLING:
		p.ai_checkRoam()
	case AI_SEARCHING:
		p.ai_checkSearching()
	case AI_ALERTED:
		p.ai_checkAlerted()
	default:
		log.AppendMessage("No CHECK func for some ai state!")
	}
	p.ai_checkNoises()
	p.ai_timeoutState()
}

func (p *pawn) ai_act() {
	switch p.ai.currentState {
	case AI_ROAM:
		p.ai_actRoam()
	case AI_PATROLLING:
		p.ai_actPatrolling()
	case AI_SEARCHING:
		p.ai_actSearching()
	case AI_ALERTED:
		p.ai_actAlerted()
	default:
		log.AppendMessage("No ACT func for some ai state!")
	}
}

func (p *pawn) ai_checkNoises() {
	for _, n := range CURRENT_MAP.noises {
		if areCoordinatesInRangeFrom(p.x, p.y, n.x, n.y, n.intensity) {
			if n.suspicious {
				p.ai.currentState = AI_SEARCHING
				p.ai.currentStateTimeoutTurn = CURRENT_TURN + 25*10
				p.ai.searchx, p.ai.searchy = n.x, n.y
			}
		}
	}
}

func (p *pawn) ai_checkRoam() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		return
	}
}

func (p *pawn) ai_actRoam() {
	ai := p.ai
	tries := 0
	for tries < 10 {
		if CURRENT_MAP.isTilePassableAndNotOccupied(p.x+ai.dirx, p.y+ai.diry) && rnd.Rand(25) > 0 {
			break
		} else {
			ai.dirx, ai.diry = rnd.RandomUnitVectorInt()
		}
		tries++
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

func (p *pawn) ai_checkSearching() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		return
	}
}

func (p *pawn) ai_actSearching() {
	ai := p.ai
	path := CURRENT_MAP.getPathFromTo(p.x, p.y, ai.searchx, ai.searchy, false)
	dirx, diry := path.GetNextStepVector()
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}

func (p *pawn) ai_checkAlerted() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		p.ai.searchx, p.ai.searchy = CURRENT_MAP.player.getCoords()
		return
	} else {
		p.ai.currentState = AI_SEARCHING
		p.ai.currentStateTimeoutTurn = CURRENT_TURN+25*10
	}
}

func (p *pawn) ai_actAlerted() {
	ai := p.ai
	var dirx, diry int
	if ai.targetPawn != nil {
		path := CURRENT_MAP.getPathFromTo(p.x, p.y, ai.targetPawn.x, ai.targetPawn.y, false)
		dirx, diry = path.GetNextStepVector()
	} else {
		path := CURRENT_MAP.getPathFromTo(p.x, p.y, ai.searchx, ai.searchy, false)
		dirx, diry = path.GetNextStepVector()
	}
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}
