package builder

import "github.com/pkg/errors"

var (
	ErrUnknownSourceType  = errors.New("unknown source type")
	ErrNoDstTypeProvided  = errors.New("no source type provided")
	ErrNoSrcFolderPath    = errors.New("no dst folder path provided")
	ErrNoDstFolderPath    = errors.New("no dst folder path provided")
	ErrLegacyAndDestEqual = errors.New("src and dst path are equal")
)
