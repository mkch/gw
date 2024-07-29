package sysutil

// Package sysutil implements utility functions for win32 syscall.
import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// IsNoError returns whether err matches(errors.Is) syscall.Error(0).
// This usually means no error during calling a syscall.
func IsNoError(err error) bool {
	return errors.Is(err, syscall.Errno(0))
}

// MustZero returns nil if r1 is 0, err otherwise.
func MustZero(r1 uintptr, r2 uintptr, err error) error {
	if r1 == 0 {
		return nil
	}
	return err
}

// MustNotZero returns (r1, nil) if r1 is not zero, (r1, err) otherwise.
func MustNotZero[T constraints.Integer | ~unsafe.Pointer](r1 uintptr, r2 uintptr, err error) (T, error) {
	if r1 != 0 {
		return T(r1), nil
	}
	return T(r1), err
}

// MustNotNegativeOne returns (r1, nil) if r1 is not -1, (r1, err) otherwise.
func MustNotNegativeOne[T constraints.Signed](r1 uintptr, r2 uintptr, err error) (T, error) {
	if T(r1) != -1 {
		return T(r1), nil
	}
	return T(r1), err
}

// MustNoError returns (r1, nil) if err is(IsNoError) something other than syscall.Errno(0), r1, err otherwise.
func MustNoError[T constraints.Integer](r1 uintptr, r2 uintptr, err error) (T, error) {
	if IsNoError(err) {
		return T(r1), nil
	}
	return T(r1), err
}

// MustTrue returns nil if r1 is not 0, err otherwise.
func MustTrue(r1 uintptr, r2 uintptr, err error) error {
	if r1 != 0 {
		return nil
	}
	return err
}

// As returns T(r1) and ignores err.
func As[T constraints.Integer](r1 uintptr, r2 uintptr, err error) T {
	return T(r1)
}

// AsBool returns whether r1 is not zero and ignores err.
func AsBool(r1 uintptr, r2 uintptr, err error) bool {
	return r1 != 0
}
