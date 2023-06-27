package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/musinit/migradaptor/builder"
)

func main() {
	var sourceType string
	var srcMigrPath string
	var dstMigrPath string
	helpPtr := flag.Bool("help", false, "print help information")
	flag.StringVar(&sourceType, "source-type", "rubenv-sql-migrate", "source library to convert from")
	flag.StringVar(&srcMigrPath, "src-folder", "src-folder", "source migrations folder")
	flag.StringVar(&dstMigrPath, "dst-folder", "dst-folder", "destination migrations folder")
	flag.Parse()

	if helpPtr != nil && *helpPtr {
		PrintHelp()
		os.Exit(0)
	}

	if err := builder.ValidateInput(&sourceType, &srcMigrPath, &dstMigrPath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "validate error: %s\n Run migrator -help for information.\n", err.Error())
		os.Exit(1)
	}

	srcType, err := builder.GetSourceType(sourceType)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "get sourceType type error: %s\n", err.Error())
		os.Exit(1)
	}

	if _, err := os.Stat(srcMigrPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "source migration directory doesn't exists\n")
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "can't get current directory: %s\n", err.Error())
		os.Exit(1)
	}
	dstMigrPath = path.Join(pwd, dstMigrPath)
	srcMigrPath = path.Join(pwd, srcMigrPath)

	if _, err := os.Stat(dstMigrPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dstMigrPath, os.ModePerm); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "create dest dir error: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		if err := builder.RemoveContents(dstMigrPath); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "clear dest migrations folder error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	files, err := os.ReadDir(srcMigrPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "read legacy migrations folder error: %s\n", err.Error())
		os.Exit(1)
	}
	maxTime := int64(0)

	for _, file := range files {
		if !builder.IsSqlMigrationFile(file.Name()) {
			continue
		}

		lf, err := os.Open(path.Join(srcMigrPath, file.Name()))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "open legacy migrations folder error: %s\n", err.Error())
			os.Exit(1)
		}
		defer func() {
			if err := lf.Close(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "closing legacy migrations folder error: %s\n", err.Error())
				os.Exit(1)
			}
		}()

		lines, err := builder.ReadFileLines(lf)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "reading legacy migrations file lines error: %s\n", err.Error())
			os.Exit(1)
		}

		var upMigr, downMigr []string
		switch srcType {
		default:
			upMigr, downMigr = builder.BuildMigrationData(lines)
		}

		timestamp, name := builder.ParseFilename(file.Name())
		if timestamp <= maxTime {
			timestamp = maxTime + 1
		}
		maxTime = timestamp

		println(fmt.Sprintf("%d : %s", timestamp, name))

		upMgrFn := fmt.Sprintf("%d_%s.up.sql", timestamp, name)
		downMgrFn := fmt.Sprintf("%d_%s.down.sql", timestamp, name)

		// create migration .up file
		if err := builder.CreateAndWrite(dstMigrPath, upMgrFn, upMigr); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "writing destination migrations error: %s\n", err.Error())
			os.Exit(1)
		}
		// migration .down file
		if err := builder.CreateAndWrite(dstMigrPath, downMgrFn, downMigr); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "writing destination migrations error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	println("finished")
}

func PrintHelp() {
	helpText := `
Usage: migradaptor [options] ...

  Migrate your sql migrations files between different lib formats.

Options:

  -source-type=rubenv-sql-migration          Source library of sql files, that need to transform.
  -src-folder="source migrations path"       Source migrations folder.
  -dst-folder="destination migrations path"  Destination migrations folder.
`
	println(helpText)
}
