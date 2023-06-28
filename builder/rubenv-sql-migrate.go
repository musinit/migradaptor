package builder

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	MigrationDownCmd  = "+migrate Down"
	MigrationUpCmd    = "+migrate Up"
	StatementBeginCmd = "StatementBegin"
	StatementEndCmd   = "StatementEnd"
	NoTransactionCmd  = "notransaction"
)

var (
	filenameReg = regexp.MustCompile(`(\d{0,15})-(.*)(.sql)`)
)

func BuildMigrationData(lines []string) ([]string, []string) {
	upLines, downLines := make([]string, 0, len(lines)/2), make([]string, 0, len(lines)/2)
	isUpTx := true
	upTransactionMode, downTransactionMode := false, false
	for _, line := range lines {
		upMigrationLine := strings.Contains(line, MigrationUpCmd)
		downMigrationLine := strings.Contains(line, MigrationDownCmd)
		switch {
		case upMigrationLine:
			if !strings.Contains(line, NoTransactionCmd) {
				upTransactionMode = true
				upLines = append(upLines, "BEGIN;\n")
			}
			isUpTx = true
		case downMigrationLine:
			if upTransactionMode {
				upLines = append(upLines, "COMMIT;\n")
			}
			if !strings.Contains(line, NoTransactionCmd) {
				downTransactionMode = true
				downLines = append(downLines, "BEGIN;\n")
			}
			isUpTx = false
		case !(strings.Contains(line, StatementBeginCmd) || strings.Contains(line, StatementEndCmd)):
			if isUpTx {
				upLines = append(upLines, line)
			} else {
				downLines = append(downLines, line)
			}
		default:
			fmt.Printf("skip line %s", line)
		}
	}
	if downTransactionMode {
		downLines = append(downLines, "\nCOMMIT;\n")
	}

	return upLines, downLines
}

func ParseFilename(filename string) (int64, string, error) {
	fileparts := filenameReg.FindAllStringSubmatch(filename, 10)
	if !isKeyExists(filenameReg, filename) {
		return 0, "", fmt.Errorf("parse fileparts: filename %s not match", filename)
	}
	fel := fileparts[0]
	// {timestamp}-{name}.sql format is expected
	// 1 - timestamp, 2 - name
	if len(fel) < 2 {
		return 0, "", errors.New("parse file")
	}
	ts := fel[1]
	name := fel[2]
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return 0, "", errors.Wrap(err, "parse timestamp")
	}
	return tsInt, name, nil
}
