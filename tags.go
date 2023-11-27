package tags

// Tag represents a string literal applied to a struct field
// see: https://pkg.go.dev/reflect#StructTag
type Tag struct {
	// Tag is the full string containing the Key, the Name and Options
	Tag string

	// Key is the tag key which can be obtained with func (StructTag) Get
	// in `json:"id,omitempty"`, the Key is "json"
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Key is "gorm"
	// in `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, the Key is "gorm"
	Key string

	// Value is obtained with func (StructTag) Lookup
	// in `json:"id,omitempty"`, the Value is "id,omitempty"
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Value is "embedded;embeddedPrefix:author_"
	// in `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, the Value is "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
	Value string

	// Name is the first part of the value obtained with func (StructTag) Lookup
	// in `json:"id,omitempty"`, the key Name "id"
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Name is "embedded"
	// in `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, the Name is "constraint"
	Name string

	// Options are the second part of the value obtained with func (StructTag) Lookup
	// Options is a list of options (Key, Value pair)
	// in `json:"id,omitempty"`, the Options is ["omitempty"]
	// in `gorm:"embedded;embeddedPrefix:author_"`, the Options is ["embeddedPrefix:author_"]
	// in `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, the Options is ["OnUpdate:CASCADE", "OnDelete:SET NULL;"]
	Options []*Option
}

// Option is a simple key, value pair
type Option struct {
	Key   string
	Value any
}

// HasOption return true if the option is found
func (tags *Tag) HasOption(option string) bool {
	return tags.GetOption(option) != nil
}

// GetOption return the option or nil
func (tags *Tag) GetOption(option string) *Option {

	for _, tagOption := range tags.Options {
		if tagOption.Key == option {
			return tagOption
		}
	}
	return nil
}
