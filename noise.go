package main

type noise struct {
	creator                     *pawn
	x, y                        int
	intensity                   int
	visual                      consoleCell
	textBubble                  string
	turnCreatedAt, duration     int
	suspicious, showOnlyNotSeen bool
}
