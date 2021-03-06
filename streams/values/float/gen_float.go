package float

import "fmt"

// SerializeFloat converts a float value to an interface representation suitable
// for marshalling into a text or binary format.
func SerializeFloat(this float64) (interface{}, error) {
	return this, nil
}

// DeserializeFloat creates float value from an interface representation that has
// been unmarshalled from a text or binary format.
func DeserializeFloat(this interface{}) (float64, error) {
	if f, ok := this.(float64); ok {
		return f, nil
	} else {
		return 0, fmt.Errorf("%v cannot be interpreted as a float64 for xsd:float", this)
	}
}

// LessFloat returns true if the left float value is less than the right value.
func LessFloat(lhs, rhs float64) bool {
	return lhs < rhs
}
