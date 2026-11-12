package overlays

import "errors"

var (
	// ErrNoExitKey indicates an overlay was shown without an exit key configured.
	ErrNoExitKey = errors.New("overlay must have an exit key defined: call SetExitKey()")

	// ErrKeyConflict indicates an overlay key conflicts with a hardcoded binding.
	ErrKeyConflict = errors.New("overlay key conflicts with hardcoded binding")
)
