package parser

import "fmt"

const (
	sequenceKey = '\ufffd'
	newKey      = '\b'
)

type SequenceTracker struct {
	sequence        int
	currentSequence string
}

func (s *SequenceTracker) IsActive(r rune) bool {
	if r == sequenceKey {
		s.sequence++

		return true
	} else if r == newKey && s.sequence == 2 {
		s.currentSequence = "<m-%s>"

		return true
	} else if r == 'k' && s.sequence == 1 {
		s.sequence++

		return true
	} else if r == 'b' && s.sequence == 2 {
		s.currentSequence = "<bs>"

		return false
	}

	return false
}

func (s SequenceTracker) Found() bool {
	return s.currentSequence != ""
}

func (s SequenceTracker) CurrentSequence(currentKey string) string {
	if s.currentSequence == "<bs>" {
		return s.currentSequence
	}

	return fmt.Sprintf(s.currentSequence, currentKey)
}

func (s *SequenceTracker) Reset() {
	s.sequence = 0
	s.currentSequence = ""
}
