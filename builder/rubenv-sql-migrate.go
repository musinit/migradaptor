package builder

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	MigrationDownCmd  = "migrate Down"
	MigrationUpCmd    = "migrate Up"
	StatementBeginCmd = "StatementBegin"
	StatementEndCmd   = "StatementEnd"
	NoTransactionCmd  = "notransaction"
)

var (
	ErrNoUpOrDownMigrationPart = errors.New("no up or down migration part")
)

var (
	filenameReg = regexp.MustCompile("(\\d{0,15})-(.*)(.sql)")
)

func BuildMigrationData(lines []string) ([]string, []string) {
	upLines, downLines := make([]string, 0, len(lines)/2), make([]string, 0, len(lines)/2)
	isUpTx := true
	upTransactionMode, downTransactionMode := false, false
	for i := range lines {
		line := lines[i]
		upMigrationLine := isSubstringExists(MigrationUpCmd, line)
		downMigrationLine := isSubstringExists(MigrationDownCmd, line)
		if upMigrationLine {
			if !isSubstringExists(NoTransactionCmd, line) {
				upTransactionMode = true
				upLines = append(upLines, "BEGIN;\n")
			}
			isUpTx = true
		} else if downMigrationLine {
			if upTransactionMode {
				upLines = append(upLines, "COMMIT;\n")
			}
			if !isSubstringExists(NoTransactionCmd, line) {
				downTransactionMode = true
				downLines = append(downLines, "BEGIN;\n")
			}
			isUpTx = false
		} else {
			if isSubstringExists(StatementBeginCmd, line) || isSubstringExists(StatementEndCmd, line) {
				upLines = append(upLines, "\n")
			} else {
				if isUpTx {
					upLines = append(upLines, line)
				} else {
					downLines = append(downLines, line)
				}
			}

			upLines = append(upLines, "\n")
		}
	}
	if downTransactionMode {
		downLines = append(downLines, "\nCOMMIT;\n")
	}

	return upLines, downLines
}

func ParseFilename(filename string) (int64, string) {
	fileparts := filenameReg.FindAllStringSubmatch(filename, 10)
	if !isKeyExists(filenameReg, filename) {
		panic("no parts")
	}
	fel := fileparts[0]
	// {timestamp}-{name}.sql format is expected
	// 1 - timestamp, 2 - name
	if len(fel) < 2 {
		panic(errors.New("can't parse file " + filename))
	}
	ts := fel[1]
	name := fel[2]
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		panic(err)
	}
	return tsInt, name
}
