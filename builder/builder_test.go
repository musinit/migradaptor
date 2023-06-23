package builder_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/musinit/migradaptor/builder"
)

func TestGetSourceType(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		source  string
		result  builder.SourceType
		wantErr error
	}{
		{
			name:    "valid rubenv-sql-migrate type",
			source:  "rubenv-sql-migrate",
			result:  builder.SourceTypeRubenvSqlMigrate,
			wantErr: nil,
		},
		{
			name:    "valid rubenv-sql-migrate type with spaces",
			source:  "  rubenv-sql-migrate ",
			result:  builder.SourceTypeRubenvSqlMigrate,
			wantErr: nil,
		},
		{
			name:    "valid rubenv-sql-migrate type capital letters",
			source:  "  RUBENV-SQL-MIGRATE ",
			result:  builder.SourceTypeRubenvSqlMigrate,
			wantErr: nil,
		},
		{
			name:    "invalid unknown type capital letters",
			source:  "  some_type ",
			wantErr: builder.ErrUnknownSourceType,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			st, err := builder.GetSourceType(tc.source)

			require.Equal(t, st, tc.result)
			require.True(t, err == tc.wantErr)
		})
	}
}

func Test_IsSqlMigrationFile(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"space",
			" ",
			false,
		},
		{
			"dirty_mess",
			"_?2_d3.]//",
			false,
		},
		{
			".sqll",
			".sqll",
			false,
		},
		{
			" .sqll ",
			".sqll ",
			false,
		},
		{
			".sql",
			".sql",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			have := builder.IsSqlMigrationFile(tc.input)
			require.True(t, have == tc.expected)
		})
	}
}

func Test_RemoveSpecialCharacters(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"space",
			" ",
			"",
		},
		{
			"\thello world\n",
			"\thello world\n",
			"hello world",
		},
		{
			"space around with special characters",
			"\t  hello world \n",
			"hello world",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			have := builder.RemoveSpecialCharacters(tc.input)
			require.True(t, have == tc.expected)
		})
	}
}

func Test_JoinMigrationData(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			"no empty",
			[]string{
				"1",
				"2",
				"3",
			},
			"123",
		},
		{
			"one empty",
			[]string{
				"",
				"2",
				"3",
			},
			"23",
		},
		{
			"all empty",
			[]string{
				"",
				"",
				"",
			},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			have := builder.JoinMigrationData(tc.input)
			require.Equal(t, have, tc.expected)
		})
	}
}

func TestValidateInput(t *testing.T) {
	type args struct {
		source     *string
		legacyPath *string
		path       *string
		err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "simple",
			args: args{
				source:     nil,
				legacyPath: nil,
				path:       nil,
			},
			wantErr: errors.Join(
				builder.ErrNoSourceProvided,
				builder.ErrNoSrcFolderPath,
				builder.ErrNoDstFolderPath,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.ValidateInput(tt.args.source, tt.args.legacyPath, tt.args.path)
			require.ErrorAs(t, tt.wantErr, err)
		})
	}
}

func TestFindUniqueConcurrentIdxStatements(t *testing.T) {
	type args struct {
		lineJoin string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "several concurrent index with drop line",
			args: args{
				lineJoin: `
				DROP INDEX companies_id_idx;
				CREATE INDEX CONCURRENTLY companies_id_idx ON companies (id);
				
				DROP INDEX companies_title_idx;
				CREATE INDEX CONCURRENTLY companies_title_idx ON companies (title);
				
				DROP INDEX clients_id_idx;
				CREATE INDEX CONCURRENTLY clients_id_idx ON clients;`,
			},
			want: []string{
				"companies_id_idx",
				"companies_title_idx",
				"clients_id_idx",
			},
		},
		{
			name: "one concurrent without drop line",
			args: args{
				lineJoin: `
				CREATE INDEX CONCURRENTLY companies_id_idx ON companies (id);`,
			},
			want: []string{
				"companies_id_idx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := builder.FindUniqueConcurrentIdxStatements(tt.args.lineJoin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUniqueConcurrentIdxStatements() = %v, want %v", got, tt.want)
			}
		})
	}
}
