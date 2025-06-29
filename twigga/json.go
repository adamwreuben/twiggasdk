package twigga

import (
	"fmt"

	"github.com/bytedance/sonic"
)

// Marshal efficiently converts a Go value to a JSON []byte.
func Marshal(data interface{}) ([]byte, error) {
	b, err := sonic.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return b, nil
}

// MustMarshal is like Marshal but panics if there's an error.
func MustMarshal(data interface{}) []byte {
	bytes, err := Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes
}

// Unmarshal takes a JSON byte slice and a pointer to a variable, and fills the variable with parsed data.
func Unmarshal(data []byte, out interface{}) error {
	if err := sonic.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// MustUnmarshal is like Unmarshal but panics if there's an error.
func MustUnmarshal(data []byte, out interface{}) {
	if err := Unmarshal(data, out); err != nil {
		panic(err)
	}
}
