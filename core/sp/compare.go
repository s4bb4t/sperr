package sp

import (
	"errors"
	"hash"
	"slices"
)

// Is returns true if the error is the same as err. Implements the errors.Is
// interface for error comparison. It returns true if the error's cause matches
// the provided error.
func (e *Error) Is(err error) bool {
	return errors.Is(e.cause, err)
}

// DeepIs traverses the entire Error chain to find a matching error. Returns true if
// the provided error is found anywhere in the chain.
func (e *Error) DeepIs(err error) bool {
	var cp = &Error{}
	*cp = *e

	var head *Error
	for {
		head = cp.Pop()
		if head == nil {
			break
		}
		if errors.Is(head.cause, err) {
			return true
		}
	}

	return false
}

// IsSP compares two SPErrors by their hash IDs.
// It returns true if both errors have the same hash ID.
// IsSP compares hashes of SpErrors, not their values or descriptions.
func (e *Error) IsSP(err *Error) bool {
	return slices.Compare(e.id.Sum(nil), err.id.Sum(nil)) == 0
}

func cmpHashes(h1, h2 hash.Hash) bool {
	return slices.Compare(h1.Sum(nil), h2.Sum(nil)) == 0
}
