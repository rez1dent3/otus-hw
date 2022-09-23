package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var ErrCantCastInt = errors.New("can't cast to int")

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

	ValidateEmpty struct {
		Code int `validate:""`
	}
)

func TestValidationErrorsCompare(t *testing.T) {
	tests := []struct {
		expected ValidationErrors
		actual   ValidationErrors
		result   bool
	}{
		{
			expected: ValidationErrors{},
			actual:   ValidationErrors{},
			result:   true,
		},
		{
			expected: ValidationErrors{
				{Field: "test1", Err: ErrIncorrectLength},
				{Field: "test", Err: ErrCantCastInt},
			},
			actual: ValidationErrors{
				{Field: "test", Err: ErrCantCastInt},
				{Field: "test1", Err: ErrIncorrectLength},
			},
			result: true,
		},
		{
			expected: ValidationErrors{
				{Field: "test1", Err: ErrIncorrectLength},
				{Field: "test", Err: ErrCantCastInt},
			},
			actual: ValidationErrors{
				{Field: "test1", Err: ErrIncorrectLength},
			},
			result: false,
		},
		{
			expected: ValidationErrors{
				{Field: "test1", Err: ErrIncorrectLength},
				{Field: "test", Err: ErrCantCastInt},
			},
			actual: nil,
			result: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			require.Equal(t, tt.result, errors.Is(tt.expected, tt.actual))
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          42,
			expectedErr: ErrNeedPassStruct,
		},
		{
			in:          &Token{},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          ValidateEmpty{},
			expectedErr: nil,
		},
		{
			in:          App{Version: "hello"},
			expectedErr: nil,
		},
		{
			in: App{Version: "hell"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrIncorrectLength},
			},
		},
		{
			in: &App{Version: "hell"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrIncorrectLength},
			},
		},
		{
			in: User{
				ID:     "",
				Email:  "",
				Phones: []string{""},
				Role:   "",
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrIncorrectLength},
				{Field: "Age", Err: ErrNumberIsLessThanMinimumAllowed},
				{Field: "Email", Err: ErrStringNotMatchRegExp},
				{Field: "Phones.0", Err: ErrIncorrectLength},
				{Field: "Role", Err: ErrStringIsNotIncludedInSet},
			},
		},
		{
			in: User{
				ID:    "",
				Age:   55,
				Email: "",
				Role:  "test",
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrIncorrectLength},
				{Field: "Age", Err: ErrNumberIsGreaterThanMaximumAllowed},
				{Field: "Email", Err: ErrStringNotMatchRegExp},
				{Field: "Role", Err: ErrStringIsNotIncludedInSet},
			},
		},
		{
			in: User{
				ID:    "",
				Age:   25,
				Email: "info@localhost.com",
				Role:  "admin",
				meta:  []byte(`{"target": "localhost"}`),
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrIncorrectLength},
			},
		},
		{
			in: User{
				ID:    "",
				Email: "info@localhost",
				Role:  "manager",
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrIncorrectLength},
				{Field: "Age", Err: ErrNumberIsLessThanMinimumAllowed},
				{Field: "Email", Err: ErrStringNotMatchRegExp},
				{Field: "Role", Err: ErrStringIsNotIncludedInSet},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "Hello World",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
				Body: "Hello World",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrNumberIsNotIncludedInSet},
			},
		},
		{
			in: Response{},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrNumberIsNotIncludedInSet},
			},
		},
	}

	validator := New()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			results := validator.Validate(tt.in)
			require.ErrorIs(t, tt.expectedErr, results)
		})
	}
}
