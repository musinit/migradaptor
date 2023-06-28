# migradaptor
> Tool for adapting migration files for different library formats. Current version allows to adapt [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate) files to [golang-migrate](https://github.com/golang-migrate).

[![Test](https://github.com/musinit/migradaptor/actions/workflows/test.yml/badge.svg)](https://github.com/musinit/migradaptor/actions/workflows/test.yml) 

About
---------
Once I faced with an issue to change our corporate migrations library in Golang from [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate)
to [golang-migrate](https://github.com/golang-migrate) and it was so hard.
There are several caveats that you should know about, so it can save you some time:
 - You have to split each sql-migrate file in 2 files: up and down migration.
 - By default, migrations in sql-migrate are running in transaction-default mode, but in golang-migrate you have to wrap it in BEGIN;COMMIT; for transactions.
 - Golang-migrate [doesn't like](https://github.com/golang-migrate/migrate/issues/731) the same timestamp for different files.
 - You [can't](https://github.com/golang-migrate/migrate/issues/284) create several indexes concurrently without adding x-multi-statement=true flag for DB connection. 
However, please note that this flag [will break](https://github.com/golang-migrate/migrate/issues/590) your CREATE FUNCTION ... AS $$ symbol.
  
That's why I decided to start this lib.
I hope there will be more sources (like rubenv/sql-migrate), so people can save time if they need to change their migration lib and adapt their migration files from one format to another.
Feel free to suggest an Issue or PR.

## Getting started
Install
```bash
go install github.com/musinit/migradaptor/...@latest
```

Use
```bash
migradaptor -src={source_folder} -dst={destination_folder}
```

## Questions or Feedback?

You can use GitHub Issues for feedback or questions.

## TODO
 - processing multiple concurrent indexes in single file, splitting the file by the number of such CREATE INDEX CONCURRENTLY commands.