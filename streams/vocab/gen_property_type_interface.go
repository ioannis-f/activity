package vocab

import "net/url"

// TypePropertyIteratorInterface represents a single value for the "type" property.
type TypePropertyIteratorInterface interface {
	// GetAnyURI returns the value of this property. When IsAnyURI returns
	// false, GetAnyURI will return an arbitrary value.
	GetAnyURI() *url.URL
	// GetIRI returns the IRI of this property. When IsIRI returns false,
	// GetIRI will return an arbitrary value.
	GetIRI() *url.URL
	// GetString returns the value of this property. When IsString returns
	// false, GetString will return an arbitrary value.
	GetString() string
	// HasAny returns true if any of the different values is set.
	HasAny() bool
	// IsAnyURI returns true if this property has a type of "anyURI". When
	// true, use the GetAnyURI and SetAnyURI methods to access and set
	// this property.
	IsAnyURI() bool
	// IsIRI returns true if this property is an IRI. When true, use GetIRI
	// and SetIRI to access and set this property
	IsIRI() bool
	// IsString returns true if this property has a type of "string". When
	// true, use the GetString and SetString methods to access and set
	// this property.
	IsString() bool
	// JSONLDContext returns the JSONLD URIs required in the context string
	// for this property and the specific values that are set. The value
	// in the map is the alias used to import the property's value or
	// values.
	JSONLDContext() map[string]string
	// KindIndex computes an arbitrary value for indexing this kind of value.
	// This is a leaky API detail only for folks looking to replace the
	// go-fed implementation. Applications should not use this method.
	KindIndex() int
	// LessThan compares two instances of this property with an arbitrary but
	// stable comparison. Applications should not use this because it is
	// only meant to help alternative implementations to go-fed to be able
	// to normalize nonfunctional properties.
	LessThan(o TypePropertyIteratorInterface) bool
	// Name returns the name of this property: "type".
	Name() string
	// Next returns the next iterator, or nil if there is no next iterator.
	Next() TypePropertyIteratorInterface
	// Prev returns the previous iterator, or nil if there is no previous
	// iterator.
	Prev() TypePropertyIteratorInterface
	// SetAnyURI sets the value of this property. Calling IsAnyURI afterwards
	// returns true.
	SetAnyURI(v *url.URL)
	// SetIRI sets the value of this property. Calling IsIRI afterwards
	// returns true.
	SetIRI(v *url.URL)
	// SetString sets the value of this property. Calling IsString afterwards
	// returns true.
	SetString(v string)
}

// Identifies the Object or Link type. Multiple values may be specified.
//
// Example 62 (https://www.w3.org/TR/activitystreams-vocabulary/#extype-jsonld):
//   {
//     "summary": "A foo",
//     "type": "http://example.org/Foo"
//   }
type TypePropertyInterface interface {
	// AppendAnyURI appends a anyURI value to the back of a list of the
	// property "type". Invalidates iterators that are traversing using
	// Prev.
	AppendAnyURI(v *url.URL)
	// AppendIRI appends an IRI value to the back of a list of the property
	// "type"
	AppendIRI(v *url.URL)
	// AppendString appends a string value to the back of a list of the
	// property "type". Invalidates iterators that are traversing using
	// Prev.
	AppendString(v string)
	// At returns the property value for the specified index. Panics if the
	// index is out of bounds.
	At(index int) TypePropertyIteratorInterface
	// Begin returns the first iterator, or nil if empty. Can be used with the
	// iterator's Next method and this property's End method to iterate
	// from front to back through all values.
	Begin() TypePropertyIteratorInterface
	// Empty returns returns true if there are no elements.
	Empty() bool
	// End returns beyond-the-last iterator, which is nil. Can be used with
	// the iterator's Next method and this property's Begin method to
	// iterate from front to back through all values.
	End() TypePropertyIteratorInterface
	// JSONLDContext returns the JSONLD URIs required in the context string
	// for this property and the specific values that are set. The value
	// in the map is the alias used to import the property's value or
	// values.
	JSONLDContext() map[string]string
	// KindIndex computes an arbitrary value for indexing this kind of value.
	// This is a leaky API method specifically needed only for alternate
	// implementations for go-fed. Applications should not use this
	// method. Panics if the index is out of bounds.
	KindIndex(idx int) int
	// Len returns the number of values that exist for the "type" property.
	Len() (length int)
	// Less computes whether another property is less than this one. Mixing
	// types results in a consistent but arbitrary ordering
	Less(i, j int) bool
	// LessThan compares two instances of this property with an arbitrary but
	// stable comparison. Applications should not use this because it is
	// only meant to help alternative implementations to go-fed to be able
	// to normalize nonfunctional properties.
	LessThan(o TypePropertyInterface) bool
	// Name returns the name of this property: "type".
	Name() string
	// PrependAnyURI prepends a anyURI value to the front of a list of the
	// property "type". Invalidates all iterators.
	PrependAnyURI(v *url.URL)
	// PrependIRI prepends an IRI value to the front of a list of the property
	// "type".
	PrependIRI(v *url.URL)
	// PrependString prepends a string value to the front of a list of the
	// property "type". Invalidates all iterators.
	PrependString(v string)
	// Remove deletes an element at the specified index from a list of the
	// property "type", regardless of its type. Panics if the index is out
	// of bounds. Invalidates all iterators.
	Remove(idx int)
	// Serialize converts this into an interface representation suitable for
	// marshalling into a text or binary format. Applications should not
	// need this function as most typical use cases serialize types
	// instead of individual properties. It is exposed for alternatives to
	// go-fed implementations to use.
	Serialize() (interface{}, error)
	// SetAnyURI sets a anyURI value to be at the specified index for the
	// property "type". Panics if the index is out of bounds. Invalidates
	// all iterators.
	SetAnyURI(idx int, v *url.URL)
	// SetIRI sets an IRI value to be at the specified index for the property
	// "type". Panics if the index is out of bounds.
	SetIRI(idx int, v *url.URL)
	// SetString sets a string value to be at the specified index for the
	// property "type". Panics if the index is out of bounds. Invalidates
	// all iterators.
	SetString(idx int, v string)
	// Swap swaps the location of values at two indices for the "type"
	// property.
	Swap(i, j int)
}