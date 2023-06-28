package builder

import (
	"errors"
	"regexp"
	"strings"
)

var (
	sqlExt                    = regexp.MustCompile(`\\.sql$`)
	concurrentIndexesGroupReg = regexp.MustCompile(`(?:(CREATE INDEX CONCURRENTLY)\\s+(?P<indexName>\\w+).+?;)`)
	specialCharsMap           = map[rune]struct{}{
		rune('\t'): {},
		rune('\n'): {},
	}
)

type SourceType string

var (
	SourceTypeRubenvSqlMigrate SourceType = "rubenv-sql-migrate"
)

func GetSourceType(sourceType string) (SourceType, error) {
	sourceType = strings.TrimSpace(sourceType)
	sourceType = strings.ToLower(sourceType)
	switch sourceType {
	case string(SourceTypeRubenvSqlMigrate):
		return SourceTypeRubenvSqlMigrate, nil
	default:
		return *(new(SourceType)), ErrUnknownSourceType
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

func ValidateInput(srcType, srcPath, dstPath *string) error {
	var errJoin error
	if srcType == nil || (srcType != nil && *srcType == "") {
		errJoin = errors.Join(errJoin, ErrNoSourceTypeProvided)
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

func FindUniqueConcurrentIdxStatements(lineJoin string) []string {
	matches := concurrentIndexesGroupReg.FindAllStringSubmatch(lineJoin, -1)

	result := make([]string, 0)

	for _, match := range matches {
		result = append(result, match[0])
	}
	return result
}
