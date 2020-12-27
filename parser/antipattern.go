package parser

type AntipatternTracker struct {
	antipatterns map[string]*Antipattern
	lastKey      string
	lastMode     string
}

func NewAntipatternTracker() *AntipatternTracker {
	return &AntipatternTracker{
		antipatterns: map[string]*Antipattern{},
	}
}

type Antipattern struct {
	Key   string `header:"pattern"`
	Count int64  `header:"count"`
}

func (t *AntipatternTracker) Track(currentKey, currentMode string) {
	// TODO: double key patterns: dddd instead of dj
	switch currentKey {
	case "h", "j", "k", "l", "b", "B", "w", "W", "e", "E":
		if currentMode != NormalMode {
			break
		}

		if t.lastKey == currentKey {
			t.addAntipatternOccurrence(t.lastKey + currentKey)
		}
	case "i", "a", "o", "O":
		if t.lastKey == "h" && currentKey == "a" {
			t.addAntipatternOccurrence(t.lastKey + currentKey)
		}

		if t.lastKey == "j" && currentKey == "O" {
			t.addAntipatternOccurrence(t.lastKey + currentKey)
		}

		if t.lastKey == "k" && currentKey == "o" {
			t.addAntipatternOccurrence(t.lastKey + currentKey)
		}

		if t.lastKey == "l" && currentKey == "i" {
			t.addAntipatternOccurrence(t.lastKey + currentKey)
		}
	case "<cr>":
		if currentMode == InsertMode && t.lastMode != InsertMode {
			switch t.lastKey {
			case "i", "a": // insert/append and enter instead of "o"
				t.addAntipatternOccurrence(t.lastKey + currentKey)
			}
		}
	}

	t.lastKey = currentKey
	t.lastMode = currentMode
}

func (t AntipatternTracker) Antipatterns() map[string]*Antipattern {
	return t.antipatterns
}

func (t *AntipatternTracker) addAntipatternOccurrence(currentKey string) {
	if _, ok := t.antipatterns[currentKey]; !ok {
		t.antipatterns[currentKey] = &Antipattern{
			Key:   currentKey,
			Count: 0,
		}
	}

	t.antipatterns[currentKey].Count++
}
