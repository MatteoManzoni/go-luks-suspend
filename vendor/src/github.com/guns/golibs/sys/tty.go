// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build linux

package sys

import (
	"syscall"
	"unsafe"
)

// SetTTYIoctl identifies the three tcsetattr ioctls described in ioctl_tty(2)
type SetTTYIoctl uintptr

// These are typed for compile time safety.
// /usr/include/asm-generic/ioctls.h:
const (
	TCSETS  SetTTYIoctl = 0x5402
	TCSETSW             = 0x5403
	TCSETSF             = 0x5404
)

// GetTTYState writes the TTY state of fd to termios.
func GetTTYState(fd uintptr, termios *syscall.Termios) error {
	_, _, err := Ioctl(fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)))
	return err
}

// SetTTYState alters the TTY state of fd to match termios.
func SetTTYState(fd uintptr, action SetTTYIoctl, termios *syscall.Termios) error {
	// tcsetattr(3):
	// Note that tcsetattr() returns success if any of the requested changes
	// could be successfully carried out. Therefore, when making multiple changes
	// it may be necessary to follow this call with a further call to tcgetattr()
	// to check that all changes have been performed successfully.
	state := syscall.Termios{}
	for {
		_, _, err := Ioctl(fd, uintptr(action), uintptr(unsafe.Pointer(termios)))
		if err != nil {
			return err
		}

		if err := GetTTYState(fd, &state); err != nil {
			return err
		}

		if state == *termios {
			return nil
		}
	}
}

// AlterTTY changes the TTY indicated by fd to the termios struct returned by
// f, which receives the current TTY state. A function is returned that will
// return the TTY to its original state if it was altered. If the TTY was not
// altered, restoreTTY will be nil.
func AlterTTY(fd uintptr, action SetTTYIoctl, f func(syscall.Termios) syscall.Termios) (restoreTTY func() error, err error) {
	oldstate := syscall.Termios{}

	if err := GetTTYState(fd, &oldstate); err != nil {
		return nil, err
	}

	restoreTTY = func() error { return SetTTYState(fd, action, &oldstate) }
	newstate := f(oldstate)

	if err := SetTTYState(fd, action, &newstate); err != nil {
		return restoreTTY, err
	}

	return restoreTTY, nil
}
