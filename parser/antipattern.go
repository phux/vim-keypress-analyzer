package parser

import "strings"

type AntipatternTracker struct {
	antipatterns        map[string]*Antipattern
	lastKey             string
	lastMode            string
	consecutiveKeyCount int64
	maxAllowedRepeats   int64
}

func NewAntipatternTracker(maxAllowedRepeats int64) *AntipatternTracker {
	return &AntipatternTracker{
		antipatterns:      map[string]*Antipattern{},
		maxAllowedRepeats: maxAllowedRepeats,
	}
}

type Antipattern struct {
	Key               string  `header:"pattern"`
	Count             int64   `header:"count"`
	TotalKeypresses   int64   `header:"total key presses"`
	AverageKeypresses float64 `header:"avg keys per occurrence"`
}

func (t *AntipatternTracker) Track(currentKey, currentMode string) {
	// TODO: multi key patterns: dwdwdw instead of d3w
	if currentMode != NormalMode && currentMode != VisualMode {
		t.consecutiveKeyCount = 0
	} else {
		t.checkNormalMode(currentKey)
	}

	if currentMode == InsertMode && t.lastMode != InsertMode && currentKey == "<cr>" {
		switch t.lastKey {
		case "I", "A": // insert/append and enter instead of "o"/"O"
			patternName := t.lastKey + currentKey
			t.addAntipatternOccurrence(patternName)
			t.antipatterns[patternName].TotalKeypresses += 2
		}
	}

	t.lastKey = currentKey
	t.lastMode = currentMode
}

func (t *AntipatternTracker) checkConsecutiveKeys(currentKey string, maxAllowedRepeats int64) {
	if t.lastKey == currentKey {
		patternName := strings.Repeat(currentKey, int(maxAllowedRepeats)+1) + "+"
		t.consecutiveKeyCount++

		if t.consecutiveKeyCount == maxAllowedRepeats {
			t.addAntipatternOccurrence(patternName)
			t.antipatterns[patternName].TotalKeypresses += maxAllowedRepeats + 1
		}

		if t.consecutiveKeyCount > maxAllowedRepeats {
			t.antipatterns[patternName].TotalKeypresses++
		}
	} else {
		t.consecutiveKeyCount = 0
	}
}

func (t *AntipatternTracker) checkNormalMode(currentKey string) {
	switch currentKey {
	case "h", "j", "k", "l", "b", "B", "w", "W", "e", "E", "x", "X":
		t.checkConsecutiveKeys(currentKey, t.maxAllowedRepeats)
	case "d":
		t.checkConsecutiveKeys(currentKey, 3)
	case "i", "a", "o", "O":
		patternName := t.lastKey + currentKey
		if t.lastKey == "h" && currentKey == "a" {
			t.addAntipatternOccurrence(patternName)
			t.antipatterns[patternName].TotalKeypresses += 2
		}

		if t.lastKey == "j" && currentKey == "O" {
			t.addAntipatternOccurrence(patternName)
			t.antipatterns[patternName].TotalKeypresses += 2
		}

		if t.lastKey == "k" && currentKey == "o" {
			t.addAntipatternOccurrence(patternName)
			t.antipatterns[patternName].TotalKeypresses += 2
		}

		if t.lastKey == "l" && currentKey == "i" {
			t.addAntipatternOccurrence(patternName)
			t.antipatterns[patternName].TotalKeypresses += 2
		}
	}
}

func (t AntipatternTracker) Antipatterns() map[string]*Antipattern {
	return t.antipatterns
}

func (t *AntipatternTracker) addAntipatternOccurrence(currentKey string) {
	if _, ok := t.antipatterns[currentKey]; !ok {
		t.antipatterns[currentKey] = &Antipattern{
			Key:             currentKey,
			Count:           0,
			TotalKeypresses: 0,
		}
	}

	t.antipatterns[currentKey].Count++
}
