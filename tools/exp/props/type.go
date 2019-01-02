package props

import (
	"fmt"
	"github.com/cjslep/activity/tools/exp/codegen"
	"github.com/dave/jennifer/jen"
	"sort"
	"strings"
	"sync"
)

const (
	typeInterfaceName   = "Type"
	extendedByMethod    = "IsExtendedBy"
	extendingMethod     = "IsExtending"
	extendsMethod       = "Extends"
	disjointWithMethod  = "IsDisjointWith"
	typeNameMethod      = "Name"
	serializeMethodName = "Serialize"
	deserializeFnName   = "Deserialize"
	typeLessMethod      = "LessThan"
)

// TypeInterface returns the Type Interface that is needed for ActivityStream
// types to compile for methods dealing with extending, in the inheritance
// sense.
func TypeInterface(pkg Package) *codegen.Interface {
	comment := fmt.Sprintf("%s represents an ActivityStreams type.", typeInterfaceName)
	funcs := []codegen.FunctionSignature{
		{
			Name:    typeNameMethod,
			Params:  nil,
			Ret:     []jen.Code{jen.String()},
			Comment: fmt.Sprintf("%s returns the ActivityStreams type name.", typeNameMethod),
		},
	}
	return codegen.NewInterface(pkg.Path(), typeInterfaceName, funcs, comment)
}

// Property represents a property of an ActivityStreams type.
type Property interface {
	GetPublicPackage() Package
	PropertyName() string
	InterfaceName() string
	SetKindFns(name string, kind, deser *jen.Statement) error
	DeserializeFnName() string
}

// TypeGenerator represents an ActivityStream type definition to generate in Go.
type TypeGenerator struct {
	pm                *PackageManager
	typeName          string
	comment           string
	properties        map[string]Property
	withoutProperties map[string]Property
	rangeProperties   []Property
	extends           []*TypeGenerator
	disjoint          []*TypeGenerator
	extendedBy        []*TypeGenerator
	m                 *ManagerGenerator
	cacheOnce         sync.Once
	cachedStruct      *codegen.Struct
}

// NewTypeGenerator creates a new generator for a specific ActivityStreams Core
// or extension type. It will return an error if there are multiple properties
// have the same Name.
//
// The TypeGenerator should be in the second pass to construct, relying on the
// fact that properties have already been constructed.
//
// The extends and disjoint parameters are allowed to be nil. These lists must
// also have unique (non-duplicated) elements. Note that the disjoint entries
// will be set up bi-directionally properly; no need to go back to an existing
// TypeGenerator to set up the link correctly.
//
// The rangeProperties list is allowed to be nil. Any passed in will properly
// have their SetKindFns bookkeeping done.
//
// All TypeGenerators must be created before the Definition method is called, to
// ensure that type extension, in the inheritence sense, is properly set up.
//
// A ManagerGenerator must be created with this type before Definition is
// called, to ensure that the serialization functions are properly set up.
func NewTypeGenerator(pm *PackageManager, typeName, comment string,
	properties, withoutProperties, rangeProperties []Property,
	extends, disjoint []*TypeGenerator) (*TypeGenerator, error) {
	t := &TypeGenerator{
		pm:                pm,
		typeName:          typeName,
		comment:           comment,
		properties:        make(map[string]Property, len(properties)),
		withoutProperties: make(map[string]Property, len(withoutProperties)),
		rangeProperties:   rangeProperties,
		extends:           extends,
		disjoint:          disjoint,
	}
	for _, property := range properties {
		if _, has := t.properties[property.PropertyName()]; has {
			return nil, fmt.Errorf("type already has property with name %q", property.PropertyName())
		}
		t.properties[property.PropertyName()] = property
	}
	for _, wop := range withoutProperties {
		if _, has := t.withoutProperties[wop.PropertyName()]; has {
			return nil, fmt.Errorf("type already has withoutproperty with name %q", wop.PropertyName())
		}
		t.withoutProperties[wop.PropertyName()] = wop
	}
	// Complete doubly-linked extends/extendedBy lists.
	for _, ext := range extends {
		ext.extendedBy = append(ext.extendedBy, t)
	}
	// Complete doubly-linked disjoint types.
	for _, disj := range disjoint {
		disj.disjoint = append(disj.disjoint, t)
	}
	return t, nil
}

// apply propagates the manager's functions referring to this type's
// implementation as if this type were a Kind.
//
// Prepares to use the manager for the Definition generation.
func (t *TypeGenerator) apply(m *ManagerGenerator) error {
	t.m = m
	// Set up Kind functions
	// Note: this "i" must be the same as the "i" in the deserialization definition.
	// TODO: Remove this kluge.
	deser := m.getPrivateDeserializationMethodForType(t).On(managerInitName())
	kind := jen.Qual(t.PublicPackage().Path(), t.InterfaceName())
	for _, p := range t.rangeProperties {
		if e := p.SetKindFns(t.TypeName(), kind, deser); e != nil {
			return e
		}
	}
	return nil
}

// Package gets this TypeGenerator's Private Package.
func (t *TypeGenerator) PrivatePackage() Package {
	return t.pm.PrivatePackage()
}

// Package gets this TypeGenerator's Public Package.
func (t *TypeGenerator) PublicPackage() Package {
	return t.pm.PublicPackage()
}

// Comment returns the comment for this type.
func (t *TypeGenerator) Comments() string {
	return t.comment
}

// TypeName returns the ActivityStreams name for this type.
func (t *TypeGenerator) TypeName() string {
	return t.typeName
}

// InterfaceName returns the interface name for this type.
func (t *TypeGenerator) InterfaceName() string {
	return fmt.Sprintf("%sInterface", t.TypeName())
}

// Extends returns the generators of types that this ActivityStreams type
// extends from.
func (t *TypeGenerator) Extends() []*TypeGenerator {
	return t.extends
}

// ExtendedBy returns the generators of types that extend from this
// ActivityStreams type.
func (t *TypeGenerator) ExtendedBy() []*TypeGenerator {
	return t.extendedBy
}

// Disjoint returns the generators of types that this ActivityStreams type is
// disjoint to.
func (t *TypeGenerator) Disjoint() []*TypeGenerator {
	return t.disjoint
}

// Properties returns the Properties of this type, mapped by their property
// name.
func (t *TypeGenerator) Properties() map[string]Property {
	return t.properties
}

// WithoutProperties returns the properties that do not apply to this type,
// mapped by their property name.
func (t *TypeGenerator) WithoutProperties() map[string]Property {
	return t.withoutProperties
}

// extendsFnName determines the name of the Extends function, which
// determines if this ActivityStreams type extends another one.
func (t *TypeGenerator) extendsFnName() string {
	return fmt.Sprintf("%s%s", t.TypeName(), extendsMethod)
}

// extendedByFnName determines the name of the ExtendedBy function, which
// determines if another ActivityStreams type extends this one.
func (t *TypeGenerator) extendedByFnName() string {
	return fmt.Sprintf("%s%s", t.TypeName(), extendedByMethod)
}

// disjointWithFnName determines the name of the DisjointWith function, which
// determines if another ActivityStreams type is disjoint with this one.
func (t *TypeGenerator) disjointWithFnName() string {
	return fmt.Sprintf("%s%s", t.TypeName(), disjointWithMethod)
}

// deserializationFnName determines the name of the deserialize function for
// this type.
func (t *TypeGenerator) deserializationFnName() string {
	return fmt.Sprintf("%s%s", deserializeFnName, t.TypeName())
}

// toInterface creates the interface version of the definition generated.
//
// Requires apply to have already been called.
//
// TODO: Delete -- redundant?
func (t *TypeGenerator) toInterface(pkg Package) *codegen.Interface {
	s := t.Definition()
	return s.ToInterface(pkg.Path(), t.InterfaceName(), "")
}

// InterfaceDefinition creates the interface of this type in the specified
// package.
//
// Requires ManagerGenerator to have been created.
func (t *TypeGenerator) InterfaceDefinition(pkg Package) *codegen.Interface {
	s := t.Definition()
	return s.ToInterface(pkg.Path(), t.InterfaceName(), t.Comments())
}

// Definition generates the golang code for this ActivityStreams type.
func (t *TypeGenerator) Definition() *codegen.Struct {
	t.cacheOnce.Do(func() {
		members := t.members()
		ser := t.serializationMethod()
		less := t.lessMethod()
		deser := t.kindDeserializationFunc()
		extendsFn, extendsMethod := t.extendsDefinition()
		t.cachedStruct = codegen.NewStruct(
			jen.Commentf(t.Comments()),
			t.TypeName(),
			[]*codegen.Method{
				t.nameDefinition(),
				extendsMethod,
				ser,
				less,
			},
			[]*codegen.Function{
				t.extendedByDefinition(),
				extendsFn,
				t.disjointWithDefinition(),
				deser,
			},
			members)
	})
	return t.cachedStruct
}

func (t *TypeGenerator) allProperties() map[string]Property {
	p := t.properties
	// Properties of parents that are extended, minus DoesNotApplyTo
	var extends []*TypeGenerator
	extends = t.getAllParentExtends(extends, t)
	for _, ext := range t.extends {
		for k, v := range ext.Properties() {
			p[k] = v
		}
	}
	for _, ext := range t.extends {
		for k, _ := range ext.WithoutProperties() {
			delete(p, k)
		}
	}
	return p
}

// sortedProperty is a slice of Properties that implements the Sort interface.
type sortedProperty []Property

func (s sortedProperty) Less(i, j int) bool {
	return s[i].PropertyName() < s[j].PropertyName()
}

func (s sortedProperty) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortedProperty) Len() int {
	return len(s)
}

func (t *TypeGenerator) members() (members []jen.Code) {
	p := t.allProperties()
	// Sort the properties for readability
	sortedMembers := make(sortedProperty, 0, len(p))
	for _, property := range p {
		sortedMembers = append(sortedMembers, property)
	}
	sort.Sort(sortedMembers)
	// Convert to jen.Code
	members = make([]jen.Code, 0, len(p))
	for _, property := range sortedMembers {
		members = append(members, jen.Id(strings.Title(property.PropertyName())).Qual(property.GetPublicPackage().Path(), property.InterfaceName()))
	}
	return
}

// nameDefinition generates the golang method for returning the ActivityStreams
// type name.
func (t *TypeGenerator) nameDefinition() *codegen.Method {
	return codegen.NewCommentedValueMethod(
		t.PrivatePackage().Path(),
		typeNameMethod,
		t.TypeName(),
		/*params=*/ nil,
		[]jen.Code{jen.String()},
		[]jen.Code{
			jen.Return(jen.Lit(t.TypeName())),
		},
		jen.Commentf("%s returns the name of this type.", typeNameMethod))
}

// getAllParentExtends recursively determines all the parent types that this
// type extends from.
func (t *TypeGenerator) getAllParentExtends(s []*TypeGenerator, tg *TypeGenerator) []*TypeGenerator {
	for _, e := range tg.Extends() {
		s = append(s, e)
		s = append(s, t.getAllParentExtends(s, e)...)
	}
	return s
}

// extendsDefinition generates the golang function for determining if this
// ActivityStreams type extends another type. It requires the Type interface.
func (t *TypeGenerator) extendsDefinition() (*codegen.Function, *codegen.Method) {
	var extends []*TypeGenerator
	extends = t.getAllParentExtends(extends, t)
	extendNames := make(map[string]struct{}, len(extends))
	for _, ext := range extends {
		extendNames[ext.TypeName()] = struct{}{}
	}
	extensions := make([]jen.Code, len(extendNames))
	for e := range extendNames {
		extensions = append(extensions, jen.Lit(e))
	}
	impl := []jen.Code{jen.Comment("Shortcut implementation: this does not extend anything."), jen.Return(jen.False())}
	if len(extensions) > 0 {
		impl = []jen.Code{jen.Id("extensions").Op(":=").Index().String().Values(extensions...),
			jen.For(jen.List(
				jen.Id("_"),
				jen.Id("ext"),
			).Op(":=").Range().Id("extensions")).Block(
				jen.If(
					jen.Id("ext").Op("==").Id("other").Dot(typeNameMethod).Call(),
				).Block(
					jen.Return(jen.True()),
				),
			),
			jen.Return(jen.False())}
	}
	f := codegen.NewCommentedFunction(
		t.PrivatePackage().Path(),
		t.extendsFnName(),
		[]jen.Code{jen.Id("other").Id(typeInterfaceName)},
		[]jen.Code{jen.Bool()},
		impl,
		jen.Commentf("%s returns true if the %s type extends from the other type.", t.extendsFnName(), t.TypeName()))
	m := codegen.NewCommentedValueMethod(
		t.PrivatePackage().Path(),
		extendingMethod,
		t.TypeName(),
		[]jen.Code{jen.Id("other").Id(typeInterfaceName)},
		[]jen.Code{jen.Bool()},
		[]jen.Code{
			jen.Return(
				jen.Id(t.extendsFnName()).Call(jen.Id("other")),
			),
		},
		jen.Commentf("%s returns true if the %s type extends from the other type.", extendingMethod, t.TypeName()))
	return f, m
}

// getAllChildrenExtendBy recursivley determines all the child types that this
// type is extended by.
func (t *TypeGenerator) getAllChildrenExtendedBy(s []string, tg *TypeGenerator) {
	for _, e := range tg.ExtendedBy() {
		s = append(s, e.TypeName())
		t.getAllChildrenExtendedBy(s, e)
	}
}

// extendedByDefinition generates the golang function for determining if
// another ActivityStreams type extends this type. It requires the Type
// interface.
func (t *TypeGenerator) extendedByDefinition() *codegen.Function {
	var extendNames []string
	t.getAllChildrenExtendedBy(extendNames, t)
	extensions := make([]jen.Code, len(extendNames))
	for i, e := range extendNames {
		extensions[i] = jen.Lit(e)
	}
	impl := []jen.Code{jen.Comment("Shortcut implementation: is not extended by anything."), jen.Return(jen.False())}
	if len(extensions) > 0 {
		impl = []jen.Code{jen.Id("extensions").Op(":=").Index().String().Values(extensions...),
			jen.For(jen.List(
				jen.Id("_"),
				jen.Id("ext"),
			).Op(":=").Range().Id("extensions")).Block(
				jen.If(
					jen.Id("ext").Op("==").Id("other").Dot(typeNameMethod).Call(),
				).Block(
					jen.Return(jen.True()),
				),
			),
			jen.Return(jen.False())}
	}
	return codegen.NewCommentedFunction(
		t.PrivatePackage().Path(),
		t.extendedByFnName(),
		[]jen.Code{jen.Id("other").Id(typeInterfaceName)},
		[]jen.Code{jen.Bool()},
		impl,
		jen.Commentf("%s returns true if the other provided type extends from the %s type.", t.extendedByFnName(), t.TypeName()))
}

// getAllChildrenDisjointWith recursivley determines all the child types that this
// type is disjoint with.
func (t *TypeGenerator) getAllDisjointWith(s []string) {
	for _, e := range t.Disjoint() {
		s = append(s, e.TypeName())
		// Get all the disjoint type's children.
		t.getAllChildrenExtendedBy(s, e)
	}
}

// disjointWithDefinition generates the golang function for determining if
// another ActivityStreams type is disjoint with this type. It requires the Type
// interface.
func (t *TypeGenerator) disjointWithDefinition() *codegen.Function {
	// TODO: Inherit disjoint from parent and the other extended types of
	// the other.
	var disjointNames []string
	t.getAllDisjointWith(disjointNames)
	disjointWith := make([]jen.Code, len(disjointNames))
	for i, d := range disjointNames {
		disjointWith[i] = jen.Lit(d)
	}
	impl := []jen.Code{jen.Comment("Shortcut implementation: is not disjoint with anything."), jen.Return(jen.False())}
	if len(disjointWith) > 0 {
		impl = []jen.Code{jen.Id("disjointWith").Op(":=").Index().String().Values(disjointWith...),
			jen.For(jen.List(
				jen.Id("_"),
				jen.Id("disjoint"),
			).Op(":=").Range().Id("disjointWith")).Block(
				jen.If(
					jen.Id("disjoint").Op("==").Id("other").Dot(typeNameMethod).Call(),
				).Block(
					jen.Return(jen.True()),
				),
			),
			jen.Return(jen.False())}
	}
	return codegen.NewCommentedFunction(
		t.PrivatePackage().Path(),
		t.disjointWithFnName(),
		[]jen.Code{jen.Id("other").Id(typeInterfaceName)},
		[]jen.Code{jen.Bool()},
		impl,
		jen.Commentf("%s returns true if the other provided type is disjoint with the %s type.", t.disjointWithFnName(), t.TypeName()))
}

// serializationMethod returns the method needed to serialize a TypeGenerator as
// a property.
func (t *TypeGenerator) serializationMethod() (ser *codegen.Method) {
	ser = codegen.NewCommentedValueMethod(
		t.PrivatePackage().Path(),
		serializeMethodName,
		t.TypeName(),
		/*params=*/ nil,
		[]jen.Code{jen.Interface(), jen.Error()},
		[]jen.Code{
			// TODO
			jen.Commentf("TODO: Serialization code for %s", t.TypeName()),
		},
		jen.Commentf("%s converts this into an interface representation suitable for marshalling into a text or binary format.", serializeMethodName))
	return
}

// lessMethod returns the method needed to compare a type with another type.
func (t *TypeGenerator) lessMethod() (less *codegen.Method) {
	less = codegen.NewCommentedValueMethod(
		t.PrivatePackage().Path(),
		typeLessMethod,
		t.TypeName(),
		[]jen.Code{
			jen.Id("o").Qual(t.PublicPackage().Path(), t.InterfaceName()),
		},
		[]jen.Code{jen.Bool()},
		[]jen.Code{
			// TODO
			jen.Commentf("TODO: Less code for %s", t.TypeName()),
		},
		jen.Commentf("%s computes if this %s is lesser, with an arbitrary but stable determination", typeLessMethod, t.TypeName()))
	return
}

// kindDeserializationFunc returns free function reference that can be used to
// treat a TypeGenerator as another property's Kind.
func (t *TypeGenerator) kindDeserializationFunc() (deser *codegen.Function) {
	deserCode := jen.Empty()
	for name, prop := range t.allProperties() {
		deserMethod := t.m.getPrivateDeserializationMethodForProperty(prop)
		deserCode = deserCode.Add(
			jen.If(
				jen.List(
					jen.Id("p"),
					jen.Err(),
				// TODO: Ensure this variable is called correctly
				).Op(":=").Add(deserMethod.Call(managerInitName(), jen.Id("m"))),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Nil(), jen.Err()),
			).Else().Block(
				jen.Id(codegen.This()).Dot(strings.Title(name)).Op("=").Op("*").Id("p"),
			).Line())
	}
	deser = codegen.NewCommentedFunction(
		t.PrivatePackage().Path(),
		t.deserializationFnName(),
		[]jen.Code{jen.Id("m").Map(jen.String()).Interface()},
		[]jen.Code{jen.Op("*").Id(t.TypeName()), jen.Error()},
		[]jen.Code{
			jen.Id(codegen.This()).Op(":=").Op("&").Id(t.TypeName()).Values(),
			deserCode,
			jen.Return(jen.Id(codegen.This()), jen.Nil()),
		},
		jen.Commentf("%s creates a %s from a map representation that has been unmarshalled from a text or binary format.", t.deserializationFnName(), t.TypeName()))
	return
}
