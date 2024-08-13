package serror

import "errors"

// sentinel or signal errors
var (
	ErrInvalidConfigPath = errors.New("error, invalid config path")
	ErrFileNotExists     = errors.New("error, a file is not exists")

	// Database errors
	ErrNilConfigStruct = errors.New("error, config structure is nil")
)
