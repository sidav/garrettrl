package main

import (
	"fmt"
	cw "github.com/sidav/golibrl/console"
)

var (
	R_VIEWPORT_WIDTH   = 40
	R_VIEWPORT_HEIGHT  = 20
	R_VIEWPORT_CURR_X  = 0
	R_VIEWPORT_CURR_Y  = 0
	R_UI_COLOR_LIGHT = cw.DARK_YELLOW
	R_UI_COLOR_DARK = cw.DARK_BLUE
	R_UI_COLOR_RUNNING = cw.RED
	RENDER_DISABLE_LOS bool
)

const (
	FogOfWarColor = cw.DARK_GRAY
	darkColor = cw.BLUE
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
	if CURRENT_MAP.player.isInLight() {
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
					renderCcellForceColor(cell, vpx, vpy, darkColor)
				}
			} else { // is in fog of war
				if d.tiles[x][y].wasSeenByPlayer {
					renderCcellForceColor(cell, vpx, vpy, FogOfWarColor)
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
			renderCcell(furniture.getStaticData().appearance, x, y)
		}
	}

	//render noises
	renderNoisesForPlayer()

	//render player
	renderPawn(d.player, false)

	renderSidebar()
	renderLog(false)

	if flush {
		cw.Flush_console()
	}
}

//func renderProjectile(p *projectile) {
//	SetColor(RED, BLACK)
//	x, y := r_CoordsToViewport(p.x, p.y)
//	PutChar('*', x, y)
//}

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
	if p == CURRENT_MAP.player && CURRENT_MAP.tiles[p.x][p.y].lightLevel == 0 {
		renderCcellForceColor(p.getStaticData().ccell, x, y, darkColor)
		return
	}
	renderCcell(p.getStaticData().ccell, x, y)
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
}

func renderNoisesForPlayer() {
	log.AppendMessagef("%d noises total", len(CURRENT_MAP.noises))
	for _, n := range CURRENT_MAP.noises {
		if !CURRENT_MAP.currentPlayerVisibilityMap[n.x][n.y] || !n.showOnlyNotSeen {
			// render only those noises in player's vicinity
			if areCoordinatesInRangeFrom(n.x, n.y, CURRENT_MAP.player.x, CURRENT_MAP.player.y, n.intensity) {
				x, y := r_CoordsToViewport(n.x, n.y)
				renderCcell(&n.visual, x, y)
			}
		}
	}
}

//func renderItem(i *i_item) {
//	SetFgColor(i.ccell.color)
//	x, y := r_CoordsToViewport(i.x, i.y)
//	PutChar(i.ccell.appearance, x, y)
//}

//func renderBullets(currCoords []*routines.Vector, currDirs []*routines.Vector, d *gameMap) {
//	renderLevel(d, false)
//
//	for i:=0; i<len(currCoords); i++ {
//		currx, curry := currCoords[i].GetRoundedCoords()
//		tox, toy := currDirs[i].GetRoundedCoords()
//		SetFgColor(YELLOW)
//		bulletRune := '*'
//		if !d.isPawnPresent(currx, curry) && !d.isTileOpaque(currx, curry) {
//			bulletRune = getTargetingChar(tox, toy)
//		}
//		x, y := r_CoordsToViewport(currx, curry)
//		PutChar(bulletRune, x, y)
//	}
//	Flush_console()
//	time.Sleep(35 * time.Millisecond)
//}

//
// UI-related stuff below
//

//func renderPlayerStats(d *gameMap) {
//	player := d.player
//	pinv := player.inventory
//	statusbarsWidth := 80 - R_VIEWPORT_WIDTH - 3
//
//	hpPercent := player.getHpPercent()
//	var hpColor int
//	switch {
//	case hpPercent < 33:
//		hpColor = RED
//		break
//	case hpPercent < 66:
//		hpColor = YELLOW
//		break
//	default:
//		hpColor = DARK_GREEN
//		break
//	}
//	SetFgColor(hpColor)
//
//	renderStatusbar(fmt.Sprintf("HP: (%d/%d)", player.hp, player.maxhp), player.hp, player.maxhp,
//		R_VIEWPORT_WIDTH+1, 0, statusbarsWidth, hpColor)
//
//	if player.wearedArmor == nil {
//		SetFgColor(BEIGE)
//		PutString("No armor", R_VIEWPORT_WIDTH+1, 1)
//	} else {
//		SetFgColor(player.wearedArmor.ccell.color)
//		renderStatusbar(fmt.Sprintf("ARMOR: (%d/%d)", player.wearedArmor.armorData.currArmor, player.wearedArmor.armorData.maxArmor),
//			player.wearedArmor.armorData.currArmor, player.wearedArmor.armorData.maxArmor, R_VIEWPORT_WIDTH+1, 1, statusbarsWidth, player.wearedArmor.ccell.color)
//	}
//
//	SetFgColor(BEIGE)
//	if player.weaponInHands != nil {
//		renderStatusbar(fmt.Sprintf("%s (%d/%d)", player.weaponInHands.name, player.weaponInHands.weaponData.ammo,
//			player.weaponInHands.weaponData.maxammo), player.weaponInHands.weaponData.ammo,
//			player.weaponInHands.weaponData.maxammo, R_VIEWPORT_WIDTH+1, 2, statusbarsWidth, DARK_YELLOW)
//	} else {
//		PutString("Barehanded", R_VIEWPORT_WIDTH+1, 2)
//	}
//
//	SetFgColor(BEIGE)
//	PutString(fmt.Sprintf("INV: %d/%d", len(pinv.items), pinv.maxItems), R_VIEWPORT_WIDTH+1, 3)
//
//	SetColor(BEIGE, BLACK)
//	ammoLine := fmt.Sprintf("BULL:%d/%d", pinv.ammo[AMMO_BULL], pinv.maxammo[AMMO_BULL])
//	PutString(ammoLine, R_VIEWPORT_WIDTH+1, 4)
//	ammoLine = fmt.Sprintf("SHLL:%d/%d", pinv.ammo[AMMO_SHEL], pinv.maxammo[AMMO_SHEL])
//	PutString(ammoLine, R_VIEWPORT_WIDTH+1, 5)
//	ammoLine = fmt.Sprintf("RCKT:%d/%d", pinv.ammo[AMMO_RCKT], pinv.maxammo[AMMO_RCKT])
//	PutString(ammoLine, R_VIEWPORT_WIDTH+1, 6)
//	ammoLine = fmt.Sprintf("CELL:%d/%d", pinv.ammo[AMMO_CELL], pinv.maxammo[AMMO_CELL])
//	PutString(ammoLine, R_VIEWPORT_WIDTH+1, 7)
//
//	timeline := fmt.Sprintf("TIME: %d.%d (%d.%d)", CURRENT_TURN/10, CURRENT_TURN%10,
//		player.playerData.lastSpentTimeAmount/10, player.playerData.lastSpentTimeAmount%10)
//	PutString(timeline, R_VIEWPORT_WIDTH+1, 9)
//
//	remEnemiesLine := fmt.Sprintf("ENEMIES LEFT: %d", len(d.pawns))
//	PutString(remEnemiesLine, R_VIEWPORT_WIDTH+1, 10)
//}
//
//func renderTargetingLine(fromx, fromy, tox, toy int, flush bool, d *gameMap) {
//	renderLevel(d, false)
//	line := routines.GetLine(fromx, fromy, tox, toy)
//	char := '?'
//	if len(line) > 1 {
//		dirVector := routines.CreateVectorByStartAndEndInt(fromx, fromy, tox, toy)
//		dirVector.TransformIntoUnitVector()
//		dirx, diry := dirVector.GetRoundedCoords()
//		char = getTargetingChar(dirx, diry)
//	}
//	if fromx == tox && fromy == toy {
//		renderPawn(d.player, true)
//	}
//	for i := 1; i < len(line); i++ {
//		x, y := line[i].X, line[i].Y
//		if d.isPawnPresent(x, y) {
//			renderPawn(d.getPawnAt(x, y), true)
//		} else {
//			SetFgColor(YELLOW)
//			if i == len(line)-1 {
//				char = 'X'
//			}
//			viewx, viewy := r_CoordsToViewport(line[i].X, line[i].Y)
//			PutChar(char, viewx, viewy)
//		}
//	}
//	if flush {
//		Flush_console()
//	}
//}
//
//func renderStatusbar(name string, curvalue, maxvalue, x, y, width, barColor int) {
//	barTitle := name
//	PutString(barTitle, x, y)
//	barWidth := width - len(name)
//	filledCells := barWidth * curvalue / maxvalue
//	barStartX := x + len(barTitle) + 1
//	for i := 0; i < barWidth; i++ {
//		if i < filledCells {
//			SetFgColor(barColor)
//			PutChar('=', i+barStartX, y)
//		} else {
//			SetFgColor(DARK_BLUE)
//			PutChar('-', i+barStartX, y)
//		}
//	}
//}
//
//func getTargetingChar(x, y int) rune {
//	if abs(x) > 1 {
//		x /= abs(x)
//	}
//	if abs(y) > 1 {
//		y /= abs(y)
//	}
//	if x == 0 {
//		return '|'
//	}
//	if y == 0 {
//		return '-'
//	}
//	if x*y == 1 {
//		return '\\'
//	}
//	if x*y == -1 {
//		return '/'
//	}
//	return '?'
//}
//
//func abs(i int) int {
//	if i < 0 {
//		return -i
//	}
//	return i
//}

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

func renderCcellForceColor(cc *consoleCell, x, y int, color int) {
	if cc.inverse {
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
