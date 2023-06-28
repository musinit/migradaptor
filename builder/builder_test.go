package builder_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/musinit/migradaptor/builder"
	"github.com/musinit/migradaptor/utils"
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
		sourceType *string
		srcPath    *string
		dstPath    *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "each argument is nil",
			args: args{
				sourceType: nil,
				srcPath:    nil,
				dstPath:    nil,
			},
			wantErr: errors.Join(
				builder.ErrNoSourceTypeProvided,
				builder.ErrNoDstFolderPath,
				builder.ErrNoSrcFolderPath,
			),
		},
		{
			name: "each argument is empty",
			args: args{
				sourceType: utils.Ptr(""),
				srcPath:    utils.Ptr(""),
				dstPath:    utils.Ptr(""),
			},
			wantErr: errors.Join(
				builder.ErrNoSourceTypeProvided,
				builder.ErrNoDstFolderPath,
				builder.ErrNoSrcFolderPath,
			),
		},
		{
			name: "source type is empty",
			args: args{
				sourceType: utils.Ptr(""),
				srcPath:    utils.Ptr("1"),
				dstPath:    utils.Ptr("1"),
			},
			wantErr: errors.Join(
				builder.ErrNoSourceTypeProvided,
			),
		},
		{
			name: "source type invalid format",
			args: args{
				sourceType: utils.Ptr("rubeenv_wrong"),
				srcPath:    utils.Ptr("2"),
				dstPath:    utils.Ptr("1"),
			},
			wantErr: errors.Join(
				builder.ErrUnknownSourceType,
			),
		},
		{
			name: "src path empty",
			args: args{
				sourceType: utils.Ptr("rubenv-sql-migrate"),
				srcPath:    utils.Ptr(""),
				dstPath:    utils.Ptr("1"),
			},
			wantErr: errors.Join(
				builder.ErrNoSrcFolderPath,
			),
		},
		{
			name: "src dst empty",
			args: args{
				sourceType: utils.Ptr("rubenv-sql-migrate"),
				srcPath:    utils.Ptr("1"),
				dstPath:    utils.Ptr(""),
			},
			wantErr: errors.Join(
				builder.ErrNoDstFolderPath,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.ValidateInput(tt.args.sourceType, tt.args.srcPath, tt.args.dstPath)
			require.Error(t, tt.wantErr, err)
		})
	}
}

func TestIsSubstringExists(t *testing.T) {
	type args struct {
		source string
		substr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "MigrationUpCmd exists",
			args: args{
				source: `
					-- +migrate Up
					CREATE SCHEMA IF NOT EXISTS gmf_go;
					SET search_path TO gmf_go;
					`,
				substr: builder.MigrationUpCmd,
			},
			want: true,
		},
		{
			name: "MigrationUpCmd does not exists",
			args: args{
				source: `
					-- +migrate
					CREATE SCHEMA IF NOT EXISTS gmf_go;
					SET search_path TO gmf_go;
					`,
				substr: builder.MigrationUpCmd,
			},
			want: false,
		},
		{
			name: "MigrationDownCmd exists",
			args: args{
				source: `
					-- +migrate Down
					CREATE SCHEMA IF NOT EXISTS gmf_go;
					SET search_path TO gmf_go;
					`,
				substr: builder.MigrationDownCmd,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings.Contains(tt.args.source, tt.args.substr); got != tt.want {
				t.Errorf("strings.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
