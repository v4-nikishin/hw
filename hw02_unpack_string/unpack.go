package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func prev(s string, idx int) rune {
	var prev rune
	for i, r := range s {
		if i == idx {
			break
		}
		prev = r
	}
	return prev
}

func isLast(s string, idx int, r rune) bool {
	return len(s) == idx+len(string(r))
}

func Unpack(s string) (string, error) {
	var b strings.Builder

	for i, r := range s {
		if i == 0 {
			if unicode.IsDigit(r) {
				return "", ErrInvalidString
			}
			if len(s) == len(string(r)) {
				return string(r), nil
			}
		} else {
			if unicode.IsDigit(r) {
				if unicode.IsDigit(prev(s, i)) {
					return "", ErrInvalidString
				}
				c, _ := strconv.Atoi(string(r))
				if c != 0 {
					b.WriteString(strings.Repeat(string(prev(s, i)), c))
				}
			} else {
				if !unicode.IsDigit(prev(s, i)) {
					b.WriteRune(prev(s, i))
				}
				if isLast(s, i, r) {
					b.WriteRune(r)
				}
			}
		}
	}
	result := b.String()
	return result, nil
}
