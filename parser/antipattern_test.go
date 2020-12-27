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
		// {
		// 	name:                 "passing an empty string",
		// 	keys:                 []string{""},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{},
		// },
		// {
		// 	name:                 "passing a single key (tracked)",
		// 	keys:                 []string{"w"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{},
		// },
		// {
		// 	name:                 "passing a single key (non-tracked)",
		// 	keys:                 []string{"q"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{},
		// },
		// {
		// 	name:                 "passing two keys (non-tracked)",
		// 	keys:                 []string{"q", "q"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{},
		// },
		// {
		// 	name: "passing a tracked key twice should track a repetition",
		// 	keys: []string{"w", "w"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{
		// 		"ww": {
		// 			Key:   "ww",
		// 			Count: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name:                 "alternating tracked key, non-tracked key, ...",
		// 	keys:                 []string{"w", "q", "w", "q", "w", "q"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{},
		// },
		// {
		// 	name:                 "alternating tracked key, different tracked key, ...",
		// 	keys:                 []string{"w", "e", "w", "j", "w", "j"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{},
		// },
		// {
		// 	name: "multiple repetitive tracked keys",
		// 	keys: []string{"w", "w", "w", "j", "j"},
		// 	expectedAntipatterns: map[string]*parser.Antipattern{
		// 		"ww": {
		// 			Key:   "ww",
		// 			Count: 2,
		// 		},
		// 		"jj": {
		// 			Key:   "jj",
		// 			Count: 1,
		// 		},
		// 	},
		// },
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
			tracker := parser.NewAntipatternTracker()
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tracker := parser.NewAntipatternTracker()
			for _, currentKey := range tt.keys {
				tracker.Track(currentKey.key, currentKey.currentMode)
			}
			require.Equal(t, tt.expectedAntipatterns, tracker.Antipatterns())
		})
	}
}