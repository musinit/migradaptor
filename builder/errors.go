package builder

import "github.com/pkg/errors"

var (
	ErrUnknownSourceType    = errors.New("unknown source type")
	ErrNoSourceTypeProvided = errors.New("no source type provided")
	ErrNoSrcFolderPath      = errors.New("no source folder path provided")
	ErrNoDstFolderPath      = errors.New("no destination folder path provided")
	ErrLegacyAndDestEqual   = errors.New("legacy and destination path are equal")
)
