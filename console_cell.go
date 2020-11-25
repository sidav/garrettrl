package main

type consoleCell struct {
	appearance    rune
	altAppearance rune
	color         int
	inverse       bool
}

func (c *consoleCell) getAppearance() rune {
	if USE_ALT_RUNES && c.altAppearance > 0 {
		return c.altAppearance
	}
	return c.appearance
}
