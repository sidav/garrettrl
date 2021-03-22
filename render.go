package main

import (
	"fmt"
	cw "github.com/sidav/golibrl/console"
	"strings"
)

type consoleRenderer struct {
	R_VIEWPORT_WIDTH         int
	R_VIEWPORT_HEIGHT        int
	R_VIEWPORT_CURR_X        int
	R_VIEWPORT_CURR_Y        int
	R_UI_COLOR_LIGHT         int
	R_UI_COLOR_DARK          int
	R_UI_COLOR_RUNNING       int
	R_LOOTABLE_CABINET_COLOR int
	RENDER_DISABLE_LOS       bool
	FogOfWarColor            int
	darkColor                int
	currentSidebarLine       int
}

func (c *consoleRenderer) initDefaults() {
	c.R_VIEWPORT_WIDTH = 40
	c.R_VIEWPORT_HEIGHT = 20
	c.R_VIEWPORT_CURR_X = 0
	c.R_VIEWPORT_CURR_Y = 0
	c.R_UI_COLOR_LIGHT = cw.YELLOW
	c.R_UI_COLOR_DARK = cw.DARK_BLUE
	c.R_UI_COLOR_RUNNING = cw.RED
	c.R_LOOTABLE_CABINET_COLOR = cw.DARK_YELLOW
	c.RENDER_DISABLE_LOS = false
	c.FogOfWarColor = cw.DARK_GRAY
	c.darkColor = cw.BLUE
}

func (c *consoleRenderer) updateBoundsIfNeccessary(force bool) {
	if cw.WasResized() || force {
		cw, ch := cw.GetConsoleSize()
		c.R_VIEWPORT_WIDTH = 2 * cw / 3
		c.R_VIEWPORT_HEIGHT = ch - len(log.Last_msgs) - 1
		//SIDEBAR_X = VIEWPORT_W + 1
		//SIDEBAR_W = cw - VIEWPORT_W - 1
		//SIDEBAR_H = ch - LOG_HEIGHT
		//SIDEBAR_FLOOR_2 = 7  // y-coord right below resources info
		//SIDEBAR_FLOOR_3 = 11 // y-coord right below "floor 2"
	}
}

//func c.R_areRealCoordsInViewport(x, y int) bool {
//	return x - c.R_VIEWPORT_CURR_X < c.R_VIEWPORT_WIDTH && y - c.R_VIEWPORT_CURR_Y < c.R_VIEWPORT_HEIGHT
//}

func (c *consoleRenderer) renderUiOutline() {
	w, _ := cw.GetConsoleSize()
	cw.SetBgColor(c.R_UI_COLOR_DARK)
	if CURRENT_MAP.player.isNotConcealed() {
		cw.SetBgColor(c.R_UI_COLOR_LIGHT)
	}
	if CURRENT_MAP.player.isRunning {
		cw.SetBgColor(c.R_UI_COLOR_RUNNING)
	}
	for x := 0; x < w; x++ {
		// cw.PutChar(' ', x, 0)
		cw.PutChar(' ', x, c.R_VIEWPORT_HEIGHT)
	}
	for y := 0; y < c.R_VIEWPORT_HEIGHT; y++ {
		cw.PutChar(' ', c.R_VIEWPORT_WIDTH, y)
	}
}

func (c *consoleRenderer) coordsToViewport(x, y int) (int, int) {
	vpx, vpy := x-c.R_VIEWPORT_CURR_X, y-c.R_VIEWPORT_CURR_Y
	if vpx >= c.R_VIEWPORT_WIDTH || vpy >= c.R_VIEWPORT_HEIGHT {
		return -1, -1
	}
	return vpx, vpy
}

func (c *consoleRenderer) updateViewportCoords(p *pawn) {
	c.R_VIEWPORT_CURR_X = p.x - c.R_VIEWPORT_WIDTH/2
	c.R_VIEWPORT_CURR_Y = p.y - c.R_VIEWPORT_HEIGHT/2
}

func (c *consoleRenderer) renderGameScreen(flush bool) {
	c.updateBoundsIfNeccessary(false)
	cw.Clear_console()
	c.updateViewportCoords(CURRENT_MAP.player)
	c.renderUiOutline()

	c.renderLevel()

	c.renderBodies()
	//render pawns
	for _, pawn := range CURRENT_MAP.pawns {
		if c.RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[pawn.x][pawn.y] {
			c.renderPawn(pawn, false)
		}
	}

	// render furniture
	c.renderFurnitures()

	//render noises
	c.renderNoisesForPlayer()

	//render player
	furnUnderPlayer := CURRENT_MAP.getFurnitureAt(CURRENT_MAP.player.x, CURRENT_MAP.player.y)
	inverse := furnUnderPlayer != nil && furnUnderPlayer.getStaticData().canBeUsedAsCover
	c.renderPawn(CURRENT_MAP.player, inverse)

	c.currentSidebarLine = 0
	c.renderSidebar()
	c.renderLog(false)

	if flush {
		cw.Flush_console()
	}
}

func (c *consoleRenderer) renderLevel() {
	// render level. vpx, vpy are viewport coords, whereas x, y are real coords.
	for x := c.R_VIEWPORT_CURR_X; x < c.R_VIEWPORT_CURR_X+c.R_VIEWPORT_WIDTH; x++ {
		for y := 0; y < c.R_VIEWPORT_CURR_Y+c.R_VIEWPORT_HEIGHT; y++ {
			vpx, vpy := c.coordsToViewport(x, y)
			if !areCoordinatesValid(x, y) {
				continue
			}
			tile := &CURRENT_MAP.tiles[x][y]
			cell := tile.getAppearance()
			// is seen right now
			if c.RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
				tile.wasSeenByPlayer = true
				if tile.lightLevel > 0 || tile.isOpaque() && !tile.isDoor() {
					c.renderCcell(cell, vpx, vpy)
				} else {
					c.renderCcellForceColor(cell, vpx, vpy, c.darkColor, false)
				}
			} else { // is in fog of war
				if tile.wasSeenByPlayer {
					c.renderCcellForceColor(cell, vpx, vpy, c.FogOfWarColor, false)
				}
			}
			cw.SetBgColor(cw.BLACK)
		}
	}
}

func (c *consoleRenderer) renderPawn(p *pawn, inverse bool) {
	x, y := c.coordsToViewport(p.x, p.y)
	if p.ai != nil {
		switch p.ai.currentState {
		//case AI_ROAM:
		//	c.renderCcellForceChar(p.getStaticData().ccell, x, y, 'R')
		//	return
		//case AI_PATROLLING:
		//	c.renderCcellForceChar(p.getStaticData().ccell, x, y, 'P')
		//	return
		case AI_ALERTED:
			c.renderCcellForceChar(p.getStaticData().ccell, x, y, '!')
			return
		case AI_SEARCHING:
			c.renderCcellForceChar(p.getStaticData().ccell, x, y, '?')
			return
		}

	}

	if p == CURRENT_MAP.player {
		playerColor := CURRENT_MAP.player.getStaticData().ccell.color
		if CURRENT_MAP.tiles[p.x][p.y].lightLevel == 0 {
			playerColor = c.darkColor
		}
		c.renderCcellForceColor(p.getStaticData().ccell, x, y, playerColor, inverse)
	} else {
		c.renderCcell(p.getStaticData().ccell, x, y)
	}
	cw.SetBgColor(cw.BLACK)
}

func (c *consoleRenderer) renderBodies() {
	for _, cBody := range CURRENT_MAP.bodies {
		if c.RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[cBody.x][cBody.y] {
			x, y := c.coordsToViewport(cBody.x, cBody.y)
			color := cBody.pawnOwner.getStaticData().ccell.color
			cw.SetFgColor(color)
			cw.PutChar('%', x, y)
		}
	}
}

func (c *consoleRenderer) renderFurnitures() {
	for _, furniture := range CURRENT_MAP.furnitures {
		if c.RENDER_DISABLE_LOS || CURRENT_MAP.currentPlayerVisibilityMap[furniture.x][furniture.y] {
			x, y := c.coordsToViewport(furniture.x, furniture.y)
			if furniture.canBeLooted() {
				c.renderCcellForceColor(furniture.getStaticData().appearance, x, y, c.R_LOOTABLE_CABINET_COLOR, false)
			} else if furniture.getCurrentLightLevel() > 0 {
				c.renderCcellForceColor(furniture.getStaticData().appearance, x, y, c.R_UI_COLOR_LIGHT, false)
			} else {
				c.renderCcell(furniture.getStaticData().appearance, x, y)
			}
		}
	}
}

func (c *consoleRenderer) renderSidebar() {
	psd := CURRENT_MAP.player.getStaticData()
	p := CURRENT_MAP.player
	cw.SetFgColor(cw.WHITE)
	if p.isRunning {
		c.addCenteredSidebarLine("!! RUNNING !!")
	} else {
		c.addCenteredSidebarLine(".. sneaking ..")
	}
	c.addAlignedSidebarLine("Health:", fmt.Sprintf("%d/%d", CURRENT_MAP.player.hp, psd.maxhp), false)
	c.addAlignedSidebarLine("Loot: ", fmt.Sprintf("%d", CURRENT_MAP.player.inv.gold), false)
	c.addCenteredSidebarLine( "Arrows:")
	for i, arrow := range p.inv.arrows {
		inverseLine := currPlayerController.currentSelectedArrowIndex == i
		c.addAlignedSidebarLine(fmt.Sprintf("%s:", strings.Replace(arrow.name, " arrow", "", 1)), fmt.Sprintf("%d", arrow.amount), inverseLine)
	}
	if len(p.inv.targetItems) > 0 {
		cw.SetColor(cw.DARK_YELLOW, cw.BLACK)
		c.addCenteredSidebarLine(fmt.Sprintf("Target items (%d/%d):", len(p.inv.targetItems), currMission.TargetNumber[currDifficultyNumber]), )
		cw.SetColor(cw.WHITE, cw.BLACK)
		for _, itemName := range p.inv.targetItems {
			c.addAlignedSidebarLine(" "+itemName, "", false)
		}
	}
}

func (c *consoleRenderer) addCenteredSidebarLine(line string) {
	wid, _ := cw.GetConsoleSize()
	cw.SetFgColor(cw.WHITE)
	if c.currentSidebarLine % 2 == 1 {
		cw.SetFgColor(cw.BEIGE)
	}
	centeredCoords := c.R_VIEWPORT_WIDTH + (wid - c.R_VIEWPORT_WIDTH)/2 - len(line)/2
	cw.PutString(line, centeredCoords, c.currentSidebarLine)
	c.currentSidebarLine++
}

func (c *consoleRenderer) addAlignedSidebarLine(leftAligned, rightAligned string, inversion bool) {
	const OFFSET = 2
	wid, _ := cw.GetConsoleSize()
	if inversion {
		cw.SetBgColor(cw.WHITE)
		cw.SetFgColor(cw.BLACK)
		if c.currentSidebarLine%2 == 1 {
			cw.SetBgColor(cw.BEIGE)
		}
	} else {
		cw.SetFgColor(cw.WHITE)
		if c.currentSidebarLine%2 == 1 {
			cw.SetFgColor(cw.BEIGE)
		}
	}
	// put left aligned
	spaces := strings.Repeat(" ", wid - c.R_VIEWPORT_WIDTH - len(leftAligned) - OFFSET*2)
	cw.PutString(leftAligned + spaces, c.R_VIEWPORT_WIDTH+OFFSET, c.currentSidebarLine)
	// put right aligned
	cw.PutString(rightAligned, wid-len(rightAligned)-OFFSET, c.currentSidebarLine)
	c.currentSidebarLine++
	cw.SetBgColor(cw.BLACK)
	cw.SetFgColor(cw.WHITE)
}

func (c *consoleRenderer) renderDamageFlash() {
	cw.SetBgColor(cw.DARK_RED)
	w, h := cw.GetConsoleSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cw.PutChar(' ', x, y)
		}
	}
	cw.Flush_console()
	cw.SetBgColor(cw.BLACK)
}

func (c *consoleRenderer) renderNoisesForPlayer() {
	// log.AppendMessagef("%d noises total", len(CURRENT_MAP.noises))
	for _, n := range CURRENT_MAP.noises {
		if !CURRENT_MAP.currentPlayerVisibilityMap[n.x][n.y] || !n.showOnlyNotSeen {
			// render only those noises in player's vicinity
			if areCoordinatesInRangeFrom(n.x, n.y, CURRENT_MAP.player.x, CURRENT_MAP.player.y, n.intensity) {
				if n.textBubble != "" {
					x, y := c.coordsToViewport(n.x, n.y)
					if n.creator != nil {
						x, y = c.coordsToViewport(n.creator.getCoords())
					}
					if x == -1 && y == -1 {
						continue
					}
					x -= len(n.textBubble) / 2
					if n.suspicious {
						if n.creator != nil {
							cw.SetColor(cw.RED, cw.DARK_GRAY)
						} else {
							cw.SetColor(cw.YELLOW, cw.DARK_GRAY)
						}
					} else {
						cw.SetColor(cw.BEIGE, cw.DARK_GRAY)
					}
					cw.PutString(n.textBubble, x, y+1)
					cw.SetBgColor(cw.BLACK)
				} else {
					x, y := c.coordsToViewport(n.x, n.y)
					c.renderCcell(&n.visual, x, y)
				}
			}
		}
	}
}

func (c *consoleRenderer) renderCursor(cx, cy int, flush bool) {
	cw.SetFgColor(cw.YELLOW)
	vx, vy := c.coordsToViewport(cx, cy)
	cw.PutChar('X', vx, vy)
	if flush {
		cw.Flush_console()
	}
}

func (c *consoleRenderer) renderLog(flush bool) {
	cw.SetFgColor(cw.RED)
	for i := 0; i < len(log.Last_msgs); i++ {
		msg := log.Last_msgs[i]
		stringToPut := msg.Message
		if msg.Count > 1 {
			stringToPut += fmt.Sprintf(" (x%d)", msg.Count)
		}
		cw.PutString(stringToPut, 0, c.R_VIEWPORT_HEIGHT+i+1)
	}
	if flush {
		cw.Flush_console()
	}
}

func (c *consoleRenderer) renderCcell(cc *consoleCell, x, y int) {
	if cc.inverse {
		cw.SetFgColor(cw.BLACK)
		cw.SetBgColor(cc.color)
	} else {
		cw.SetFgColor(cc.color)
		cw.SetBgColor(cw.BLACK)
	}
	cw.PutChar(cc.getAppearance(), x, y)
}

func (c *consoleRenderer) renderCcellForceColor(cc *consoleCell, x, y int, color int, forceInverse bool) {
	if cc.inverse || forceInverse {
		cw.SetFgColor(cw.BLACK)
		cw.SetBgColor(color)
	} else {
		cw.SetFgColor(color)
		cw.SetBgColor(cw.BLACK)
	}
	cw.PutChar(cc.getAppearance(), x, y)
}

func (c *consoleRenderer) renderCcellForceChar(cc *consoleCell, x, y int, char rune) {
	if cc.inverse {
		cw.SetFgColor(cw.BLACK)
		cw.SetBgColor(cc.color)
	} else {
		cw.SetFgColor(cc.color)
		cw.SetBgColor(cw.BLACK)
	}
	cw.PutChar(char, x, y)
}

// puts wrapped text in rectangle.
func (c *consoleRenderer) putTextInRect(text string, x, y, w int) {
	if w == 0 {
		w, _ = cw.GetConsoleSize()
	}
	cx, cy := x, y
	splittedText := strings.Split(text, " ")
	for _, word := range splittedText {
		if cx-x+len(word) > w || word == "\\n" || word == "\n" {
			cx = 0
			cy += 1
		}
		if word != "\\n" && word != "\n" {
			cw.PutString(word, cx, cy)
			cx += len(word) + 1
		}
	}
}

func (r *consoleRenderer) renderBuyMenu(bm *buyMenu) {
	cw.Clear_console()
	cw.SetColor(cw.BLACK, cw.DARK_GREEN)
	cw.PutString("     BLACK MARKET     ", 0, 0)
	cw.PutString(fmt.Sprintf("I have %d gold.", bm.currentGold), 0, 1)
	for i := range bm.itemsNames {
		if i == bm.cursorPosition {
			cw.SetColor(cw.BLACK, cw.BEIGE)
		} else {
			cw.SetColor(cw.BEIGE, cw.BLACK)
		}
		cw.PutString(fmt.Sprintf("%s (%d gold)    < %d >", bm.itemsNames[i], bm.itemsCosts[i], bm.itemsBought[i]),
			0, i+2)
	}
	cw.Flush_console()
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
