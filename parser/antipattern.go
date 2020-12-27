package parser

type AntipatternTracker struct {
	repetitions map[string]*Repetition
	lastKey     string
	lastMode    string
}

func NewAntipatternTracker() *AntipatternTracker {
	return &AntipatternTracker{
		repetitions: map[string]*Repetition{},
	}
}

type Repetition struct {
	Key   string `header:"pattern"`
	Count int64  `header:"count"`
}

func (t *AntipatternTracker) Track(currentKey, currentMode string) {
	// TODO: double key patterns: dddd instead of ddj
	switch currentKey {
	case "h", "j", "k", "l", "b", "B", "w", "W", "e", "E":
		if t.lastKey == currentKey {
			t.addRepetition(t.lastKey + currentKey)
		}
	case "<cr>":
		if currentMode == InsertMode && t.lastMode != InsertMode {
			switch t.lastKey {
			case "i", "a": // insert/append and enter instead of "o"
				t.addRepetition(t.lastKey + currentKey)
			}
		}
	}

	t.lastKey = currentKey
	t.lastMode = currentMode
}

func (t AntipatternTracker) Repetitions() map[string]*Repetition {
	return t.repetitions
}

func (t *AntipatternTracker) addRepetition(currentKey string) {
	if _, ok := t.repetitions[currentKey]; !ok {
		t.repetitions[currentKey] = &Repetition{
			Key:   currentKey,
			Count: 0,
		}
	}

	t.repetitions[currentKey].Count++
}
