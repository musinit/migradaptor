package tests_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/musinit/migradaptor/builder"
)

func TestDbmate_BuildMigrationData(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		sqlLines   string
		resultUp   string
		resultDown string
	}{
		{
			name: "simple one line",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);
					
					
					-- migrate:down
					DROP TABLE companies;`,
			resultUp: `BEGIN;
						CREATE TABLE companies (id int, title string);
						COMMIT;`,
			resultDown: `BEGIN;
						DROP TABLE companies;
						COMMIT;
						`,
		},
		{
			name: "simple two lines",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					
					
					-- migrate:down
					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
					`,
			resultUp: `BEGIN;
						CREATE TABLE companies (id int, title string);
						CREATE INDEX companies_title_idx on companies (title);
						COMMIT;`,
			resultDown: `BEGIN;
						DROP INDEX IF EXISTS companies_title_idx;
						DROP TABLE IF EXISTS companies;
						COMMIT;
						`,
		},
		{
			name: "two lines - up without transactions",
			sqlLines: `-- migrate:up transaction:false
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					
					
					-- migrate:down
					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
					`,
			resultUp: `CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);`,
			resultDown: `BEGIN;
						DROP INDEX IF EXISTS companies_title_idx;
						DROP TABLE IF EXISTS companies;
						COMMIT;
						`,
		},
		{
			name: "two lines - down without transactions",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					
					
					-- migrate:down transaction:false
					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
					`,
			resultUp: `BEGIN;
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					COMMIT;`,
			resultDown: `
						DROP INDEX IF EXISTS companies_title_idx;
						DROP TABLE IF EXISTS companies;
						`,
		},
		{
			name: "two lines - down without transactions",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					
					
					-- migrate:down transaction:false
					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
					`,
			resultUp: `BEGIN;
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					COMMIT;`,
			resultDown: `
						DROP INDEX IF EXISTS companies_title_idx;
						DROP TABLE IF EXISTS companies;
						`,
		},
		{
			name: "two lines - statement in up",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);

					-- +goose StatementBegin
					CREATE OR REPLACE FUNCTION do_something()
					returns void AS $$
					DECLARE
					  create_query text;
					BEGIN
					  -- Do something here
					END;
					$$
					language plpgsql;
					-- +goose StatementEnd
					
					
					-- migrate:down
					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
					`,
			resultUp: `BEGIN;
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					CREATE OR REPLACE FUNCTION do_something()
					returns void AS $$
					DECLARE
					  create_query text;
					BEGIN
					  -- Do something here
					END;
					$$
					language plpgsql;
					COMMIT;`,
			resultDown: `BEGIN;
						DROP INDEX IF EXISTS companies_title_idx;
						DROP TABLE IF EXISTS companies;
						COMMIT;`,
		},
		{
			name: "two lines - statement in down with no transaction",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					
					
					-- migrate:down transaction:false

					-- +goose StatementBegin

					CREATE OR REPLACE FUNCTION do_something()
					returns void AS $$
					DECLARE
					  create_query text;
					BEGIN
					  -- Do something here
					END;
					$$
					language plpgsql;
					-- +goose StatementEnd

					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
					`,
			resultUp: `BEGIN;
					CREATE TABLE companies (id int, title string);
					CREATE INDEX companies_title_idx on companies (title);
					COMMIT;`,
			resultDown: `
					CREATE OR REPLACE FUNCTION do_something()
					returns void AS $$
					DECLARE
					  create_query text;
					BEGIN
					  -- Do something here
					END;
					$$
					language plpgsql;
					DROP INDEX IF EXISTS companies_title_idx;
					DROP TABLE IF EXISTS companies;
						`,
		},
		{
			name: "every command on one line - statement in down with no transaction",
			sqlLines: `-- migrate:up
					CREATE TABLE companies (id int, title string);CREATE INDEX companies_title_idx on companies (title);

					-- migrate:down transaction:false

					-- +goose StatementBegin

					CREATE OR REPLACE FUNCTION do_something() returns void AS $$ DECLARE
					  create_query text;
					BEGIN
					END;
					$$
					language plpgsql;
					-- +goose StatementEnd

					DROP INDEX IF EXISTS companies_title_idx;DROP TABLE IF EXISTS companies;
					`,
			resultUp: `BEGIN;
					CREATE TABLE companies (id int, title string);CREATE INDEX companies_title_idx on companies (title);
					COMMIT;`,
			resultDown: `
					CREATE OR REPLACE FUNCTION do_something() returns void AS $$ DECLARE
					  create_query text;
					BEGIN
					END;
					$$
					language plpgsql;
					DROP INDEX IF EXISTS companies_title_idx;DROP TABLE IF EXISTS companies;
						`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			upLines, downLines := builder.BuildMigrationData(strings.Split(c.sqlLines, "\n"))
			upJoin := builder.JoinMigrationData(upLines)
			downJoin := builder.JoinMigrationData(downLines)
			upGot := builder.RemoveSpecialCharacters(upJoin)
			downGot := builder.RemoveSpecialCharacters(downJoin)
			upWant := builder.RemoveSpecialCharacters(c.resultUp)
			downWant := builder.RemoveSpecialCharacters(c.resultDown)

			require.Equal(t, upWant, upGot)
			require.Equal(t, downWant, downGot)
		})
	}
}
