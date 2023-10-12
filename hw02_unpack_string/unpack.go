package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrInvalidDigit  = errors.New("invalid digit")
)

func Unpack(str string) (string, error) {
	if len(str) < 1 {
		return "", nil
	}

	var unpacked strings.Builder
	runes := []rune(str)

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	for i, symbol := range runes {
		if unicode.IsDigit(runes[i]) {
			if i > 1 && unicode.IsDigit(runes[i-1]) {
				return "", ErrInvalidString
			}
			continue
		}

		if i+1 < len(runes) {
			nextSymbol := runes[i+1]

			if unicode.IsDigit(symbol) && unicode.IsDigit(nextSymbol) {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(nextSymbol) {
				digit, err := strconv.ParseInt(string(nextSymbol), 10, 32)
				if err != nil {
					return "", ErrInvalidDigit
				}
				repeated := strings.Repeat(string(symbol), int(digit))
				unpacked.WriteString(repeated)
				continue
			}
		}

		unpacked.WriteString(string(symbol))
	}

	return unpacked.String(), nil
}
