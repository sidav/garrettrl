package main

type missionTypeCode uint8

const (
	MISSION_STEAL_KEY_ITEMS missionTypeCode = iota
	MISSION_STEAL_MINIMUM_LOOT
)

type Mission struct {
	BriefingText   string
	DebriefingText string

	MissionType            missionTypeCode
	DifficultyChoosingStr  string
	DifficultyLevelsNames  []string
	TargetNumber           []int
	AdditionalGuardsNumber []int
	Rewards                []int
	TotalLoot              []int
}

func (m *Mission) readFromFile(filename string) {

}

//func main() {
//	mis := Mission{
//		BriefingText:
//		"Things got worse lately. After all my struggle, after all the second chances I'm standing right where I belong. " +
//			"The factory I worked in has been shut down. No more job, no more clean hands. Guess my old job still trying " +
//			"to get me. \n " +
//			"I have almost no money left. It will be not so much longer before I will start starving - and " +
//			"then what? Return to being a hired knife? I hate even thinkng of this. I'll better end up picking pockets on" +
//			" a market square. Better being a thief than a murderer. \n " +
//			"And that thought have led me to an idea. A friend of mine mentioned something about his lord selling his mansion. " +
//			"With all the moving, security will be less strict than usual. The goodies will be moved in three days. I don't think the lord " +
//			"needs his goodies too badly, huh? \n " +
//			"Thus said, I have three days for the job. I don't need much - a thousand worth of gold will be enough for some time. " +
//			"I can move out just tonight, or tomorrow, or even later, but the earlier I move out, " +
//			"the more luxuries will be there. And more luxuries means more guards.",
//		MissionType:            MISSION_STEAL_MINIMUM_LOOT,
//		DifficultyChoosingStr:  "When should I go?",
//		DifficultyLevelsNames:  []string{"As late as possible", "Tomorrow night", "Tonight"},
//		TargetNumber:           []int{1000, 1000, 1000},
//		AdditionalGuardsNumber: []int{0, 3, 6},
//		Rewards:                []int{0, 0, 0},
//		TotalLoot:              []int{1000, 1100, 1200},
//	}
//	file, _ := json.MarshalIndent(mis, " ", " ")
//	_ = ioutil.WriteFile("mission.json", file, 0644)
//}
