package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/musinit/migradaptor/builder"
)

func main() {
	source := flag.String("source", "source", "rubenv-sql-migrate")
	sourceMigrationPath := flag.String("source-folder", "legacy_migrations", "source migrations folder")
	destMigrationPath := flag.String("dest-folder", "migrations", "destination migrations folder")
	flag.Parse()
	sourceType, err := builder.GetSourceType(*source)
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(*sourceMigrationPath)
	if err != nil {
		log.Fatal(err)
	}
	maxTime := int64(0)

	for _, file := range files {
		if !builder.IsSqlMigrationFile(file.Name()) {
			continue
		}

		lf, err := os.Open(path.Join(*sourceMigrationPath, file.Name()))
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := lf.Close(); err != nil {
				panic(err)
			}
		}()

		lines, err := builder.ReadFileLines(lf)
		if err != nil {
			panic(err)
		}

		var upMigr, downMigr bytes.Buffer
		switch sourceType {
		default:
			upMigr, downMigr = builder.BuildMigrationDataBuffer(lines)
		}

		timestamp, name := builder.ParseFilename(file.Name())
		if timestamp <= maxTime {
			timestamp = maxTime + 1
		}
		maxTime = timestamp

		println(fmt.Sprintf("%d : %s", timestamp, name))

		upMgrFile := fmt.Sprintf("%d_%s.up.sql", timestamp, name)
		downMgrFile := fmt.Sprintf("%d_%s.down.sql", timestamp, name)

		// create migration .up file
		if err := builder.CreateAndWrite(*destMigrationPath, upMgrFile, &upMigr); err != nil {
			panic(err)
		}
		// migration .down file
		if err := builder.CreateAndWrite(*destMigrationPath, downMgrFile, &downMigr); err != nil {
			panic(err)
		}
	}

	println("finished")
}
