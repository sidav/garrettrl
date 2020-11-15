package main

type noise struct {
	x, y                        int
	intensity                   int
	visual                      consoleCell
	turnCreatedAt, duration     int
	suspicious, showOnlyNotSeen bool
}
