package oututil

import (
	"regexp"
	"strconv"
)

// Direction impies the text direction
type Direction byte

const (
	// LTR is left to right
	LTR Direction = iota
	// RTL is right to left
	RTL
)

// SourceError represents an errror in the srouce code
type SourceError struct {
	File         string
	Line, Column int
	Message      string
}

// Point is a coordinate in text with a length and direction
type Point struct {
	Origin int64
	Length int64
	Dir    Direction
}

// ScanSourceError takes the log of a process and
// returns it's sourcecode errors
func ScanSourceError(message string) []SourceError {
	/*
	   Parse a message that is almost the standard in error messages that are outputed by
	   most modern compilers and tools that work with source code

	   format is either

	   xxx.yyy:01:01: some message
	   {filename}.{fileext}:{line}:{column}: {message}

	   or
	   xxx.yyy:01: some message
	   {filename}.{fileext}:{line}: {message}

	*/
	var validMessage = regexp.MustCompile(`([[:alnum:]]+.[[:alnum:]]+):([0-9]+):([0-9]+)?:? (.*)`)
	var errors []SourceError
	messages := validMessage.FindAllStringSubmatch(message, -1)
	for _, message := range messages {
		e := SourceError{
			File: message[1],
		}
		if line, err := strconv.Atoi(message[2]); err == nil {
			e.Line = line
		}
		if column, err := strconv.Atoi(message[3]); err == nil {
			e.Column = column
		} else {
			e.Column = -1
			e.Message = message[3]
		}
		if len(message) > 3 {
			e.Message = message[4]
		}
		errors = append(errors, e)
	}
	return errors
}
