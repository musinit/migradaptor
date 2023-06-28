package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/musinit/migradaptor/builder"
)

func Test_ParseFilename(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		expectedTimestamp int64
		extectedName      string
		expectedErr       error
	}{
		{
			"1-test.sql",
			"1-test.sql",
			1,
			"test",
			nil,
		},
		{
			"1-1_initial.sql",
			"1-1_initial.sql",
			1,
			"1_initial",
			nil,
		},
		{
			"1-_.sql",
			"1-_.sql",
			1,
			"_",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, name, err := builder.ParseFilename(tc.input)
			require.True(t, ts == tc.expectedTimestamp)
			require.True(t, name == tc.extectedName)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
