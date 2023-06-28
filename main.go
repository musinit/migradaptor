package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"github.com/musinit/migradaptor/builder"
)

func main() {
	var (
		sourceType  string
		srcMigrPath string
		dstMigrPath string
		flgVersion  bool
		helpPtr     bool
	)
	flag.BoolVar(&flgVersion, "version", false, "if true, print version and exit")
	flag.BoolVar(&helpPtr, "help", false, "print help information")
	flag.StringVar(&sourceType, "source-type", "rubenv-sql-migrate", "source library to convert from")
	flag.StringVar(&srcMigrPath, "src", "src", "source migrations folder")
	flag.StringVar(&dstMigrPath, "dst", "dst", "destination migrations folder")
	flag.Parse()

	switch {
	case flgVersion:
		println(GetVersion())
		os.Exit(0)
	case helpPtr:
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
		_, _ = fmt.Fprint(os.Stderr, "source migration directory doesn't existss")
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
		_, _ = fmt.Fprintf(os.Stderr, "read src migrations folder error: %s\n", err.Error())
		os.Exit(1)
	}
	maxTime := int64(0)

	for _, file := range files {
		if !builder.IsSqlMigrationFile(file.Name()) {
			continue
		}

		lf, err := os.Open(path.Join(srcMigrPath, file.Name()))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "open src migrations folder error: %s\n", err.Error())
			os.Exit(1)
		}
		defer func() {
			if err := lf.Close(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "closing src migrations folder error: %s\n", err.Error())
				os.Exit(1)
			}
		}()

		lines, err := builder.ReadFileLines(lf)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "reading src migrations file lines error: %s\n", err.Error())
			os.Exit(1)
		}

		var upMigr, downMigr []string
		switch srcType {
		default:
			upMigr, downMigr = builder.BuildMigrationData(lines)
		}

		timestamp, name, err := builder.ParseFilename(file.Name())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "filename parsing error: %s\n", err.Error())
			os.Exit(1)
		}
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

  -source-type=rubenv-sql-migration   Source library of sql files, that need to transform.
  -src="source migrations path"       Source migrations folder.
  -dst="destination migrations path"  Destination migrations folder.
`
	println(helpText)
}

func GetVersion() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok && buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}
	return "dev"
}
