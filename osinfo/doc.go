// Package osinfo exposes typed helpers describing the operating system
// the current process is running on.
//
// Every other package in onyx asks osinfo instead of switching on
// runtime.GOOS directly. This guarantees that platform facts have a single
// source of truth.
package osinfo
