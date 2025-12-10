package gen

// FieldInfo holds metadata about a single struct field.
type FieldInfo struct {
	// Name is the original struct field name, e.g., "Host".
	Name string

	// Type is the original type of the field in the struct, e.g., "string", "*int".
	Type string

	// Func is the generated function name, e.g., "WithHost", "AddIP".
	Func string

	// ParamType is the type used for the generated function's parameter.
	// For slices, it's the element type. For pointers, it's the underlying type.
	ParamType string

	// Comments are the comments from the original struct field.
	Comments []string

	// DefaultVal is the default value as a Go code literal, from the `default` tag.
	DefaultVal string

	// IsSlice is true if the field is a slice.
	IsSlice bool

	// IsPointer is true if the field is a pointer.
	IsPointer bool

	// IsMap is true if the field is a map.
	IsMap bool

	// KeyType is the key type of the map.
	KeyType string

	// ValueType is the value type of the map.
	ValueType string
}

type TemplateData struct {
	Package    string
	Imports    []string
	StructName string
	Fields     []FieldInfo
}
