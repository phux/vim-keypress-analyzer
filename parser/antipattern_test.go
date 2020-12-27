package parser_test

import (
	"testing"

	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/stretchr/testify/require"
)

func TestAntipatternTracker_Track_NormalMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedAntipatterns map[string]*parser.Antipattern
		name                 string
		keys                 []string
	}{
		{
			name:                 "passing an empty string",
			keys:                 []string{""},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name:                 "passing a single key (tracked)",
			keys:                 []string{"w"},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name:                 "passing a single key (non-tracked)",
			keys:                 []string{"q"},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name:                 "passing two keys (non-tracked)",
			keys:                 []string{"q", "q"},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name: "passing a tracked key thrice should track a repetition",
			keys: []string{"w", "w", "w"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"www+": {
					Key:   "www+",
					Count: 1,
				},
			},
		},
		{
			name: "passing a tracked key more than thrice should track only a single repetition",
			keys: []string{"w", "w", "w", "w"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"www+": {
					Key:   "www+",
					Count: 1,
				},
			},
		},
		{
			name:                 "ddd is not tracked",
			keys:                 []string{"d", "d", "d"},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name: "dddd is tracked",
			keys: []string{"d", "d", "d", "d"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"dddd+": {
					Key:   "dddd+",
					Count: 1,
				},
			},
		},
		{
			name:                 "alternating tracked key, non-tracked key, ...",
			keys:                 []string{"w", "q", "w", "q", "w", "q"},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name:                 "alternating tracked key, different tracked key, ...",
			keys:                 []string{"w", "e", "w", "j", "w", "j"},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name: "multiple repetitive tracked keys",
			keys: []string{"w", "w", "w", "j", "j", "j"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"www+": {
					Key:   "www+",
					Count: 1,
				},
				"jjj+": {
					Key:   "jjj+",
					Count: 1,
				},
			},
		},
		{
			name: "two times repetitive tracked keys",
			keys: []string{"w", "w", "w", "j", "w", "w", "w"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"www+": {
					Key:   "www+",
					Count: 2,
				},
			},
		},
		{
			name: "move right (l) and press i instead of just a",
			keys: []string{"l", "i"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"li": {
					Key:   "li",
					Count: 1,
				},
			},
		},
		{
			name: "move left (h) and press a instead of just i",
			keys: []string{"h", "a"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"ha": {
					Key:   "ha",
					Count: 1,
				},
			},
		},
		{
			name: "move down (j) and press O instead of just o",
			keys: []string{"j", "O"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"jO": {
					Key:   "jO",
					Count: 1,
				},
			},
		},
		{
			name: "move up (k) and press o instead of just O",
			keys: []string{"k", "o"},
			expectedAntipatterns: map[string]*parser.Antipattern{
				"ko": {
					Key:   "ko",
					Count: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			maxAllowedRepeats := int64(2)
			tracker := parser.NewAntipatternTracker(maxAllowedRepeats)

			for i := range tt.keys {
				tracker.Track(tt.keys[i], parser.NormalMode)
			}

			require.Equal(t, tt.expectedAntipatterns, tracker.Antipatterns())
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
		expectedAntipatterns map[string]*parser.Antipattern
		name                 string
		keys                 []params
	}{
		{
			name: "insert and enter",
			keys: []params{
				{key: "i", currentMode: parser.NormalMode},
				{key: "<cr>", currentMode: parser.InsertMode},
			},
			expectedAntipatterns: map[string]*parser.Antipattern{
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
			expectedAntipatterns: map[string]*parser.Antipattern{
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
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
		{
			name: "pressing li already in insert mode",
			keys: []params{
				{key: "l", currentMode: parser.InsertMode},
				{key: "i", currentMode: parser.InsertMode},
			},
			expectedAntipatterns: map[string]*parser.Antipattern{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			maxAllowedRepeats := int64(2) // doesn't matter for this test
			tracker := parser.NewAntipatternTracker(maxAllowedRepeats)

			for _, currentKey := range tt.keys {
				tracker.Track(currentKey.key, currentKey.currentMode)
			}
			require.Equal(t, tt.expectedAntipatterns, tracker.Antipatterns())
		})
	}
}
