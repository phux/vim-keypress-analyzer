package parser

import (
	"io"
	"io/ioutil"
	"strconv"
	"unicode"

	"github.com/phux/vimkeypressanalyzer/tree"
	"github.com/pkg/errors"
)

const (
	NormalMode            = "normal"
	InsertMode            = "insert"
	CommandMode           = "command"
	VisualMode            = "visual"
	CharEsc               = '\x1b'
	CharReadableEsc       = "<esc>"
	CharEnter             = '\x0D'
	CharReadableEnter     = "<cr>"
	CharSpace             = " "
	CharReadableSpace     = "<space>"
	CharCC                = '\x03'
	CharReadableCC        = "<c-c>"
	CharCV                = '\x16'
	CharReadableCV        = "<c-v>"
	CharTab               = '\x09'
	CharReadableTab       = "<tab>"
	CharBackspace         = '\b'
	CharReadableBackspace = "<bs>"
)

var controlCodeToHumanReadable = map[rune]string{
	CharEsc:       CharReadableEsc,
	CharEnter:     CharReadableEnter,
	CharCC:        CharReadableCC,
	CharCV:        CharReadableCV,
	CharTab:       CharReadableTab,
	CharBackspace: CharReadableBackspace,
	'\x00':        "^@",
	'\x01':        "^A",
	'\x02':        "^B",
	'\x04':        "^D",
	'\x05':        "^E",
	'\x06':        "^F",
	'\x07':        "^G",
	'\x0a':        "^J",
	'\x0b':        "^K",
	'\x0c':        "^L",
	'\x0e':        "^N",
	'\x0f':        "^O",
	'\x10':        "^P",
	'\x11':        "^Q",
	'\x12':        "^R",
	'\x13':        "^S",
	'\x14':        "^T",
	'\x15':        "^U",
	'\x17':        "^W",
	'\x18':        "^X",
	'\x19':        "^Y",
	'\x1a':        "^Z",
	'\x1c':        "^\\",
	'\x1d':        "^]",
}

type Parser struct {
	currentMode        string
	previousKey        string
	enableAntipatterns bool
	isSearchActive     bool
	isMotionActive     bool
}

func NewParser(enableAntipatterns bool) *Parser {
	return &Parser{
		enableAntipatterns: enableAntipatterns,
		currentMode:        NormalMode,
	}
}

func (p *Parser) Parse(r io.Reader, excludeModes []string) (*Result, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read input")
	}

	input := string(raw)
	result := NewResult()
	keymapNode := tree.NewNode("")
	modeCountNode := tree.NewNode("")
	sequenceTracker := &SequenceTracker{}

	maxAllowedRepeats := int64(2)
	antipatternTracker := NewAntipatternTracker(maxAllowedRepeats)

NextKey:
	for _, r := range input {
		if sequenceTracker.IsActive(r) {
			continue
		}

		var currentKey string
		if sequenceTracker.Found() {
			currentKey = sequenceTracker.CurrentSequence(toReadable(r))
		} else {
			currentKey = toReadable(r)
		}

		sequenceTracker.Reset()

		currentMode := p.currentMode

		modeCountNode.AddOrIncrementChild(currentMode)
		p.setNewMode(currentKey)

		p.previousKey = currentKey

		for _, excludeMode := range excludeModes {
			if currentMode == excludeMode {
				continue NextKey
			}
		}

		keymapNode.AddOrIncrementChild(currentKey)

		if p.enableAntipatterns {
			antipatternTracker.Track(currentKey, currentMode)
		}
	}

	result.KeyMap = keymapNode
	result.ModeCount = modeCountNode
	result.Antipatterns = antipatternTracker.Antipatterns()

	return result, nil
}

func (p *Parser) setNewMode(currentKey string) {
	if currentKey == CharReadableEsc || currentKey == CharReadableCC {
		p.currentMode = NormalMode
		p.isSearchActive = false

		return
	}

	if currentKey == CharReadableEnter {
		if p.currentMode == CommandMode {
			p.currentMode = NormalMode
		}

		p.isSearchActive = false
		p.isMotionActive = false

		return
	}

	if p.isSearchActive {
		return
	}

	if p.isMotionActive {
		p.isMotionActive = false

		return
	}

	switch currentKey {
	case "/", "?":
		if p.currentMode == NormalMode || p.currentMode == VisualMode {
			p.isSearchActive = true
		}

		return
	case "f", "F", "t", "T":
		if p.currentMode == NormalMode || p.currentMode == VisualMode {
			p.isMotionActive = true

			return
		}
	case "i", "I", "a", "A", "o", "O", "C":
		if p.currentMode == NormalMode {
			p.currentMode = InsertMode
		}
	case "c":
		if p.currentMode == NormalMode && p.previousKey == "c" {
			p.currentMode = InsertMode
		}
	case ":":
		if p.currentMode == NormalMode || p.currentMode == VisualMode {
			p.currentMode = CommandMode
		}
	case "v", "V", CharReadableCV:
		switch p.currentMode {
		case NormalMode:
			p.currentMode = VisualMode
		case VisualMode:
			p.currentMode = NormalMode
		}
	}
}

func toReadable(r rune) string {
	str := string(r)
	if !unicode.IsControl(r) {
		if str == CharSpace {
			return CharReadableSpace
		}

		return str
	}

	tmp := strconv.QuoteRune(r)
	if val, ok := controlCodeToHumanReadable[r]; ok {
		tmp = val
	}

	return tmp
}
