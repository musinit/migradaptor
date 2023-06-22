package builder

import "github.com/pkg/errors"

var (
	ErrUnknownSourceType  = errors.New("unknown source type")
	ErrNoSourceProvided   = errors.New("no source provided. Run migradaptor -help for more information")
	ErrNoSrcFolderPath    = errors.New("no source folder path provided. Run migradaptor -help for more information")
	ErrNoDstFolderPath    = errors.New("no destination folder path provided. Run migradaptor -help for more information")
	ErrLegacyAndDestEqual = errors.New("legacy and destination path are equal. Run migradaptor -help for more information")
)
