package main

type missionTypeCode uint8

const (
	MISSION_STEAL_KEY_ITEMS missionTypeCode = iota
	MISSION_STEAL_MINIMUM_LOOT
)

type Mission struct {
	BriefingText           string
	MissionType            missionTypeCode
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
//		"Things got worse lately. I have almost no money left. It will be not so much longer before I will end up\\n" +
//			"picking pockets on a marketplace. I hate even the thinkng of this. Yet, it is not that bad. A friend of mine " +
//			"mentioned something about his lord leaving his mansion for some trip.",
//		MissionType:            MISSION_STEAL_MINIMUM_LOOT,
//		TargetNumber:           []int{1000, 1000, 1000},
//		AdditionalGuardsNumber: []int{0, 3, 6},
//		Rewards:                []int{0, 0, 0},
//		TotalLoot:              []int{1000, 1100, 1200},
//	}
//	file, _ := json.MarshalIndent(mis, " ", " ")
//	_ = ioutil.WriteFile("mission.json", file, 0644)
//}
