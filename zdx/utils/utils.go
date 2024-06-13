package utils

import (
	"encoding/json"
)

// Float64 is a custom type that forces JSON unmarshalling to interpret numbers as float64
type Float64 float64

func (f *Float64) UnmarshalJSON(b []byte) error {
	var number float64
	err := json.Unmarshal(b, &number)
	if err != nil {
		return err
	}
	*f = Float64(number)
	return nil
}
