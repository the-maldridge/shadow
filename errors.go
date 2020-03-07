package shadow

import "errors"

var (
	// ErrWrongNumFields is returned when a record which should
	// have a fixed number of fields is parsed with an input which
	// has an improper number of fields.
	ErrWrongNumFields = errors.New("wrong number of fields in provided record")

	// ErrNotANumber is returned when a field that must be a
	// number cannot be parsed as one.
	ErrNotANumber = errors.New("atoi failed during numerical parse")
)
