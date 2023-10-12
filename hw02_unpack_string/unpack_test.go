package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abcd", expected: "abcd"},
		{input: "abccd", expected: "abccd"},
		{input: "aaa0b", expected: "aab"},
		{input: "", expected: ""},
		{input: "b2c1r3l0", expected: "bbcrrr"},
		{input: "при4ветт0", expected: "приииивет"},
		{input: "你2好5", expected: "你你好好好好好"},
		{input: "🥺5смайл2🥺1", expected: "🥺🥺🥺🥺🥺смайлл🥺"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
