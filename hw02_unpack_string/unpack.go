package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrInvalidDigit  = errors.New("invalid digit")
)

func Unpack(str string) (string, error) {
	var unpackedString []byte

	for index := 0; index < len(str); index++ {
		currentCharacter := str[index]
		currentIsDigit := unicode.IsDigit(rune(currentCharacter))

		if currentIsDigit {
			return "", ErrInvalidString
		}

		if index+1 == len(str) {
			unpackedString = append(unpackedString, currentCharacter)
			break
		}

		nextCharacter := str[index+1]
		nextIsDigit := unicode.IsDigit(rune(nextCharacter))

		if currentIsDigit && nextIsDigit {
			return "", ErrInvalidString
		}

		if nextIsDigit {
			digit, err := strconv.ParseInt(string(nextCharacter), 10, 32)
			if err != nil {
				return "", ErrInvalidDigit
			}

			for count := 0; int64(count) < digit; count++ {
				unpackedString = append(unpackedString, currentCharacter)
			}

			index++
			continue
		}
		unpackedString = append(unpackedString, currentCharacter)
	}

	return string(unpackedString), nil
}
