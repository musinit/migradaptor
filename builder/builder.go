package builder

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	sqlExt          = regexp.MustCompile(`.sql$`)
	specialCharsMap = map[rune]struct{}{
		rune('\t'): {},
		rune('\n'): {},
	}
)

type Cmd interface {
	~string
}

type DstType string

var (
	DstTypeSqlMigrate DstType = "golang-migrate"
)

func GetDstType(sourceType string) (DstType, error) {
	sourceType = strings.TrimSpace(sourceType)
	sourceType = strings.ToLower(sourceType)
	switch sourceType {
	case string(DstTypeSqlMigrate):
		return DstTypeSqlMigrate, nil
	default:
		return *(new(DstType)), ErrUnknownSourceType
	}
}

func IsSqlMigrationFile(filename string) bool {
	return isKeyExists(sqlExt, filename)
}

func isKeyExists(reg *regexp.Regexp, source string) bool {
	fileparts := reg.FindAllStringSubmatch(source, -1)
	return len(fileparts) != 0
}

func RemoveSpecialCharacters(s string) string {
	result := s
	for k := range specialCharsMap {
		result = strings.ReplaceAll(result, string(k), "")
	}
	result = strings.TrimSpace(result)
	return result
}

func JoinMigrationData(lines []string) string {
	result := strings.Builder{}
	for i := range lines {
		if lines[i] != "" {
			result.WriteString(lines[i])
		}
	}
	return result.String()
}

func ValidateInput(dstType, srcPath, dstPath *string) error {
	var errJoin error
	if dstType == nil || (dstType != nil && *dstType == "") {
		errJoin = errors.Join(errJoin, ErrNoDstTypeProvided)
	}
	if srcPath == nil || (srcPath != nil && *srcPath == "") {
		errJoin = errors.Join(errJoin, ErrNoSrcFolderPath)
	}
	if dstPath == nil || (dstPath != nil && *dstPath == "") {
		errJoin = errors.Join(errJoin, ErrNoDstFolderPath)
	}
	if errJoin != nil {
		return errJoin
	}
	if *srcPath == *dstPath {
		return ErrLegacyAndDestEqual
	}
	return nil
}

func IsContainsCmd[T Cmd](src string, substrs ...T) bool {
	for _, substr := range substrs {
		if strings.Contains(src, string(substr)) {
			return true
		}
	}
	return false
}

func BuildMigrationData(lines []string) ([]string, []string) {
	upLines, downLines := make([]string, 0, len(lines)/2), make([]string, 0, len(lines)/2)
	isUpTx := true
	upTransactionMode, downTransactionMode := false, false
	for _, line := range lines {
		upMigrationLine := IsContainsCmd(line,
			string(SqlMigrateCmdMigrationUp),
			string(DbmateCmdMigrationUp),
			string(GooseCmdMigrationUp),
		)
		downMigrationLine := IsContainsCmd(line,
			string(SqlMigrateCmdMigrationDown),
			string(DbmateCmdMigrationDown),
			string(GooseCmdMigrationDown),
		)
		switch {
		case upMigrationLine:
			if !(IsContainsCmd(line, SqlMigrateCmdNoTransaction) ||
				IsContainsCmd(line, DbmateCmdNoTransaction) ||
				IsContainsCmd(line, GooseCmdNoTransaction)) {
				upTransactionMode = true
				upLines = append(upLines, "BEGIN;\n")
			}
			isUpTx = true
		case downMigrationLine:
			if upTransactionMode {
				upLines = append(upLines, "COMMIT;\n")
			}
			if !(IsContainsCmd(line, SqlMigrateCmdNoTransaction) ||
				IsContainsCmd(line, DbmateCmdNoTransaction) ||
				IsContainsCmd(line, GooseCmdNoTransaction)) {
				downTransactionMode = true
				downLines = append(downLines, "BEGIN;\n")
			}
			isUpTx = false
		case IsContainsCmd(line, SqlMigrateCmdStatementBegin) || IsContainsCmd(line, SqlMigrateCmdStatementEnd) ||
			IsContainsCmd(line, GooseCmdStatementBegin) || IsContainsCmd(line, GooseCmdStatementEnd):
			fmt.Printf("skip line %s", line)
		default:
			if isUpTx {
				upLines = append(upLines, line)
			} else {
				downLines = append(downLines, line)
			}

		}
	}
	if downTransactionMode {
		downLines = append(downLines, "\nCOMMIT;\n")
	}

	return upLines, downLines
}
