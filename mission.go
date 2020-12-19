package main

type missionTypeCode uint8

const (
	MISSION_STEAL_TARGET_ITEMS missionTypeCode = iota
	MISSION_STEAL_MINIMUM_LOOT
)

type Mission struct {
	BriefingText   string
	DebriefingText string

	MissionType            missionTypeCode
	DifficultyChoosingStr  string
	DifficultyLevelsNames  []string
	TargetNumber           []int
	TargetItemsNames       []string
	AdditionalGuardsNumber []int
	Rewards                []int
	TotalLoot              []int
}

func (m *Mission) readFromFile(filename string) {

}

//func main() {
//	mis := Mission{
//		BriefingText:
//		"Well, in this very,very long briefing I got to know that I have to steal lord Manderley's family jewels. All three of them.",
//		MissionType:            MISSION_STEAL_MINIMUM_LOOT,
//		DifficultyChoosingStr:  "When should I go?",
//		DifficultyLevelsNames:  []string{"As late as possible", "Tomorrow night", "Tonight"},
//		TargetItemsNames: 			[]string{"Family Diamond", "Family Emerald", "Family Ruby"},
//		AdditionalGuardsNumber: []int{0, 3, 6},
//		Rewards:                []int{500, 750, 800},
//		TotalLoot:              []int{300, 400, 500},
//	}
//	file, _ := json.MarshalIndent(mis, " ", " ")
//	_ = ioutil.WriteFile("mission.json", file, 0644)
//}
