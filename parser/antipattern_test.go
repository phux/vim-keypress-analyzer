package parser_test

import (
	"testing"

	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/stretchr/testify/require"
)

func TestAntipatternTracker_Track_NormalMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedRepetitions map[string]*parser.Repetition
		name                string
		keys                []string
	}{
		{
			name:                "passing an empty string",
			keys:                []string{""},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
		{
			name:                "passing a single key (tracked)",
			keys:                []string{"w"},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
		{
			name:                "passing a single key (non-tracked)",
			keys:                []string{"q"},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
		{
			name:                "passing two keys (non-tracked)",
			keys:                []string{"q", "q"},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
		{
			name: "passing a tracked key twice should track a repetition",
			keys: []string{"w", "w"},
			expectedRepetitions: map[string]*parser.Repetition{
				"ww": {
					Key:   "ww",
					Count: 1,
				},
			},
		},
		{
			name:                "alternating tracked key, non-tracked key, ...",
			keys:                []string{"w", "q", "w", "q", "w", "q"},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
		{
			name:                "alternating tracked key, different tracked key, ...",
			keys:                []string{"w", "e", "w", "j", "w", "j"},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
		{
			name: "multiple repetitive tracked keys",
			keys: []string{"w", "w", "w", "j", "j"},
			expectedRepetitions: map[string]*parser.Repetition{
				"ww": {
					Key:   "ww",
					Count: 2,
				},
				"jj": {
					Key:   "jj",
					Count: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tracker := parser.NewAntipatternTracker()
			for _, currentKey := range tt.keys {
				tracker.Track(currentKey, parser.NormalMode)
			}
			require.Equal(t, tt.expectedRepetitions, tracker.Repetitions())
		})
	}
}

func TestAntipatternTracker_Track_InsertMode(t *testing.T) {
	t.Parallel()

	type params struct {
		currentMode string
		key         string
	}

	tests := []struct {
		expectedRepetitions map[string]*parser.Repetition
		name                string
		keys                []params
	}{
		{
			name: "insert and enter",
			keys: []params{
				{key: "i", currentMode: parser.NormalMode},
				{key: "<cr>", currentMode: parser.InsertMode},
			},
			expectedRepetitions: map[string]*parser.Repetition{
				"i<cr>": {
					Key:   "i<cr>",
					Count: 1,
				},
			},
		},
		{
			name: "append and enter",
			keys: []params{
				{key: "a", currentMode: parser.NormalMode},
				{key: "<cr>", currentMode: parser.InsertMode},
			},
			expectedRepetitions: map[string]*parser.Repetition{
				"a<cr>": {
					Key:   "a<cr>",
					Count: 1,
				},
			},
		},
		{
			name: "pressing a<cr> already in insert mode",
			keys: []params{
				{key: "a", currentMode: parser.InsertMode},
				{key: "<cr>", currentMode: parser.InsertMode},
			},
			expectedRepetitions: map[string]*parser.Repetition{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tracker := parser.NewAntipatternTracker()
			for _, currentKey := range tt.keys {
				tracker.Track(currentKey.key, currentKey.currentMode)
			}
			require.Equal(t, tt.expectedRepetitions, tracker.Repetitions())
		})
	}
}
