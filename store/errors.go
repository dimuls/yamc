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
		s += ": " + e.Details
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

	// cleaning errors
	ErrFailToCreateCleaning  = e(40, "fail to create cleaning")
	ErrFailToStartCleaning   = e(41, "fail to start cleaning")
	ErrCleaningNotStartedYet = e(42, "cleaning not started yet")
	ErrFailToStopCleaning    = e(43, "fail to stop cleaning")

	// dumping errors
	ErrFailToCreateDumping = e(50, "fail to create dumping")
	ErrFailToStartDumping  = e(51, "fail to start dumping")
	ErrDumperNotStartedYet = e(52, "dumping not started yet")
	ErrFailToStopDumping   = e(53, "fail to stop dumping")

	// load errors
	ErrFailOpenDumpFile     = e(60, "fail to open dump file")
	ErrFailToDumpItems      = e(61, "fail to dump items")
	ErrFailToDecodeDumpFile = e(61, "fail to decode dump file")
	ErrFailToCloseDumpFile  = e(62, "fail to close dump file")
)
