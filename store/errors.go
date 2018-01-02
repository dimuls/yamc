package store

// StoreError is a store error
type StoreError struct {
	Err     string
	Code    int
	Details string
}

// detailed adds details to error
func (e StoreError) detailed(details string) StoreError {
	e.Details = details
	return e
}

// Error returns error's string representation
func (e StoreError) Error() string {
	s := e.Err
	if e.Details != "" {
		s += ":" + e.Details
	}
	return s
}

// e is StoreError constructor
func e(code int, err string) StoreError {
	return StoreError{Err: err, Code: code}
}

var (
	// store construction errors
	ErrInvalidParams = e(1, "invalid params")
	ErrNilClock      = e(2, "nil clock")

	// item errors
	ErrNotKeyItem  = e(10, "not key item")
	ErrNotListItem = e(11, "not list item")
	ErrNotDictItem = e(12, "not dict item")

	// not exists errors
	ErrKeyNotExists       = e(20, "key not exists")
	ErrListIndexNotExists = e(21, "list index not exists")
	ErrDictKeyNotExists   = e(22, "dict key not exists")

	// other errors
	ErrInvalidListIndex = e(30, "invalid list index")

	// cleaner errors
	ErrFailedToCreateCleaner = e(40, "failed to create cleaner")
	ErrFailedToStartCleaner  = e(41, "failed to start cleaner")
	ErrCleanerNotStartedYet  = e(42, "cleaner not started yet")
	ErrFailedToStopCleaner   = e(43, "failed to stop cleaner")
)
