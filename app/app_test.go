package app_test

import (
	"strings"
	"testing"

	"github.com/phux/vimkeypressanalyzer/app"
	"github.com/phux/vimkeypressanalyzer/parser"
	"github.com/stretchr/testify/require"
)

func TestApp_Analyze(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                            string
		input                           string
		expectedTotalCount              int64
		expectedTotalCountWithoutInsert int64
	}{
		{
			name:                            "Empty input",
			input:                           "",
			expectedTotalCount:              0,
			expectedTotalCountWithoutInsert: 0,
		},
		{
			name:                            "single character input",
			input:                           "j",
			expectedTotalCount:              1,
			expectedTotalCountWithoutInsert: 1,
		},
		{
			name:                            "ji<esc>",
			input:                           "ji" + string(parser.CharEsc),
			expectedTotalCount:              3,
			expectedTotalCountWithoutInsert: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := parser.NewParser(false)
			a := app.NewApp(p)
			input := strings.NewReader(tt.input)

			limit := int64(0)
			excludeModes := []string{parser.InsertMode}
			result, err := a.Analyze(input, limit, excludeModes)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTotalCount, result.TotalKeypresses)
			require.Equal(t, tt.expectedTotalCountWithoutInsert, result.TotalKeypressesWithoutExcludedModes)
		})
	}
}
