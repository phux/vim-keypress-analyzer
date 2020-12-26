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
	CharBackspace         = '\x08'
	CharReadableBackspace = "<bs"
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
	currentMode string
	previousKey string
}

func NewParser() *Parser {
	return &Parser{
		currentMode: NormalMode,
	}
}

func (p *Parser) Parse(r io.Reader) (*Result, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read input")
	}

	input := string(raw)
	result := NewResult()
	rootNode := tree.NewNode("")

	for _, r := range input {
		currentKey := toReadable(r)
		currentMode := p.currentMode
		result.IncrModeCount(currentMode)
		p.setNewMode(currentKey)

		// either we were not in insert mode
		// OR
		// we escaped from insert mode with the current key
		if currentMode != InsertMode {
			rootNode.AddOrIncrementChild(currentKey)
			p.previousKey = currentKey
		}
	}

	result.KeyMap = rootNode

	return result, nil
}

func (p *Parser) setNewMode(currentKey string) {
	switch currentKey {
	case CharReadableEsc:
		if p.currentMode != NormalMode {
			p.currentMode = NormalMode
		}
	case CharReadableEnter:
		if p.currentMode == CommandMode {
			p.currentMode = NormalMode
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
