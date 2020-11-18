package main

func applyArrowEffect(arrowName string, x, y int) {
	switch arrowName {
	case "Water arrow":
		furn := CURRENT_MAP.getFurnitureAt(x, y)
		if furn != nil && furn.getCurrentLightLevel() > 0 && furn.getStaticData().isExtinguishable {
			furn.isLit = false
		}
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       5,
			visual:          consoleCell{},
			textBubble:      "Splash!",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	case "Gas arrow":
		for i := x-1; i <= x+1; i++ {
			for j := y-1; j <= y+1; j++ {
				pawnAt := CURRENT_MAP.getPawnAt(i, j)
				if pawnAt != nil && pawnAt != CURRENT_MAP.player {
					newBody := pawnAt.createBody(rnd.RandInRange(10, 15) * 10)
					CURRENT_MAP.bodies = append(CURRENT_MAP.bodies, newBody)
					CURRENT_MAP.removePawn(pawnAt)
				}
			}
		}
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       5,
			duration:        10,
			visual:          consoleCell{},
			textBubble:      "* Fssss *",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	case "Explosive arrow":
	case "Noise arrow":
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       15,
			duration:        50,
			visual:          consoleCell{},
			textBubble:      "*SCREECH*",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	default:
		log.AppendMessage("Unknown arrow: " + arrowName)
	}
}
