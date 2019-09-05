// +build !windows,!linux,!darwin

package golog

func NewPrinter() *PlainPrinter {
	return NewPlainPrinter()
}
