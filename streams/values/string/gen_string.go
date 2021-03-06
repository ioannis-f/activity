package string

import "fmt"

// SerializeString converts a string value to an interface representation suitable
// for marshalling into a text or binary format.
func SerializeString(this string) (interface{}, error) {
	return this, nil
}

// DeserializeString creates string value from an interface representation that
// has been unmarshalled from a text or binary format.
func DeserializeString(this interface{}) (string, error) {
	if s, ok := this.(string); ok {
		return s, nil
	} else {
		return "", fmt.Errorf("%v cannot be interpreted as a string for xsd:string", this)
	}
}

// LessString returns true if the left string value is less than the right value.
func LessString(lhs, rhs string) bool {
	return lhs < rhs
}
