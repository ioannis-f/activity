package propertylatitude

import (
	"fmt"
	float "github.com/go-fed/activity/streams/values/float"
	vocab "github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// LatitudeProperty is the functional property "latitude". It is permitted to be a
// single default-valued value type.
type LatitudeProperty struct {
	floatMember    float64
	hasFloatMember bool
	unknown        interface{}
	iri            *url.URL
	alias          string
}

// DeserializeLatitudeProperty creates a "latitude" property from an interface
// representation that has been unmarshalled from a text or binary format.
func DeserializeLatitudeProperty(m map[string]interface{}, aliasMap map[string]string) (*LatitudeProperty, error) {
	alias := ""
	if a, ok := aliasMap["https://www.w3.org/TR/activitystreams-vocabulary"]; ok {
		alias = a
	}
	propName := "latitude"
	if len(alias) > 0 {
		// Use alias both to find the property, and set within the property.
		propName = fmt.Sprintf("%s:%s", alias, "latitude")
	}
	i, ok := m[propName]

	if ok {
		if s, ok := i.(string); ok {
			u, err := url.Parse(s)
			// If error exists, don't error out -- skip this and treat as unknown string ([]byte) at worst
			// Also, if no scheme exists, don't treat it as a URL -- net/url is greedy
			if err == nil && len(u.Scheme) > 0 {
				this := &LatitudeProperty{
					alias: alias,
					iri:   u,
				}
				return this, nil
			}
		}
		if v, err := float.DeserializeFloat(i); err == nil {
			this := &LatitudeProperty{
				alias:          alias,
				floatMember:    v,
				hasFloatMember: true,
			}
			return this, nil
		}
		this := &LatitudeProperty{
			alias:   alias,
			unknown: i,
		}
		return this, nil
	}
	return nil, nil
}

// NewLatitudeProperty creates a new latitude property.
func NewLatitudeProperty() *LatitudeProperty {
	return &LatitudeProperty{alias: ""}
}

// Clear ensures no value of this property is set. Calling IsFloat afterwards will
// return false.
func (this *LatitudeProperty) Clear() {
	this.unknown = nil
	this.iri = nil
	this.hasFloatMember = false
}

// Get returns the value of this property. When IsFloat returns false, Get will
// return any arbitrary value.
func (this LatitudeProperty) Get() float64 {
	return this.floatMember
}

// GetIRI returns the IRI of this property. When IsIRI returns false, GetIRI will
// return any arbitrary value.
func (this LatitudeProperty) GetIRI() *url.URL {
	return this.iri
}

// HasAny returns true if the value or IRI is set.
func (this LatitudeProperty) HasAny() bool {
	return this.IsFloat() || this.iri != nil
}

// IsFloat returns true if this property is set and not an IRI.
func (this LatitudeProperty) IsFloat() bool {
	return this.hasFloatMember
}

// IsIRI returns true if this property is an IRI.
func (this LatitudeProperty) IsIRI() bool {
	return this.iri != nil
}

// JSONLDContext returns the JSONLD URIs required in the context string for this
// property and the specific values that are set. The value in the map is the
// alias used to import the property's value or values.
func (this LatitudeProperty) JSONLDContext() map[string]string {
	m := map[string]string{"https://www.w3.org/TR/activitystreams-vocabulary": this.alias}
	var child map[string]string

	/*
	   Since the literal maps in this function are determined at
	   code-generation time, this loop should not overwrite an existing key with a
	   new value.
	*/
	for k, v := range child {
		m[k] = v
	}
	return m
}

// KindIndex computes an arbitrary value for indexing this kind of value. This is
// a leaky API detail only for folks looking to replace the go-fed
// implementation. Applications should not use this method.
func (this LatitudeProperty) KindIndex() int {
	if this.IsFloat() {
		return 0
	}
	if this.IsIRI() {
		return -2
	}
	return -1
}

// LessThan compares two instances of this property with an arbitrary but stable
// comparison. Applications should not use this because it is only meant to
// help alternative implementations to go-fed to be able to normalize
// nonfunctional properties.
func (this LatitudeProperty) LessThan(o vocab.LatitudePropertyInterface) bool {
	// LessThan comparison for if either or both are IRIs.
	if this.IsIRI() && o.IsIRI() {
		return this.iri.String() < o.GetIRI().String()
	} else if this.IsIRI() {
		// IRIs are always less than other values, none, or unknowns
		return true
	} else if o.IsIRI() {
		// This other, none, or unknown value is always greater than IRIs
		return false
	}
	// LessThan comparison for the single value or unknown value.
	if !this.IsFloat() && !o.IsFloat() {
		// Both are unknowns.
		return false
	} else if this.IsFloat() && !o.IsFloat() {
		// Values are always greater than unknown values.
		return false
	} else if !this.IsFloat() && o.IsFloat() {
		// Unknowns are always less than known values.
		return true
	} else {
		// Actual comparison.
		return float.LessFloat(this.Get(), o.Get())
	}
}

// Name returns the name of this property: "latitude".
func (this LatitudeProperty) Name() string {
	return "latitude"
}

// Serialize converts this into an interface representation suitable for
// marshalling into a text or binary format. Applications should not need this
// function as most typical use cases serialize types instead of individual
// properties. It is exposed for alternatives to go-fed implementations to use.
func (this LatitudeProperty) Serialize() (interface{}, error) {
	if this.IsFloat() {
		return float.SerializeFloat(this.Get())
	} else if this.IsIRI() {
		return this.iri.String(), nil
	}
	return this.unknown, nil
}

// Set sets the value of this property. Calling IsFloat afterwards will return
// true.
func (this *LatitudeProperty) Set(v float64) {
	this.Clear()
	this.floatMember = v
	this.hasFloatMember = true
}

// SetIRI sets the value of this property. Calling IsIRI afterwards will return
// true.
func (this *LatitudeProperty) SetIRI(v *url.URL) {
	this.Clear()
	this.iri = v
}