package main

type missionTypeCode uint8

const (
	MISSION_STEAL_KEY_ITEMS missionTypeCode = iota
	MISSION_STEAL_MINIMUM_LOOT
)

type Mission struct {
	MissionBriefingText string
	MissionType         missionTypeCode
	MissionRewards      []int
	MissionTotalLoot    int
}

func (m *Mission) readFromFile(filename string) {

}
