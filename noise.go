package main

type noise struct {
	x, y      int
	intensity int
	visual    consoleCell
	suspicious, showOnlyNotSeen bool
}
