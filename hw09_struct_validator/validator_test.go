package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}
	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Dmitriy Morozov",
				Age:    50,
				Email:  "info@otus.com",
				Role:   "admin",
				Phones: []string{"89224818241"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Dmitriy Morozov",
				Age:    50,
				Email:  "infootus.com",
				Role:   "admin",
				Phones: []string{"89224818241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrInvalidRegexp,
				},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Dmitriy Morozov",
				Age:    50,
				Email:  "info@otus.com",
				Role:   "manager",
				Phones: []string{"89224818241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   ErrInvalidIn,
				},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Dmitriy Morozov",
				Age:    1,
				Email:  "info@otus.com",
				Role:   "admin",
				Phones: []string{"89224818241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrInvalidMin,
				},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Dmitriy Morozov",
				Age:    120,
				Email:  "info@otus.com",
				Role:   "admin",
				Phones: []string{"89224818241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrInvalidMax,
				},
			},
		},
		{
			in: User{
				ID:     "12345678901234567890123456789012345", // ID length is not 36
				Name:   "Denis Morozov",
				Age:    18,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"89224818241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrInvalidLen,
				},
			},
		},
		{
			in: User{
				ID:     "12345678901234567890123456789012345", // ID length is not 36
				Name:   "Denis Morozov",
				Age:    18,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"892", "89224818241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrInvalidLen,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrInvalidLen,
				},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrInvalidLen,
				},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 300,
				Body: "OK",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrInvalidIn,
				},
			},
		},
	}

	for i, testCase := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			testCase := testCase
			t.Parallel()

			err := Validate(testCase.in)
			require.Equal(t, testCase.expectedErr, err)
		})
	}
}
