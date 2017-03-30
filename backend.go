package main

import "github.com/a-h/ver/git"

// Backend represents a backend for recording whether git operations have been logged.
type Backend interface {
	GetLease() (id string, err error)
	ExtendLease(id string) (ok bool, err error)
	CancelLease() (ok bool, err error)
	IsProcessed(bool, error)
	MarkProcessed(commit git.Commit) error
}
