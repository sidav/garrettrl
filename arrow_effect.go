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
