package builder

import "github.com/pkg/errors"

var (
	ErrUnknownSourceType  = errors.New("unknown source type")
	ErrNoSourceProvided   = errors.New("no source provided")
	ErrNoSrcFolderPath    = errors.New("no source folder path provided")
	ErrNoDstFolderPath    = errors.New("no destination folder path provided")
	ErrLegacyAndDestEqual = errors.New("legacy and destination path are equal")
)
