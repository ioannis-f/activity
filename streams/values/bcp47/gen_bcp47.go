package bcp47

import "fmt"

// SerializeBcp47 converts a bcp47 value to an interface representation suitable
// for marshalling into a text or binary format.
func SerializeBcp47(this string) (interface{}, error) {
	return this, nil
}

// DeserializeBcp47 creates bcp47 value from an interface representation that has
// been unmarshalled from a text or binary format.
func DeserializeBcp47(this interface{}) (string, error) {
	if s, ok := this.(string); ok {
		return s, nil
	} else {
		return "", fmt.Errorf("%v cannot be interpreted as a string for bcp47 languagetag", this)
	}
}

// LessBcp47 returns true if the left bcp47 value is less than the right value.
func LessBcp47(lhs, rhs string) bool {
	return lhs < rhs
}
