package builder

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type SqlMigrateCmd string

var (
	SqlMigrateCmdMigrationDown  SqlMigrateCmd = "+migrate Down"
	SqlMigrateCmdMigrationUp    SqlMigrateCmd = "+migrate Up"
	SqlMigrateCmdStatementBegin SqlMigrateCmd = "StatementBegin"
	SqlMigrateCmdStatementEnd   SqlMigrateCmd = "StatementEnd"
	SqlMigrateCmdNoTransaction  SqlMigrateCmd = "notransaction"
)

var (
	filenameReg = regexp.MustCompile(`(\d{0,15})(-|_)(.*)(.sql)`)
)

func ParseFilename(filename string) (int64, string, error) {
	fileparts := filenameReg.FindAllStringSubmatch(filename, 10)
	if !isKeyExists(filenameReg, filename) {
		return 0, "", fmt.Errorf("parse fileparts: filename %s not match", filename)
	}
	fel := fileparts[0]
	// {timestamp}[-|_]{name}.sql format is expected
	// 1 - timestamp, 2 - name
	if len(fel) < 3 {
		return 0, "", errors.New("parse file")
	}
	ts := fel[1]
	name := fel[3]
	println(ts)
	println(name)
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return 0, "", errors.Wrap(err, "parse timestamp")
	}
	return tsInt, name, nil
}
