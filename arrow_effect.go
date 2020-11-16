package main

func applyArrowEffect(arrowName string, x, y int) {
	switch arrowName {
	case "Water arrow":
		furn := CURRENT_MAP.getFurnitureAt(x, y)
		if furn != nil && furn.getCurrentLightLevel() > 0 && furn.getStaticData().isExtinguishable {
			furn.isLit = false
			CURRENT_MAP.createNoise(&noise{
				creator:         nil,
				x:               furn.x,
				y:               furn.y,
				intensity:       5,
				visual:          consoleCell{},
				textBubble:      "",
				suspicious:      true,
				showOnlyNotSeen: true,
			})
		}
	case "Gas arrow":
	case "Explosive arrow":
	case "Noise arrow":
	default:
		log.AppendMessage("Unknown arrow: " + arrowName)
	}
}
