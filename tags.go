// Package tags parses Go struct tags and splits each value into a name and options.
//
// A value is split with the separators , ; and :
//
//   - the name is the part before the first , or ; and may hold spaces and non-ASCII letters (First Name, Año);
//   - options follow, separated by , or ; and are either a key (omitempty) or a key:value pair (default:anonymous);
//   - a value holding a : before any , or ; has no name, only options (column:Last Name).
//
// A separator belongs to a value when it is not followed by a letter, a digit or an underscore:
// split:,;| gives the value ",;|". The constraint option keeps its key:value pairs as a single value:
// constraint:OnUpdate:CASCADE,OnDelete:SET NULL gives the value "OnUpdate:CASCADE,OnDelete:SET NULL".
//
// A name cannot hold , ; or :
//
// Example usage:
//
//	tgs, _ := tags.Parse(`json:"id,omitempty" excel:"First Name,default:anonymous"`)
//	tgs[1].Name                          // "First Name"
//	tgs[1].GetOption("default").Value    // "anonymous"
//
//	field, _ := reflect.TypeOf(User{}).FieldByName("Name")
//	tag := tags.Lookup(field, "excel")   // nil when the field has no excel tag
package tags

// Tag is one key:"value" pair of a struct tag, see https://pkg.go.dev/reflect#StructTag
type Tag struct {
	// Tag is the whole pair: json:"id,omitempty"
	Tag string

	// Key is the tag key, as given to StructTag.Get
	// in `json:"id,omitempty"`, the Key is "json"
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Key is "gorm"
	Key string

	// Value is the tag value, as returned by StructTag.Get
	// in `json:"id,omitempty"`, the Value is "id,omitempty"
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Value is "embedded;embeddedPrefix:author_"
	Value string

	// Name is the first part of the value, before the first , or ;
	// in `json:"id,omitempty"`, the Name is "id"
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Name is "embedded"
	// in `excel:"column:Last Name"`, the Name is "" (the value holds options only)
	Name string

	// Options are the parts of the value after the Name
	// in `json:"id,omitempty"`, the Options are [omitempty]
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Options are [embeddedPrefix:author_]
	// in `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, the Options are
	// [constraint:OnUpdate:CASCADE,OnDelete:SET NULL;] and the Name is ""
	Options []*Option
}

// Option is an option of a tag value
type Option struct {
	// Key is the option name: omitempty, default
	Key string
	// Value is the string after the first : of the option, or nil when the option has no value
	Value any
}

// HasOption reports whether the tag has the option
func (t *Tag) HasOption(option string) bool {
	return t.GetOption(option) != nil
}

// GetOption returns the first option with this key, or nil
func (t *Tag) GetOption(option string) *Option {
	for _, tagOption := range t.Options {
		if tagOption.Key == option {
			return tagOption
		}
	}
	return nil
}
