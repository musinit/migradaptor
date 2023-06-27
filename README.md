# Migradaptor

About
---------
Once I faced with a issue to change our corporate migrations library in Golang from [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate)
to [golang-migrate](https://github.com/golang-migrate) and it was so hard.
There are several caveats that you should know about, so it can save you some time:
 - You have to split each sql-migrate file in 2 files: up and down migration.
 - By default migrations in sql-migrate are running in transaction-default mode, but in golang-migrate you have to wrap it in BEGIN;COMMIT; for transactions.
 - Golang-migrate [doesn't like](https://github.com/golang-migrate/migrate/issues/731) the same timestamp for different files.
 - You [can't](https://github.com/golang-migrate/migrate/issues/284) create several indexes concurrently without adding x-multi-statement=true flag for DB connection. 
However, please note that this flag [will break](https://github.com/golang-migrate/migrate/issues/590) your CREATE FUNCTION ... AS $$ symbol.
  
That's why I decided to start this lib.
I hope there will be more sources (like rubenv/sql-migrate), so people can save time if they need to change their migration lib and adapt their migration files from one format to another.
Feel free to suggest an Issue or PR.

## Getting started
Install
```bash
$ go install github.com/musinit/migradaptor/...@latest
```

Use
```bash
$ migradaptor -src-folder={source_folder} -dst-folder={destination_folder}
```

TODO
 - processing multiple concurrent indexes in single file, splitting the file by the number of such CREATE CONCURRENT INDEX.