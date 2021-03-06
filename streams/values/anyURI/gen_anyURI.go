package anyuri

import (
	"fmt"
	"net/url"
)

// SerializeAnyURI converts a anyURI value to an interface representation suitable
// for marshalling into a text or binary format.
func SerializeAnyURI(this *url.URL) (interface{}, error) {
	return this.String(), nil
}

// DeserializeAnyURI creates anyURI value from an interface representation that
// has been unmarshalled from a text or binary format.
func DeserializeAnyURI(this interface{}) (*url.URL, error) {
	var u *url.URL
	var err error
	if s, ok := this.(string); ok {
		u, err = url.Parse(s)
		if err != nil {
			err = fmt.Errorf("%v cannot be interpreted as a xsd:anyURI: %s", this, err)
		} else if len(u.Scheme) == 0 {
			err = fmt.Errorf("%v cannot be interpreted as a xsd:anyURI: no scheme", this)
		}
	} else {
		err = fmt.Errorf("%v cannot be interpreted as a string for xsd:anyURI", this)
	}
	return u, err
}

// LessAnyURI returns true if the left anyURI value is less than the right value.
func LessAnyURI(lhs, rhs *url.URL) bool {
	return lhs.String() < rhs.String()
}
