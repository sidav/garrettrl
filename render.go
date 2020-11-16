package main

import (
	"fmt"
	cw "github.com/sidav/golibrl/console"
)

var (
	R_VIEWPORT_WIDTH         = 40
	R_VIEWPORT_HEIGHT        = 20
	R_VIEWPORT_CURR_X        = 0
	R_VIEWPORT_CURR_Y        = 0
	R_UI_COLOR_LIGHT         = cw.YELLOW
	R_UI_COLOR_DARK          = cw.DARK_BLUE
	R_UI_COLOR_RUNNING       = cw.RED
	R_LOOTABLE_CABINET_COLOR = cw.DARK_YELLOW
	RENDER_DISABLE_LOS       bool
)

const (
	FogOfWarColor = cw.DARK_GRAY
	darkColor     = cw.BLUE
)

func updateBoundsIfNeccessary(force bool) {
	if cw.WasResized() || force {
		cw, ch := cw.GetConsoleSize()
		R_VIEWPORT_WIDTH = 2 * cw / 3
		R_VIEWPORT_HEIGHT = ch - len(log.Last_msgs) - 1
		//SIDEBAR_X = VIEWPORT_W + 1
		//SIDEBAR_W = cw - VIEWPORT_W - 1
		//SIDEBAR_H = ch - LOG_HEIGHT
		//SIDEBAR_FLOOR_2 = 7  // y-coord right below resources info
		//SIDEBAR_FLOOR_3 = 11 // y-coord right below "floor 2"
	}
}

//func r_areRealCoordsInViewport(x, y int) bool {
//	return x - R_VIEWPORT_CURR_X < R_VIEWPORT_WIDTH && y - R_VIEWPORT_CURR_Y < R_VIEWPORT_HEIGHT
//}

func r_renderUiOutline() {
	w, _ := cw.GetConsoleSize()
	cw.SetBgColor(R_UI_COLOR_DARK)
	if CURRENT_MAP.player.isNotConcealed() {
		cw.SetBgColor(R_UI_COLOR_LIGHT)
	}
	if CURRENT_MAP.player.isRunning {
		cw.SetBgColor(R_UI_COLOR_RUNNING)
	}
	for x := 0; x < w; x++ {
		// cw.PutChar(' ', x, 0)
		cw.PutChar(' ', x, R_VIEWPORT_HEIGHT)
	}
	for y := 0; y < R_VIEWPORT_HEIGHT; y++ {
		cw.PutChar(' ', R_VIEWPORT_WIDTH, y)
	}
}

func r_CoordsToViewport(x, y int) (int, int) {
	vpx, vpy := x-R_VIEWPORT_CURR_X, y-R_VIEWPORT_CURR_Y
	if vpx >= R_VIEWPORT_WIDTH || vpy >= R_VIEWPORT_HEIGHT {
		return -1, -1
	}
	return vpx, vpy
}

func updateViewportCoords(p *pawn) {
	R_VIEWPORT_CURR_X = p.x - R_VIEWPORT_WIDTH/2
	R_VIEWPORT_CURR_Y = p.y - R_VIEWPORT_HEIGHT/2
}

func renderLevel(d *gameMap, flush bool) {
	updateBoundsIfNeccessary(false)
	cw.Clear_console()
	updateViewportCoords(d.player)
	r_renderUiOutline()
	// render level. vpx, vpy are viewport coords, whereas x, y are real coords.
	for x := R_VIEWPORT_CURR_X; x < R_VIEWPORT_CURR_X+R_VIEWPORT_WIDTH; x++ {
		for y := 0; y < R_VIEWPORT_CURR_Y+R_VIEWPORT_HEIGHT; y++ {
			vpx, vpy := r_CoordsToViewport(x, y)
			if !areCoordinatesValid(x, y) {
				continue
			}
			cell := d.tiles[x][y].getAppearance()
			// is seen right now
			if RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
				d.tiles[x][y].wasSeenByPlayer = true
				if d.tiles[x][y].lightLevel > 0 || d.tiles[x][y].isOpaque() {
					renderCcell(cell, vpx, vpy)
				} else {
					renderCcellForceColor(cell, vpx, vpy, darkColor, false)
				}
			} else { // is in fog of war
				if d.tiles[x][y].wasSeenByPlayer {
					renderCcellForceColor(cell, vpx, vpy, FogOfWarColor, false)
				}
			}
			cw.SetBgColor(cw.BLACK)
		}
	}
	//render items
	//for _, item := range d.items {
	//	if RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[item.x][item.y] {
	//		renderItem(item)
	//	}
	//}

	//render pawns
	for _, pawn := range d.pawns {
		if RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[pawn.x][pawn.y] {
			renderPawn(pawn, false)
		}
	}

	// render furniture
	for _, furniture := range d.furnitures {
		if RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[furniture.x][furniture.y] {
			x, y := r_CoordsToViewport(furniture.x, furniture.y)
			if furniture.canBeLooted() {
				renderCcellForceColor(furniture.getStaticData().appearance, x, y, R_LOOTABLE_CABINET_COLOR, false)
			} else {
				renderCcell(furniture.getStaticData().appearance, x, y)
			}
		}
	}

	//render noises
	renderNoisesForPlayer()

	//render player
	furnUnderPlayer := d.getFurnitureAt(d.player.x, d.player.y)
	inverse := furnUnderPlayer != nil && furnUnderPlayer.getStaticData().canBeUsedAsCover
	renderPawn(d.player, inverse)

	renderSidebar()
	renderLog(false)

	if flush {
		cw.Flush_console()
	}
}

func renderPawn(p *pawn, inverse bool) {
	x, y := r_CoordsToViewport(p.x, p.y)
	if p.ai != nil {
		switch p.ai.currentState {
		case AI_ALERTED:
			renderCcellForceChar(p.getStaticData().ccell, x, y, '!')
			return
		case AI_SEARCHING:
			renderCcellForceChar(p.getStaticData().ccell, x, y, '?')
			return
		}

	}
	playerColor := CURRENT_MAP.player.getStaticData().ccell.color
	if p == CURRENT_MAP.player {
		if CURRENT_MAP.tiles[p.x][p.y].lightLevel == 0 {
			playerColor = darkColor
		}
		renderCcellForceColor(p.getStaticData().ccell, x, y, playerColor, inverse)
	} else {
		renderCcell(p.getStaticData().ccell, x, y)
	}
	cw.SetBgColor(cw.BLACK)
}

func renderSidebar() {
	psd := CURRENT_MAP.player.getStaticData()
	p := CURRENT_MAP.player
	cw.SetFgColor(cw.WHITE)
	if p.isRunning {
		cw.PutString(fmt.Sprintf("!! RUNNING !!"), R_VIEWPORT_WIDTH+1, 0)
	} else {
		cw.PutString(fmt.Sprintf(".. sneaking .."), R_VIEWPORT_WIDTH+1, 0)
	}
	cw.PutString(fmt.Sprintf("Health: %d/%d", CURRENT_MAP.player.hp, psd.maxhp), R_VIEWPORT_WIDTH+1, 1)
	cw.PutString(fmt.Sprintf("Loot: %d", CURRENT_MAP.player.inv.gold), R_VIEWPORT_WIDTH+1, 2)
	for i, arrow := range p.inv.arrows {
		cw.PutString(fmt.Sprintf("%s: %d", arrow.name, arrow.amount), R_VIEWPORT_WIDTH+1, 4+i)
	}
}

func renderNoisesForPlayer() {
	log.AppendMessagef("%d noises total", len(CURRENT_MAP.noises))
	for _, n := range CURRENT_MAP.noises {
		if !CURRENT_MAP.currentPlayerVisibilityMap[n.x][n.y] || !n.showOnlyNotSeen {
			// render only those noises in player's vicinity
			if areCoordinatesInRangeFrom(n.x, n.y, CURRENT_MAP.player.x, CURRENT_MAP.player.y, n.intensity) {
				if n.textBubble != "" {
					x, y := r_CoordsToViewport(n.x, n.y)
					x -= len(n.textBubble)/2
					cw.SetColor(cw.BEIGE, cw.DARK_GRAY)
					cw.PutString(n.textBubble, x, y+1)
					cw.SetBgColor(cw.BLACK)
				} else {
					x, y := r_CoordsToViewport(n.x, n.y)
					renderCcell(&n.visual, x, y)
				}
			}
		}
	}
}

func renderLog(flush bool) {
	cw.SetFgColor(cw.RED)
	for i := 0; i < len(log.Last_msgs); i++ {
		cw.PutString(log.Last_msgs[i].Message, 0, R_VIEWPORT_HEIGHT+i+1)
	}
	if flush {
		cw.Flush_console()
	}
}

func renderCcell(cc *consoleCell, x, y int) {
	if cc.inverse {
		cw.SetFgColor(cw.BLACK)
		cw.SetBgColor(cc.color)
	} else {
		cw.SetFgColor(cc.color)
		cw.SetBgColor(cw.BLACK)
	}
	cw.PutChar(cc.appearance, x, y)
}

func renderCcellForceColor(cc *consoleCell, x, y int, color int, forceInverse bool) {
	if cc.inverse || forceInverse {
		cw.SetFgColor(cw.BLACK)
		cw.SetBgColor(color)
	} else {
		cw.SetFgColor(color)
		cw.SetBgColor(cw.BLACK)
	}
	cw.PutChar(cc.appearance, x, y)
}

func renderCcellForceChar(cc *consoleCell, x, y int, char rune) {
	if cc.inverse {
		cw.SetFgColor(cw.BLACK)
		cw.SetBgColor(cc.color)
	} else {
		cw.SetFgColor(cc.color)
		cw.SetBgColor(cw.BLACK)
	}
	cw.PutChar(char, x, y)
}

//func renderLine(char rune, fromx, fromy, tox, toy int, flush, exceptFirstAndLast bool) {
//	line := routines.GetLine(fromx, fromy, tox, toy)
//	SetFgColor(RED)
//	if exceptFirstAndLast {
//		for i := 1; i < len(line)-1; i++ {
//			PutChar(char, line[i].X, line[i].Y)
//		}
//	} else {
//		for i := 0; i < len(line); i++ {
//			PutChar(char, line[i].X, line[i].Y)
//		}
//	}
//	if flush {
//		Flush_console()
//	}
//}
