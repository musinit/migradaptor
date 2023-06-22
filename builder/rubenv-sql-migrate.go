package builder

import (
	"bytes"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	filenameReg      = regexp.MustCompile("(\\d{0,15})-(.*)(.sql)")
	noTransactionReg = regexp.MustCompile("notransaction")
	statementBegin   = regexp.MustCompile("StatementBegin")
	statementEnd     = regexp.MustCompile("StatementEnd")
	migrationUpReg   = regexp.MustCompile("migrate Up")
	migrationDownReg = regexp.MustCompile("migrate Down")
)

func BuildMigrationData(lines []string) ([]string, []string) {
	upLines, downLines := make([]string, 0, len(lines)/2), make([]string, 0, len(lines)/2)
	isUpTx := true
	upTransactionMode, downTransactionMode := false, false
	for i := range lines {
		line := lines[i]
		upMigrationLine := isKeyExists(migrationUpReg, line)
		downMigrationLine := isKeyExists(migrationDownReg, line)
		if upMigrationLine {
			if !isKeyExists(noTransactionReg, line) {
				upTransactionMode = true
				upLines = append(upLines, "BEGIN;\n")
			}
			isUpTx = true
		} else if downMigrationLine {
			if upTransactionMode {
				upLines = append(upLines, "COMMIT;\n")
			}
			if !isKeyExists(noTransactionReg, line) {
				downTransactionMode = true
				downLines = append(downLines, "BEGIN;\n")
			}
			isUpTx = false
		} else {
			if isKeyExists(statementBegin, line) || isKeyExists(statementEnd, line) {
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

func BuildMigrationDataBuffer(lines []string) (bytes.Buffer, bytes.Buffer) {
	var upMigr bytes.Buffer
	var downMigr bytes.Buffer
	isUpTx := true
	upTransactionMode, downTransactionMode := false, false
	for i := range lines {
		line := lines[i]
		upMigrationLine := isKeyExists(migrationUpReg, line)
		downMigrationLine := isKeyExists(migrationDownReg, line)
		if upMigrationLine {
			if !isKeyExists(noTransactionReg, line) {
				upTransactionMode = true
				upMigr.Write([]byte("BEGIN;\n"))
			}
			isUpTx = true
		} else if downMigrationLine {
			if upTransactionMode {
				upMigr.Write([]byte("COMMIT;\n"))
			}
			if !isKeyExists(noTransactionReg, line) {
				downTransactionMode = true
				downMigr.Write([]byte("BEGIN;\n"))
			}
			isUpTx = false
		} else {
			if isKeyExists(statementBegin, line) || isKeyExists(statementEnd, line) {
				upMigr.Write([]byte("\n"))
			} else {
				if isUpTx {
					upMigr.Write([]byte(line))
				} else {
					downMigr.Write([]byte(line))
				}
			}

			upMigr.Write([]byte("\n"))
		}
	}
	if downTransactionMode {
		downMigr.Write([]byte("\n"))
		downMigr.Write([]byte("COMMIT;\n"))
	}
	return upMigr, downMigr
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
