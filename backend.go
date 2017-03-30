package main

// Backend represents a backend for recording whether git operations have been logged.
type Backend interface {
	GetLease() (id string, err error)
	ExtendLease(id string) (ok bool, err error)
	CancelLease() (err error)
	IsProcessed(hash string) (bool, error)
	MarkProcessed(hash string) error
}
